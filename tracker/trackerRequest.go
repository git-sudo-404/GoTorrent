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
	donwloaded int64
	left       int64
	compact    int64 // 1 or 0
	noPeerId   int64
	event      *Event
	ip         *string
	numwant    *int64
	key        *string
	trackerId  *string
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
	tr.donwloaded = downloaded
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
	URLString.WriteString(fmt.Sprintf("?info_hash=%s", tr.infoHash))
	URLString.WriteString(fmt.Sprintf("?peer_id=%s", tr.peerId))
	URLString.WriteString(fmt.Sprintf("?port=%d", tr.port))
	URLString.WriteString(fmt.Sprintf("?uploaded=%d", tr.uploaded))
	URLString.WriteString(fmt.Sprintf("?donwloaded=%d", tr.donwloaded))
	URLString.WriteString(fmt.Sprintf("?left=%d", tr.left))
	URLString.WriteString(fmt.Sprintf("?compact=%d", tr.compact))
	URLString.WriteString(fmt.Sprintf("?no_peer_id=%d", tr.noPeerId))
	if tr.event != nil {
		URLString.WriteString(fmt.Sprintf("?event=%s", *tr.event))
	}
	if tr.ip != nil {
		URLString.WriteString(fmt.Sprintf("?ip=%s", *tr.ip))
	}
	if tr.numwant != nil {
		URLString.WriteString(fmt.Sprintf("?numwant=%d", *tr.numwant))
	}
	if tr.key != nil {
		URLString.WriteString(fmt.Sprintf("?key=%s", *tr.key))
	}
	if tr.trackerId != nil {
		URLString.WriteString(fmt.Sprintf("?trackerid=%s", *tr.trackerId))
	}

	urlEncoder := urlencoder.NewURLEncoder()
	urlEncoder.EncodeString(URLString.String())

	return urlEncoder.String(), nil
}
