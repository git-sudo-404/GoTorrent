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
	"fmt"
	"gotorrent/internal/bencode"
	"net/http"
	"time"
)

func populateInitialTrackerRequestParams(tr *TrackerRequest, mi *MetaInfo, client *Client) {
	fmt.Println("[LOG] Populating TrackerRequest...")
	tr.SetInfoHash(mi)
	tr.SetPeerId(string(client.peerId[:]))
	tr.SetPort(int64(client.port))
	tr.SetUploaded(int64(0))
	tr.SetDownloaded(int64(0))
	tr.SetLeft(mi.length)
	tr.SetCompact(1)
	tr.SetNoPeerId(int64(1))
	tr.SetEvent(STARTED)
}

func sendTrackerRequest(trackerRequestURL *string) *TrackerResponse {

	trackerResponse := NewTrackerResponse()

	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	req, err := http.NewRequest(http.MethodGet, *trackerRequestURL, nil)
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println("[LOG] Sending TrackerRequest...")
	fmt.Println("TRACKER REQUEST URL : ", *trackerRequestURL)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error tracker server return : ", resp.Body)
	}

	fmt.Println("STATUS:", resp.Status)

	buf := bufio.NewReader(resp.Body)

	decodedResponse, err := bencode.Decode(buf)
	if err != nil {
		panic(err)
	}

	fmt.Println(decodedResponse)

	return trackerResponse
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

	fmt.Printf("INFO RAW BYTES : % x\n", metaInfo.GetRawInfo())
	fmt.Printf("INFO HASH string: %q\n", metaInfo.GetInfoHash())
	fmt.Printf("INFO HASH bytes: % x\n", trackerRequest.infoHash)
	fmt.Printf("INFO HASH string: %q\n", trackerRequest.infoHash)
	fmt.Println("PORT:", client.port)

	trackerRequestURL, err := trackerRequest.GetURLEncodedRequestString(metaInfo)
	if err != nil {
		panic(err)
	}

	sendTrackerRequest(&trackerRequestURL)

}
