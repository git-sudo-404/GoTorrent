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
	"fmt"
	"net"
	"os"
	"time"
)

// NOTE: additional functions like getCompletePieces(), getIncompletePieces(), or getPartialPieces(), etc./ might shall be implemented as needed
type PieceBuffer struct {
	block              [][]byte
	offset             []int
	bitfield           []byte // represents the completely downloaded pieces if the bit is set
	uploadedBytes      uint64
	downloadedBytes    uint64
	pieceToPeer        [][20]byte // represents piece -> peerId mapping denoting a particular piece has been requested by the mentioned peer
	downloadInProgress []byte     // represents whether or not a piece is being downloaded block by block by a peer , it is set as long as the peer connection is alive and is being actively downloaded
}

func NewPieceBuffer(metaInfo *MetaInfo) *PieceBuffer {
	pieceLength, _ := metaInfo.GetPieceLength()
	pieceNums := metaInfo.length / pieceLength

	blocks := make([][]byte, pieceNums)
	for i := range blocks {
		blocks[i] = make([]byte, pieceLength)
	}
	offset := make([]int, pieceNums)
	for i := 0; i < int(pieceNums); i++ {
		offset[i] = 0
	}
	return &PieceBuffer{
		block:              blocks,
		offset:             offset,
		bitfield:           make([]byte, pieceNums),
		uploadedBytes:      0,
		downloadedBytes:    0,
		pieceToPeer:        make([][20]byte, pieceNums),
		downloadInProgress: make([]byte, pieceNums),
	}
}

func (pb *PieceBuffer) writeBlcoksToFile(metaInfo *MetaInfo) {
	fileName := metaInfo.name
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println("[ERROR] Error creating the downloaded file")
		os.Exit(1)
	}
	writer := bufio.NewWriter(file)
	for pieceIndex, _ := range pb.block {
		writer.Write(pb.block[pieceIndex])
	}
	fmt.Println("[LOG] File Downloaded...")
}

type Peer struct {
	ip              net.IP
	port            int64
	peerId          [20]byte
	conn            net.Conn
	bitfield        []byte
	am_choking      bool
	am_interested   bool
	peer_choking    bool
	peer_interested bool
	alive           time.Time
	reqBlock        chan struct{}
}

type Client struct {
	clientId        [20]byte
	port            int64
	peers           []Peer
	trackerInterval int64
	completePeers   int64
	incompletePeers int64
}

func NewClient(metaInfo *MetaInfo) *Client {
	var peerId [20]byte
	prefix := []byte("-GT0001-")

	copy(peerId[:], prefix)

	processId := os.Getpid()
	timeStamp := time.Now().UnixMilli()

	data := fmt.Sprintf("%d-%d", processId, timeStamp)
	hash := sha1.Sum([]byte(data))

	copy(peerId[8:], hash[:12])

	return &Client{
		clientId:        peerId,
		port:            6881,
		peers:           []Peer{},
		trackerInterval: -1,
		completePeers:   -1,
		incompletePeers: -1,
	}
}
