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

import "encoding/binary"

// Helper Function to set/unset bitfield bits
func setBit(pieceIndex int64, bitfield []byte) {
	byteIndex := pieceIndex / 8
	bitIndex := pieceIndex % 8
	bitfield[byteIndex] |= 1 << bitIndex
}

func unsetBit(pieceIndex int64, bitfield []byte) {
	byteIndex := pieceIndex / 8
	bitIndex := pieceIndex % 8
	bitfield[byteIndex] ^= 1 << bitIndex
}

func hasBit(pieceIndex int64, bitfield []byte) bool {
	byteIndex := pieceIndex / 8
	bitIndex := pieceIndex % 8
	return bitfield[byteIndex]&(1<<bitIndex) == 1
}

// keep-alive: <len=0000>
// The keep-alive message is a message with zero bytes, specified with the length prefix set to zero. There is no message ID and no payload.
func NewKeepAliveMessage() [4]byte {
	return [4]byte{0, 0, 0, 0}
}

// choke: <len=0001><id=0>
// The choke message is fixed-length and has no payload.
func NewChokeMessage() [5]byte {
	return [5]byte{0, 0, 0, 1, 0}
}

// unchoke: <len=0001><id=1>
// The unchoke message is fixed-length and has no payload.
func NewUnchokeMessage() [5]byte {
	return [5]byte{0, 0, 0, 1, 1}
}

// interested: <len=0001><id=2>
// The interested message is fixed-length and has no payload.
func NewInterestedMessage() [5]byte {
	return [5]byte{0, 0, 0, 1, 2}
}

// not interested: <len=0001><id=3>
// The not interested message is fixed-length and has no payload.
func NewNotInterestedMessage() [5]byte {
	return [5]byte{0, 0, 0, 1, 3}
}

// have: <len=0005><id=4><piece index>
// The have message is fixed length. The payload is the zero-based index of a piece that has just been successfully downloaded and verified via the hash.
func NewHaveMessage(pieceIndex int32) [9]byte {
	msg := [9]byte{0, 0, 0, 1, 1}
	binary.BigEndian.PutUint32(msg[5:], uint32(pieceIndex))
	return msg
}

func NewHandshakeRequest(metaInfo *MetaInfo, client *Client) [68]byte {
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

	return peerHandshakeRequestMsg
}
