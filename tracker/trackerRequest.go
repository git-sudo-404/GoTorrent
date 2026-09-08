/*
 * Copyright (c) 2026 git-sudo-404 <https://github.com/git-sudo-404/GoTorrent.git>
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 * THE SOFTWARE.
 */

package tracker

import (
	"crypto/sha1"
	"fmt"
	"gotorrent/bencode"
	metainfo "gotorrent/meta-info"
	urlencoder "gotorrent/url-encoder"
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
	ip         *string
	numwant    *int64
	key        *string
	trackerId  *string
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

func (tr *TrackerRequest) SetInfoHash(mi *metainfo.MetaInfo) *TrackerRequest {
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
	tr.peerId = peerId
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

func (tr *TrackerRequest) SetIp(ip string) *TrackerRequest {
	tr.ip = &ip
	return tr
}

func (tr *TrackerRequest) SetNumwant(numwant int64) *TrackerRequest {
	tr.numwant = &numwant
	return tr
}

func (tr *TrackerRequest) SetKey(key string) *TrackerRequest {
	tr.key = &key
	return tr
}

func (tr *TrackerRequest) SetTrackerId(trackerId string) *TrackerRequest {
	tr.trackerId = &trackerId
	return tr
}

func (tr *TrackerRequest) GetURLEncodedRequestString(mi *metainfo.MetaInfo) (string, error) {
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
	if tr.ip != nil {
		fmt.Fprintf(&URLString, "&ip=%s", *tr.ip)
	}
	if tr.numwant != nil {
		fmt.Fprintf(&URLString, "&numwant=%d", *tr.numwant)
	}
	if tr.key != nil {
		fmt.Fprintf(&URLString, "&key=%s", *tr.key)
	}
	if tr.trackerId != nil {
		fmt.Fprintf(&URLString, "&trackerid=%s", *tr.trackerId)
	}

	urlEncoder := urlencoder.NewURLEncoder()
	urlEncoder.EncodeString(URLString.String())

	return urlEncoder.String(), nil
}
