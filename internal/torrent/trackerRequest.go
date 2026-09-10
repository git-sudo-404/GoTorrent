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
	"crypto/sha1"
	"fmt"
	bencode "gotorrent/internal/bencode"
	urlencoder "gotorrent/internal/url-encoder"
	"strings"
)

type Event string

const (
	STARTED   Event = "started"
	STOPPED   Event = "stopped"
	COMPLETED Event = "completed"
)

type TrackerRequest struct {
	infoHash   string
	peerId     string
	port       int64
	uploaded   int64
	downloaded int64
	left       int64
	compact    int64 // 1 or 0
	noPeerId   int64
	event      *Event
}

func CreateNewTrackerRequest() *TrackerRequest {
	return &TrackerRequest{
		infoHash:   "",
		peerId:     "",
		port:       0,
		uploaded:   0,
		downloaded: 0,
		left:       0,
		compact:    0,
		noPeerId:   0,
	}
}

func (tr *TrackerRequest) SetInfoHash(mi *MetaInfo) *TrackerRequest {
	infoDict, err := mi.GetInfoDict()
	if err != nil {
		return nil
	}

	encoder := bencode.CreateNewEncoder()
	encoder.Encode(infoDict)
	infoBencoded := encoder.Bytes()

	infoHashed := sha1.Sum(infoBencoded)

	urlencoder := urlencoder.NewURLEncoder()
	urlencoder.EncodeBytes(infoHashed[:])

	infoURLEncoded := urlencoder.String()
	tr.infoHash = string(infoURLEncoded)

	return tr
}

// NOTE: The peerId has to be constructed only once during the client startup
// which will be done and be stored in the client struct
func (tr *TrackerRequest) SetPeerId(peerId string) *TrackerRequest {
	urlencoder := urlencoder.NewURLEncoder()
	urlencoder.EncodeString(peerId)
	tr.peerId = urlencoder.String()
	return tr
}

// port got from client
func (tr *TrackerRequest) SetPort(port int64) *TrackerRequest {
	tr.port = port
	return tr
}

func (tr *TrackerRequest) SetUploaded(uploaded int64) *TrackerRequest {
	tr.uploaded = uploaded
	return tr
}

func (tr *TrackerRequest) SetDownloaded(downloaded int64) *TrackerRequest {
	tr.downloaded = downloaded
	return tr
}

func (tr *TrackerRequest) SetLeft(left int64) *TrackerRequest {
	tr.left = left
	return tr
}

func (tr *TrackerRequest) SetCompact(compact int64) *TrackerRequest {
	tr.compact = compact
	return tr
}

func (tr *TrackerRequest) SetNoPeerId(noPeerId int64) *TrackerRequest {
	tr.noPeerId = noPeerId
	return tr
}

func (tr *TrackerRequest) SetEvent(event Event) *TrackerRequest {
	tr.event = &event
	return tr
}

func (tr *TrackerRequest) GetURLEncodedRequestString(mi *MetaInfo) (string, error) {
	var URLString strings.Builder

	baseURL, err := mi.GetAnnounce()
	if err != nil {
		return "", nil
	}

	URLString.WriteString(baseURL)
	fmt.Fprintf(&URLString, "?info_hash=%s", tr.infoHash)
	fmt.Fprintf(&URLString, "&peer_id=%s", tr.peerId)
	fmt.Fprintf(&URLString, "&port=%d", tr.port)
	fmt.Fprintf(&URLString, "&uploaded=%d", tr.uploaded)
	fmt.Fprintf(&URLString, "&downloaded=%d", tr.downloaded)
	fmt.Fprintf(&URLString, "&left=%d", tr.left)
	fmt.Fprintf(&URLString, "&compact=%d", tr.compact)
	fmt.Fprintf(&URLString, "&no_peer_id=%d", tr.noPeerId)
	if tr.event != nil {
		fmt.Fprintf(&URLString, "&event=%s", *tr.event)
	}

	return URLString.String(), nil
}
