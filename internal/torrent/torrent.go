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
	"os"
	"strconv"
	"sync"
	"time"
)

func sendTrackerRequestAndHandleResponse(client *Client, metaInfo *MetaInfo) {

	trackerRequest := NewTrackerRequest()
	fmt.Println("[LOG] Populating TrackerRequest...")
	trackerRequest.SetInfoHash(metaInfo)
	trackerRequest.SetPeerId(string(client.clientId[:]))
	trackerRequest.SetPort(int64(client.port))
	trackerRequest.SetUploaded(int64(0))
	trackerRequest.SetDownloaded(int64(0))
	trackerRequest.SetLeft(metaInfo.length)
	trackerRequest.SetCompact(int64(0))
	trackerRequest.SetNoPeerId(int64(1))
	trackerRequest.SetEvent(STARTED)
	trackerRequestURL, err := trackerRequest.GetURLEncodedRequestString(metaInfo)
	if err != nil {
		panic(err)
	}

	httpClient := &http.Client{
		Timeout: 15 * time.Second,
	}
	req, err := http.NewRequest(http.MethodGet, trackerRequestURL, nil)
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
	// fmt.Println("[LOG] STATUS:", resp.Status)

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
		for _, peerListItem := range peerList {
			if peerMap, ok := peerListItem.(map[string]any); ok {
				ip, ok := peerMap["ip"].(string)
				if !ok {
					fmt.Println("[ERROR] ip not of type string from tracker response")
				}
				port, ok := peerMap["port"].(int64)
				if !ok {
					fmt.Println("[ERROR] port not of type integer from tracker response")
				}
				bitfieldLen := (len(metaInfo.pieces) + 7) / 8
				client.peers = append(client.peers, Peer{
					ip:           net.ParseIP(ip),
					port:         port,
					mu:           sync.Mutex{},
					canReqBlocks: make(chan struct{}),
					bitfield:     make([]byte, bitfieldLen),
				})
			}
		}
	} else {
		fmt.Println("[ERROR] Wrong data type in peerList")
	}
	fmt.Println("[LOG] Updated Client Peers from tracker Response")
}

func sendHandshakeRequest(wg *sync.WaitGroup, peer *Peer, metaInfo *MetaInfo, peerHandshakeRequestMsg [68]byte) {

	defer wg.Done()

	conn, err := net.Dial("tcp", net.JoinHostPort(peer.ip.String(), strconv.Itoa(int(peer.port))))
	if err != nil {
		fmt.Println("[ERROR] Handshake Failed:", err)
		return
	}

	peer.conn = conn

	fmt.Printf("[LOG] Initiating Handshake with  Peer : %s on Port : %d\n", peer.ip.String(), peer.port)
	n, err := conn.Write(peerHandshakeRequestMsg[:])
	if err != nil {
		fmt.Println("[ERROR] Write Failed", err)
		return
	}
	fmt.Println("[DEBUG] Sent", n, "bytes")

	readBuffer := make([]byte, 68)
	n, err = io.ReadFull(conn, readBuffer)
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

	copy(peer.peerId[:], peerId[:])
	fmt.Println("[LOG] Recieved Peer Id ")
}

func sendHandshakeRequestToAllPeers(client *Client, metaInfo *MetaInfo) {
	fmt.Println("[LOG] Initiating Peer Handshake ...")

	pstr := []byte("BitTorrent protocol") //NOTE: This string is not arbitrary
	reserved := make([]byte, 8)
	info_hash := metaInfo.GetInfoHash()
	peer_id := client.clientId

	peerHandshakeRequestMsg := [68]byte{}
	peerHandshakeRequestMsg[0] = uint8(19)
	copy(peerHandshakeRequestMsg[1:20], pstr)
	copy(peerHandshakeRequestMsg[20:28], reserved)
	copy(peerHandshakeRequestMsg[28:48], info_hash[:])
	copy(peerHandshakeRequestMsg[48:68], peer_id[:])

	var wg sync.WaitGroup

	for peerIndex, _ := range client.peers {
		peer := &client.peers[peerIndex]
		wg.Add(1)
		go sendHandshakeRequest(&wg, peer, metaInfo, peerHandshakeRequestMsg)
	}

	wg.Wait()
}

func handlePeer(peerMessagingWG *sync.WaitGroup, peer *Peer, metaInfo *MetaInfo, blocks [][]byte, blockOffset map[int]int, downloadedBytes *uint64, uploadedBytes *uint64) {
	defer peerMessagingWG.Done()
	for {
		// read the 4 byte message len
		msgLengthBytes := [4]byte{}
		io.ReadFull(peer.conn, msgLengthBytes[:])
		msgLength := binary.BigEndian.Uint32(msgLengthBytes[:])

		if msgLength == 0 { // keep-alive msg
			peer.alive = time.Now()
			continue
		}

		msgIdByte := [1]byte{}
		io.ReadFull(peer.conn, msgIdByte[:])
		msgId := msgIdByte[0]

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
		case 13:
			if msgId == 6 { // request msg
				fmt.Println("[LOG] Reecieved a request message ...")
				reqMsgPayload := [12]byte{}
				io.ReadFull(peer.conn, reqMsgPayload[:])
				// index := binary.BigEndian.Uint32(reqMsgPayload[0:4])
				// begin := binary.BigEndian.Uint32(reqMsgPayload[4:8])
				// length := binary.BigEndian.Uint32(reqMsgPayload[8:12])
				// // piecePayload := blocks[index][begin : begin+length]
				// pieceMessage := make([]byte,length + 4 + 1)
			} else if msgId == 8 {
				cancelMsgPayload := make([]byte, msgLength-1)
				io.ReadFull(peer.conn, cancelMsgPayload)
				//TODO: implement the cancel logic
			}
		default:
			switch msgId {
			case 5: // bitfield msg
				fmt.Println("[LOG] Recieved a bitfield Msg")
				bitfieldBytes := make([]byte, msgLength-1)
				io.ReadFull(peer.conn, bitfieldBytes[:])
				peer.bitfield = bitfieldBytes
				peer.canReqBlocks <- struct{}{}
			case 7: // piece msg
				fmt.Println("[LOG] Recieved a piece message...")
				pieceMsgPayload := make([]byte, msgLength-1) // since , 1 byte is already consumed in msgId
				io.ReadFull(peer.conn, pieceMsgPayload)
				index := binary.BigEndian.Uint32(pieceMsgPayload[0:4])
				begin := binary.BigEndian.Uint32(pieceMsgPayload[4:8])
				blockLen := msgLength - (1 + 4 + 4)
				for i := range blockLen {
					blocks[index][begin+i] = pieceMsgPayload[8+i]
				}
				blockOffset[int(index)] = int(begin + blockLen)
				fmt.Println("[LOG] Piece successfully recieved ")
				*downloadedBytes += uint64(blockLen)
				if blockOffset[int(index)] < int(metaInfo.pieceLength) {
					fmt.Println("[LOG] Request next block offsett")
					peer.canReqBlocks <- struct{}{}
				}
				fmt.Println("[LOG] Downloaded Bytes : ", *downloadedBytes)
			}
		}
	}
}

func handlePeerMessaging(torrentWG *sync.WaitGroup, client *Client, metaInfo *MetaInfo) {
	fmt.Println("[LOG] Starting Peer Message Handling...")
	var peerMessagingWG sync.WaitGroup
	for peerIndex, _ := range client.peers {
		peer := &client.peers[peerIndex]
		peerMessagingWG.Add(1)
		go handlePeer(&peerMessagingWG, peer, metaInfo, client.blocks, client.blockOffset, &client.downloadedBytes, &client.uploadedBytes)
	}
	peerMessagingWG.Wait()
	torrentWG.Done()
}

func downloadFromPeer(downloadBlocksWG *sync.WaitGroup, peer *Peer, metaInfo *MetaInfo, blocks [][]byte, blockOffset map[int]int, downloadedBytes *uint64, uploadedBytes *uint64) {
	defer downloadBlocksWG.Done()
	for range peer.canReqBlocks {
		fmt.Println("[LOG] Requested a block from peer")
		for i := 0; i < len(peer.bitfield); i++ {
			for j := 0; j < 8; j++ {
				blockIndex := uint32(i*8 + j)
				if !hasBit(int64(blockIndex), peer.bitfield) {
					continue
				}
				if blockOffset[int(blockIndex)] >= int(metaInfo.pieceLength) {
					continue
				}
				reqMsgBytes := [17]byte{}
				binary.BigEndian.PutUint32(reqMsgBytes[0:4], uint32(13))
				reqMsgBytes[4] = byte(6) // msgId
				binary.BigEndian.PutUint32(reqMsgBytes[5:9], blockIndex)
				binary.BigEndian.PutUint32(reqMsgBytes[9:13], uint32(blockOffset[int(blockIndex)]))
				binary.BigEndian.PutUint32(reqMsgBytes[13:17], uint32(1<<14)) // 16 KB
				_, err := peer.conn.Write(reqMsgBytes[:])
				if err != nil {
					fmt.Println("[ERROR] Error when writing req msg", err)
				}
				fmt.Println("[LOG] Sent a piece req message : ", reqMsgBytes)
				goto nextTime
			}
		}
	nextTime:
	}
}

func DownloadBlocks(torrentWG *sync.WaitGroup, client *Client, metaInfo *MetaInfo) {
	var downloadBlocksWG sync.WaitGroup
	interestedMessage := NewInterestedMessage()
	for peerIndex, _ := range client.peers {
		peer := &client.peers[peerIndex]
		downloadBlocksWG.Add(1)
		peer.conn.Write(interestedMessage[:])
		go downloadFromPeer(&downloadBlocksWG, peer, metaInfo, client.blocks, client.blockOffset, &client.downloadedBytes, &client.uploadedBytes)
	}
	downloadBlocksWG.Wait()
	torrentWG.Done()
}

func StartTorrent(metaInfoFilePath string, destinationFilePath string) {

	fmt.Println("[LOG] Parsing MetaInfo...")
	metaInfo, err := CreateMetaInfoFromFile(metaInfoFilePath)
	if err != nil {
		fmt.Println("[ERROR] Error when parsing MetaInfo file : ", err)
		os.Exit(1)
	}

	client := NewClient(metaInfo)

	sendTrackerRequestAndHandleResponse(client, metaInfo)

	sendHandshakeRequestToAllPeers(client, metaInfo)

	var torrentWG sync.WaitGroup
	torrentWG.Add(1)
	go DownloadBlocks(&torrentWG, client, metaInfo)
	torrentWG.Add(1)
	go handlePeerMessaging(&torrentWG, client, metaInfo)
	torrentWG.Wait()
}
