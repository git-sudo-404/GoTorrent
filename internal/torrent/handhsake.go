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

type PeerHandshakeRequest struct {
	pstr     []byte   //  identifier of the protocol
	reserved [8]byte  // eight (8) reserved bytes. All current implementations use all zeroes. Each bit in these bytes can be used to change the behavior of the protocol. An email from Bram suggests that trailing bits should be used first, so that leading bits may be used to change the meaning of trailing bits.
	infoHash [20]byte // 20-byte SHA1 hash of the info key in the metainfo file. This is the same info_hash that is transmitted in tracker requests.
	peerId   [20]byte // 20-byte string used as a unique ID for the client. This is usually the same peer_id that is transmitted in tracker requests (but not always e.g. an anonymity option in Azureus).
}

func (phr *PeerHandshakeRequest) SetPstr(pstr []byte) *PeerHandshakeRequest {
	phr.pstr = pstr
	return phr
}

func (phr *PeerHandshakeRequest) SetReserved(reserved [8]byte) *PeerHandshakeRequest {
	phr.reserved = reserved
	return phr
}

func (phr *PeerHandshakeRequest) SetInfoHash(infoHash [20]byte) *PeerHandshakeRequest {
	phr.infoHash = infoHash
	return phr
}

func (phr *PeerHandshakeRequest) SetPeerId(peerId [20]byte) *PeerHandshakeRequest {
	phr.peerId = peerId
	return phr
}

func (phr *PeerHandshakeRequest) Serialize() []byte {
	buf := make([]byte, 49+len(phr.pstr))
	buf[0] = byte(len(phr.pstr))
	copy(buf[1:], phr.pstr)
	copy(buf[1+len(phr.pstr):], phr.reserved[:])
	copy(buf[1+len(phr.pstr)+8:], phr.infoHash[:])
	copy(buf[1+len(phr.pstr)+28:], phr.peerId[:])
	return buf
}
