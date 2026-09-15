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

package main

import (
	"crypto/tls"
	"fmt"
	"gotorrent/internal/torrent"
	"io"
	"net/http"
	"os"
	"strconv"
	"sync"
)

// func listenForIncomingRequests(wg *sync.WaitGroup, address string) {
// 	defer wg.Done()
// for {
// 	lisntener, err := net.Listen("tcp", address)
// 	if err != nil {
// 		panic(err)
// 	}
// }
// }

func main() {
	metaInfoFile := os.Args[1]
	metaInfo, err := torrent.CreateMetaInfoFromFile(metaInfoFile)
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup

	fmt.Println(metaInfo)

	port := "6881" // Standard BitTorrent listening port
	// address := ":" + port

	// wg.Add(1)
	// go listenForIncomingRequests(&wg, address)

	portInt, err := strconv.Atoi(port)
	if err != nil {
		panic(err)
	}

	client := torrent.NewClient()

	length, err := metaInfo.GetLength()
	if err != nil {
		panic(err)
	}

	peerId := client.GetPeerId()

	trackerRequest := torrent.CreateNewTrackerRequest().
		SetInfoHash(metaInfo).
		SetPeerId(string(peerId[:])).
		SetPort(int64(portInt)).
		SetUploaded(int64(0)).
		SetDownloaded(int64(0)).
		SetLeft(int64(length)).
		SetCompact(int64(0)).
		SetNoPeerId(int64(0)).
		SetEvent(torrent.STARTED)

	trackerRequestURL, err := trackerRequest.GetURLEncodedRequestString(metaInfo)
	if err != nil {
		panic(err)
	}

	fmt.Println("\n\n......Making GET Request.......\n\n")
	req, err := http.NewRequest("GET", trackerRequestURL, nil)
	req.Header.Set("User-Agent", "BitTorrent/1.0 (gotorrent)")
	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		panic(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Status: %s\n", resp.Status)
	fmt.Printf("Raw Tracker Response:\n%s\n", string(body))

	wg.Wait()

}
