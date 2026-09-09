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

package torrent

import (
	"crypto/sha1"
	"fmt"
	"gotorrent/bencode"
	urlencoder "gotorrent/url-encoder"
	"testing"
)

func TestGetURLEncodedRequestString(t *testing.T) {
	metaInfo, err := CreateTestMetaInfo()
	if err != nil {
		panic(err)
	}
	client := NewClient()
	peerIdBytes := client.GetPeerId()
	peerId := string(peerIdBytes[:])
	port := 8080
	uploaded := 1024
	downloaded := 1024
	left := 1024
	compact := 1
	noPeerId := 1
	event := STARTED

	trackerRequest := CreateNewTrackerRequest().
		SetInfoHash(metaInfo).
		SetPeerId(peerId).
		SetPort(int64(port)).
		SetUploaded(int64(uploaded)).
		SetDownloaded(int64(downloaded)).
		SetLeft(int64(left)).
		SetCompact(int64(compact)).
		SetNoPeerId(int64(noPeerId)).
		SetEvent(event)

	got, _ := trackerRequest.GetURLEncodedRequestString(metaInfo)

	baseURL, _ := metaInfo.GetAnnounce()

	infoDict, err := metaInfo.GetInfoDict()
	if err != nil {
		panic(err)
	}

	// enode the info dict
	encoder := bencode.CreateNewEncoder()
	encoder.Encode(infoDict)
	infoBencoded := encoder.Bytes()

	infoHashed := sha1.Sum(infoBencoded)

	urlencoder := urlencoder.NewURLEncoder()
	urlencoder.EncodeBytes(infoHashed[:])

	infoURLEncoded := urlencoder.String()
	infoHash := string(infoURLEncoded)

	urlencoder.Reset()

	// urlEncode the peerId
	urlencoder.EncodeString(peerId)
	peerId = urlencoder.String()

	want := fmt.Sprintf("%s?info_hash=%s&peer_id=%s&port=%d&uploaded=%d&downloaded=%d&left=%d&compact=%d&no_peer_id=%d&event=%s",
		baseURL, infoHash, peerId, port, uploaded, downloaded, left, compact, noPeerId, event)
	if got != want {
		t.Errorf("\nGOT  : %v\nwant : %v", got, want)
	}
}
