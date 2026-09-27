/*
 * MIT License
 *
 * Copyright (c) 2026 git-sudo-404
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * Of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * Copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * Copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING, BUT NOT LIMITED TO, THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES, OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package torrent

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"gotorrent/internal/bencode"
	"io"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"
)

func populateInitialTrackerRequestParams(tr *TrackerRequest, mi *MetaInfo, client *Client) {
	fmt.Println("[LOG] Populating TrackerRequest...")
	tr.SetInfoHash(mi)
	tr.SetPeerId(string(client.clientId[:]))
	tr.SetPort(int64(client.port))
	tr.SetUploaded(int64(0))
	tr.SetDownloaded(int64(0))
	tr.SetLeft(mi.length)
	tr.SetCompact(int64(0))
	tr.SetNoPeerId(int64(1))
	tr.SetEvent(STARTED)
}

func sendTrackerRequest(trackerRequestURL *string, client *Client) {

	httpClient := &http.Client{
		Timeout: 5 * time.Second,
	}
	req, err := http.NewRequest(http.MethodGet, *trackerRequestURL, nil)
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println("[LOG] Sending TrackerRequest...")
	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Println(err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error tracker server return : ", resp.Body)
	}

	fmt.Println("[LOG] STATUS:", resp.Status)

	buf := bufio.NewReader(resp.Body)

	decodedResponse, err := bencode.Decode(buf)
	if err != nil {
		panic(err)
	}

	// fmt.Println(decodedResponse)
	interval, present := decodedResponse["interval"]
	if !present {
		panic("interval not present in the tracker response")
	}
	//NOTE: Tracker Id is optioanl
	// trackerId, present := decodedResponse["tracker id"]
	// if !present {
	// 	panic("tracker id not present in tracker response")
	// }
	complete, present := decodedResponse["complete"]
	if !present {
		panic("complete field not present in tracker response")
	}
	incomplete, present := decodedResponse["incomplete"]
	if !present {
		panic("Field : incomplete , not present in tracker response")
	}
	peers, present := decodedResponse["peers"]
	if !present {
		panic("peers not present in tracker response")
	}

	// fmt.Println("[LOG] Tracker Responded with the following : ....")
	// fmt.Println("[LOG] interval : ", interval)
	// fmt.Println("[LOG] complete : ", complete)
	// fmt.Println("[LOG] incomplete : ", incomplete)
	// fmt.Println("[LOG] peers : ", peers)

	if val, ok := interval.(int64); ok {
		client.trackerInterval = val
	} else {
		fmt.Println("[ERROR] : interval of tracker response not of type int64")
	}

	if val, ok := complete.(int64); ok {
		client.completePeers = val
	} else {
		fmt.Println("[ERROR] : complete from tracker response is not of type int64")
	}

	if val, ok := incomplete.(int64); ok {
		client.incompletePeers = val
	} else {
		fmt.Println("[ERROR] : incomplete of tracker response not of type int64")
	}

	if peerList, ok := peers.([]any); ok {
		// fmt.Println("[DEBUG PRINT] PEER LIST : ", peerList)
		for _, peerListItem := range peerList {
			// fmt.Println("[DEBUG PRINT] PEERLISTITEM : ", peerListItem)
			if peerMap, ok := peerListItem.(map[string]any); ok {
				ip, ok := peerMap["ip"].(string)
				// fmt.Println("[DEBUG PRINT] IP : ", peerMap["ip"])
				// fmt.Println("[DEBUG PRINT] PORT : ", peerMap["port"])
				if !ok {
					fmt.Println("[ERROR] ip not of type string from tracker response")
				}
				port, ok := peerMap["port"].(int64)
				if !ok {
					fmt.Println("[ERROR] port not of type integer from tracker response")
				}
				client.peers = append(client.peers, Peer{
					ip:   net.ParseIP(ip),
					port: port,
				})
			}
		}
	} else {
		fmt.Println("[ERROR] Wrong data type in peerList")
	}

	fmt.Println("[LOG] Updated Client Peers from tracker Response")

}

func intitaiteHandshakeWithPeers(client *Client, metaInfo *MetaInfo) {
	fmt.Println("[LOG] Initiating Peer Handshake ...")
	peerHandshakeRequest := NewPeerHandshakeRequest()

	pstr := []byte("BitTorrent protocol") //NOTE: This string is not arbitrary
	reserved := make([]byte, 8)
	info_hash := metaInfo.GetInfoHash()
	peer_id := client.clientId

	var wg sync.WaitGroup

	peerHandshakeRequest.SetPstr(pstr)
	peerHandshakeRequest.SetReserved([8]byte(reserved))
	peerHandshakeRequest.SetInfoHash(info_hash)
	peerHandshakeRequest.SetPeerId(peer_id)

	peerHandshakeRequestBytes := peerHandshakeRequest.Serialize()

	for peerIndex, peer := range client.peers {
		wg.Add(1)
		go func() {
			conn, err := net.Dial("tcp", net.JoinHostPort(peer.ip.String(), strconv.Itoa(int(peer.port))))
			if err != nil {
				fmt.Println("[ERROR] Handshake Failed:", err)
				return
			}
			client.peers[peerIndex].conn = conn
			defer func() {
				wg.Done()
			}()
			fmt.Printf("[LOG] Initiating Handshake with  Peer : %s on Port : %d\n", peer.ip.String(), peer.port)
			n, err := conn.Write(peerHandshakeRequestBytes)
			if err != nil {
				fmt.Println("[ERROR] Write Failed", err)
				return
			}
			fmt.Println("[DEBUG] Sent", n, "bytes")
			readBuffer := make([]byte, 68)
			n, err = conn.Read(readBuffer)
			if err != nil {
				fmt.Println("[ERROR] Read Failed", err)
				return
			}
			fmt.Println("[DEBUG] Recieved", n, "bytes")

			peerInfoHash := readBuffer[28:48]
			peerId := readBuffer[48:68]

			for i := 0; i < 20; i++ {
				if peerInfoHash[i] != metaInfo.infoHash[i] {
					fmt.Println("[ERROR] Info Hash Not matching with peer")
					return
				}
			}

			copy(client.peers[peerIndex].peerId[:], peerId[:])
			fmt.Println("[LOG] Recieved Peer Id : ", client.peers[peerIndex].peerId)

			bitfieldLen := (len(metaInfo.pieces) + 7) / 8
			client.peers[peerIndex].bitfield = make([]byte, bitfieldLen)

		}()
	}

	wg.Wait()

}

// func getPeerBitFeilds(client *Client) {
// 	for peerIndex, _ := range client.peers {
//
// 		peer := &client.peers[peerIndex]
//
// 		lenBuffer := [4]byte{}
// 		_, err := io.ReadFull(peer.conn, lenBuffer[:])
// 		if err != nil {
// 			fmt.Println("[ERROR] Error reading the length Prefix of the bitfield message", err)
// 		}
//
// 		messageId := [1]byte{}
// 		io.ReadFull(peer.conn, messageId[:])
//
// 		if binary.BigEndian.Uint32(messageId[:]) != 5 {
// 			return
// 		}
//
// 		length := binary.BigEndian.Uint32(lenBuffer[:])
// 		bitfieldRawBytes := make([]byte, length)
//
// 		_, err = io.ReadFull(peer.conn, bitfieldRawBytes)
// 		if err != nil {
// 			fmt.Println("[ERROR] Error while reading bitfield message from peer", err)
// 		}
//
// 	}
// }

func handlePeer(wg *sync.WaitGroup, peer *Peer) {
	defer func() {
		wg.Done()
	}()

	for {

		// read the 4 byte message len
		msgLengthBytes := [4]byte{}
		io.ReadFull(peer.conn, msgLengthBytes[:])
		msgLength := binary.BigEndian.Uint32(msgLengthBytes[:])

		if msgLength == 0 { // keep-alive msg
			peer.alive = time.Now()
		}

		msgIdBytes := [1]byte{}
		io.ReadFull(peer.conn, msgIdBytes[:])
		msgId := binary.BigEndian.Uint32(msgIdBytes[:])

		switch msgLength {
		case 1:
			switch msgId {
			case 0: // choke
				peer.am_choking = true
			case 1: // unchoke
				peer.am_choking = false
			case 2: // interested
				peer.peer_interested = true
			case 3: // not-interested
				peer.peer_interested = false
			}
		case 5:
			if msgId == 4 { // have msg
				pieceIndexBytes := [4]byte{}
				io.ReadFull(peer.conn, pieceIndexBytes[:])
				pieceIndex := binary.BigEndian.Uint32(pieceIndexBytes[:])
				setBit(int64(pieceIndex), peer.bitfield)
			}
		}

	}

}

func startPeerMessaging(client *Client, metaInfo *MetaInfo) {
	var wg sync.WaitGroup
	for peerIndex, _ := range client.peers {
		wg.Add(1)
		go handlePeer(&wg, &client.peers[peerIndex])
	}
	wg.Wait()
}

func StartTorrent(metaInfoFilePath string, destinationFilePath string) {

	fmt.Println("[LOG] Starting Torrent ...")

	fmt.Println("[LOG] Parsing MetaInfo")
	metaInfo, err := CreateMetaInfoFromFile(metaInfoFilePath)
	if err != nil {
		panic(err)
	}

	client := NewClient()
	trackerRequest := NewTrackerRequest()

	populateInitialTrackerRequestParams(trackerRequest, metaInfo, client)

	trackerRequestURL, err := trackerRequest.GetURLEncodedRequestString(metaInfo)
	if err != nil {
		panic(err)
	}

	sendTrackerRequest(&trackerRequestURL, client)

	//NOTE: This function handles only the client peers's handshake , need to implement the logic to handle the incoming handshake requests from other peers
	intitaiteHandshakeWithPeers(client, metaInfo)

	startPeerMessaging(client, metaInfo)

}
