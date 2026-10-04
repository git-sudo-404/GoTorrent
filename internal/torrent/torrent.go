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
	"crypto/sha1"
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

	fmt.Printf("[DEBUG] peers type: %T\n", peers)
	fmt.Printf("[DEBUG] peers value: %#v\n", peers)

	if peerList, ok := peers.([]any); ok {
		fmt.Println("[LOG] Parsing peers address from peers map")
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
					ip:       net.ParseIP(ip),
					port:     port,
					reqBlock: make(chan struct{}),
					bitfield: make([]byte, bitfieldLen),
				})
			}
		}
	} else if peerString, ok := peers.(string); ok {
		fmt.Println("[LOG] Parsing peers address from Peer String")
		if len(peerString)%6 != 0 {
			fmt.Println("[ERROR] peerString returned by Tracker is incomplete")
			os.Exit(1)
		}
		offset := 0
		bitfieldLen := (len(metaInfo.pieces) + 7) / 8
		for {
			if offset >= len(peerString) {
				break
			}
			port, _ := strconv.ParseInt(peerString[offset+4:offset+6], 10, 16)
			client.peers = append(client.peers, Peer{
				ip:       net.ParseIP(peerString[offset : offset+4]),
				port:     int64(port),
				reqBlock: make(chan struct{}),
				bitfield: make([]byte, bitfieldLen),
			})
			offset += 6
		}
	} else {
		fmt.Println("[ERROR] Wrong data type in peerList")
	}
	fmt.Println("[LOG] Updated Client Peers from tracker Response")
}

func sendHandshakeRequest(wg *sync.WaitGroup, peer *Peer, metaInfo *MetaInfo, peerHandshakeRequestMsg [68]byte, handshakeRetriesChan chan [20]byte) {

	defer wg.Done()

	conn, err := net.Dial("tcp", net.JoinHostPort(peer.ip.String(), strconv.Itoa(int(peer.port))))
	if err != nil {
		fmt.Println("[ERROR] Handshake Failed:", err)
	}

	peer.conn = conn

	fmt.Printf("[LOG] Initiating Handshake with  Peer : %s on Port : %d\n", peer.ip.String(), peer.port)
	n, err := conn.Write(peerHandshakeRequestMsg[:])
	if err != nil {
		fmt.Println("[ERROR] Write Failed", err)
		handshakeRetriesChan <- peer.peerId
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

func retryHandhshake(retryWG *sync.WaitGroup, peer *Peer, handshakeRequestMsg [68]byte, retryCount int, removePeerChan chan [20]byte) {

	for i := 0; i < retryCount; i++ {
		n, err := peer.conn.Write(handshakeRequestMsg[:])
		if err == nil && n == 68 {
			return
		}
	}
	removePeerChan <- peer.peerId
}

func handFailedHandshakes(handshakeRetriesChan chan [20]byte, metaInfo *MetaInfo, client *Client) {
	var retryWG sync.WaitGroup
	removePeerChan := make(chan [20]byte, len(client.peers))
	handshakeRequestMsg := NewHandshakeRequest(metaInfo, client)
	for failedPeerId := range handshakeRetriesChan {
		var peer *Peer
		for peerIndex, _ := range client.peers {
			for i := 0; i < 20; i++ {
				if client.peers[peerIndex].peerId[i] != failedPeerId[i] {
					goto checkNextPeer
				}
			}
			peer = &client.peers[peerIndex]
			break
		checkNextPeer:
		}
		retryHandhshake(&retryWG, peer, handshakeRequestMsg, 3, removePeerChan)
	}
	for failedPeerId := range removePeerChan {
		for peerIndex, _ := range client.peers {
			for i := 0; i < 20; i++ {
				if client.peers[peerIndex].peerId[i] != failedPeerId[i] {
					goto nextPeer
				}
			}
			// just a faster way to remove the failed peer from the slice , by replacing it with the last element and then slicing off the last element
			client.peers[peerIndex] = client.peers[len(client.peers)-1]
			client.peers = client.peers[:len(client.peers)-1]
			fmt.Println("[LOG] Removed Peer from Client Peers List after 3 failed handshake attempts")
			break
		nextPeer:
		}
	}
}

func sendHandshakeRequestToAllPeers(client *Client, metaInfo *MetaInfo) {
	fmt.Println("[LOG] Initiating Peer Handshake ...")

	handshakeRequest := NewHandshakeRequest(metaInfo, client)

	var wg sync.WaitGroup
	handshakeRetriesChan := make(chan [20]byte, len(client.peers))

	//TODO: Add logic to handle the failing handshakes , increase the timeouts , do 3 retries , use cahnnels to collect the failing peers alone and retry them

	for peerIndex, _ := range client.peers {
		peer := &client.peers[peerIndex]
		wg.Add(1)
		go sendHandshakeRequest(&wg, peer, metaInfo, handshakeRequest, handshakeRetriesChan)
	}

	go handFailedHandshakes(handshakeRetriesChan, metaInfo, client)

	wg.Wait()
}

func writeToPeerWithRetries(peer *Peer, msg []byte, retryCount int) {
	for i := 0; i < retryCount; i++ {
		n, err := peer.conn.Write(msg)
		if err != nil || n < len(msg) {
			fmt.Println("[LOG] An error occured while writing to peer , retrying again ...")
			continue
		}
		fmt.Println("[LOG] Sent ", len(msg), " bytes to peer")
		return
	}
}

func serveRequestedBlockToPeer(reqMsg [12]byte, peer *Peer, pieceBuffer *PieceBuffer) {

	pieceIndex := binary.BigEndian.Uint32(reqMsg[0:4])
	blockBegin := binary.BigEndian.Uint32(reqMsg[4:8])
	blockLen := binary.BigEndian.Uint32(reqMsg[8:12])
	piecePayload := pieceBuffer.block[pieceIndex][blockBegin : blockBegin+blockLen]
	pieceMessage := make([]byte, blockLen+4+1)
	binary.BigEndian.AppendUint32(pieceMessage[0:4], blockLen+1)
	pieceMessage[4] = 7
	copy(pieceMessage[5:5+blockLen], piecePayload)

	writeToPeerWithRetries(peer, piecePayload, 3)

}

func verifyPieceHash(pieceBuffer *PieceBuffer, pieceIndex int, metaInfo *MetaInfo) bool {
	pieceSha1Sum := sha1.Sum(pieceBuffer.block[pieceIndex])
	if pieceSha1Sum != metaInfo.pieces[pieceIndex] {
		return false
	}
	return true
}

func handlePeer(peerMessagingWG *sync.WaitGroup, peer *Peer, metaInfo *MetaInfo, pieceBuffer *PieceBuffer) {
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
				fmt.Println("[LOG] Recieved a request message ...")
				reqMsg := [12]byte{}
				io.ReadFull(peer.conn, reqMsg[:])
				if peer.peer_choking {
					continue
				}
				serveRequestedBlockToPeer(reqMsg, peer, pieceBuffer)
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
				fmt.Println("[LOG] Starting to request blocks from peer after a bitfield msg...")
				peer.reqBlock <- struct{}{}
			case 7: // piece msg
				// fmt.Println("[LOG] Recieved a piece message...")
				pieceMsgPayload := make([]byte, msgLength-1) // since , 1 byte is already consumed in msgId
				io.ReadFull(peer.conn, pieceMsgPayload)
				pieceIndex := binary.BigEndian.Uint32(pieceMsgPayload[0:4])
				blockBegin := binary.BigEndian.Uint32(pieceMsgPayload[4:8])
				blockLen := msgLength - (1 + 4 + 4)
				for i := range blockLen {
					pieceBuffer.block[pieceIndex][blockBegin+i] = pieceMsgPayload[8+i]
				}
				pieceBuffer.offset[int(pieceIndex)] = int(blockBegin + blockLen)
				// fmt.Println("[LOG] Piece successfully recieved : ", blockLen, " bytes")
				pieceBuffer.downloadedBytes += uint64(blockLen)
				if pieceBuffer.offset[pieceIndex] == int(metaInfo.pieceLength) {
					if !verifyPieceHash(pieceBuffer, int(pieceIndex), metaInfo) {
						fmt.Println("[LOG] Downloaded Piece Hash didn't match , discarding it ...")
						pieceBuffer.offset[pieceIndex] = 0
						pieceBuffer.downloadedBytes -= uint64(metaInfo.pieceLength)
						pieceBuffer.pieceToPeer[pieceIndex] = [20]byte{}
						return
					}
					fmt.Println("[LOG] piece ", pieceIndex, " successfully downloaded")
					setBit(int64(pieceIndex), pieceBuffer.bitfield)
				} else {
					// fmt.Println("[LOG] Requesting next block...")
				}
				peer.reqBlock <- struct{}{}
				// fmt.Println("[LOG] Total Downloaded Bytes : ", pieceBuffer.downloadedBytes, " Remaining : ", metaInfo.length-int64(pieceBuffer.downloadedBytes))
				fmt.Println("[LOG] ", float64(pieceBuffer.downloadedBytes)/float64(metaInfo.length)*100, "% downlaoded")
				if pieceBuffer.downloadedBytes == uint64(metaInfo.length) {
					fmt.Println("[LOG] File Downloaded....")
					pieceBuffer.writeBlcoksToFile(metaInfo)
				}
			}
		}
	}
}

func handleMessageFromPeers(torrentWG *sync.WaitGroup, client *Client, metaInfo *MetaInfo, pieceBuffer *PieceBuffer) {
	fmt.Println("[LOG] Starting Peer Message Handling...")
	var peerMessagingWG sync.WaitGroup
	for peerIndex, _ := range client.peers {
		peer := &client.peers[peerIndex]
		peerMessagingWG.Add(1)
		go handlePeer(&peerMessagingWG, peer, metaInfo, pieceBuffer)
	}
	peerMessagingWG.Wait()
	torrentWG.Done()
}

func downloadFromPeer(downloadBlocksWG *sync.WaitGroup, peer *Peer, metaInfo *MetaInfo, pieceBuffer *PieceBuffer) {
	defer downloadBlocksWG.Done()
	for range peer.reqBlock {
		if peer.am_choking {
			continue
		}
		for i := 0; i < len(peer.bitfield); i++ {
			for j := 0; j < 8; j++ {
				pieceIndex := uint32(i*8 + j)
				// check if the client already downloaded that piece
				if hasBit(int64(pieceIndex), pieceBuffer.bitfield) {
					continue
				}
				// check if the peer has that piece
				if !hasBit(int64(pieceIndex), peer.bitfield) {
					continue
				}
				// check if the piece is already being dowbloaded
				if hasBit(int64(pieceIndex), pieceBuffer.downloadInProgress) {
					// if already being downloaded , check if piece is already being downloaded by some other peer
					if pieceBuffer.pieceToPeer[pieceIndex] != peer.peerId && pieceBuffer.pieceToPeer[pieceIndex] != [20]byte{} {
						continue
					}
				}
				if pieceBuffer.pieceToPeer[pieceIndex] == [20]byte{} {
					pieceBuffer.pieceToPeer[pieceIndex] = peer.peerId
					setBit(int64(pieceIndex), pieceBuffer.downloadInProgress)
				}
				reqMsgBytes := [17]byte{}
				binary.BigEndian.PutUint32(reqMsgBytes[0:4], uint32(13))
				reqMsgBytes[4] = byte(6) // msgId
				binary.BigEndian.PutUint32(reqMsgBytes[5:9], pieceIndex)
				binary.BigEndian.PutUint32(reqMsgBytes[9:13], uint32(pieceBuffer.offset[int(pieceIndex)]))
				binary.BigEndian.PutUint32(reqMsgBytes[13:17], uint32(1<<14)) // 16 KB
				_, err := peer.conn.Write(reqMsgBytes[:])
				if err != nil {
					fmt.Println("[ERROR] Error when writing req msg", err)
				}
				// fmt.Println("[LOG] Sent a piece req message : ", reqMsgBytes)
				goto waitUntilNextReqCanBeSent
			}
		}
	waitUntilNextReqCanBeSent:
	}
}

func DownloadBlocks(torrentWG *sync.WaitGroup, client *Client, metaInfo *MetaInfo, pieceBuffer *PieceBuffer) {
	var downloadBlocksWG sync.WaitGroup
	interestedMessage := NewInterestedMessage()
	for peerIndex, _ := range client.peers {
		peer := &client.peers[peerIndex]
		n, err := peer.conn.Write(interestedMessage[:])
		if err != nil {
			fmt.Println("[ERROR] Can't send interested message to peer", err)
			return
		}
		if n != 5 {
			fmt.Println("[ERROR] Can't write the full interested message to peer , sent ", n, " bytes")
		}
		downloadBlocksWG.Add(1)
		go downloadFromPeer(&downloadBlocksWG, peer, metaInfo, pieceBuffer)
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
	pieceBuffer := NewPieceBuffer(metaInfo)

	sendTrackerRequestAndHandleResponse(client, metaInfo)

	sendHandshakeRequestToAllPeers(client, metaInfo)

	var torrentWG sync.WaitGroup
	torrentWG.Add(1)
	go DownloadBlocks(&torrentWG, client, metaInfo, pieceBuffer)
	torrentWG.Add(1)
	go handleMessageFromPeers(&torrentWG, client, metaInfo, pieceBuffer)
	torrentWG.Wait()
}
