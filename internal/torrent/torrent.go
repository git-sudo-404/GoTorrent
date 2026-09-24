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
	"net"
	"net/http"
	"time"
)

func populateInitialTrackerRequestParams(tr *TrackerRequest, mi *MetaInfo, client *Client) {
	fmt.Println("[LOG] Populating TrackerRequest...")
	tr.SetInfoHash(mi)
	tr.SetPeerId(string(client.clientId[:]))
	tr.SetPort(int64(client.port))
	tr.SetUploaded(int64(0))
	tr.SetDownloaded(int64(0))
	tr.SetLeft(mi.length)
	tr.SetCompact(int64(0))
	tr.SetNoPeerId(int64(1))
	tr.SetEvent(STARTED)
}

func sendTrackerRequest(trackerRequestURL *string, client *Client) {

	httpClient := &http.Client{
		Timeout: 5 * time.Second,
	}
	req, err := http.NewRequest(http.MethodGet, *trackerRequestURL, nil)
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

	// fmt.Println("[LOG] Tracker Responded with the following : ....")
	// fmt.Println("[LOG] interval : ", interval)
	// fmt.Println("[LOG] complete : ", complete)
	// fmt.Println("[LOG] incomplete : ", incomplete)
	// fmt.Println("[LOG] peers : ", peers)

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

	if peerList, ok := peers.([]any); ok {
		// fmt.Println("[DEBUG PRINT] PEER LIST : ", peerList)
		for _, peerListItem := range peerList {
			// fmt.Println("[DEBUG PRINT] PEERLISTITEM : ", peerListItem)
			if peerMap, ok := peerListItem.(map[string]any); ok {
				ip, ok := peerMap["ip"].(string)
				// fmt.Println("[DEBUG PRINT] IP : ", peerMap["ip"])
				// fmt.Println("[DEBUG PRINT] PORT : ", peerMap["port"])
				if !ok {
					fmt.Println("[ERROR] ip not of type string from tracker response")
				}
				port, ok := peerMap["port"].(int64)
				if !ok {
					fmt.Println("[ERROR] port not of type integer from tracker response")
				}
				client.peers = append(client.peers, Peer{
					ip:   net.ParseIP(ip),
					port: port,
				})
			}
		}
	} else {
		fmt.Println("[ERROR] Wrong data type in peerList")
	}

	fmt.Println("[LOG] Updated Client Peers from tracker Response")

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

	trackerRequestURL, err := trackerRequest.GetURLEncodedRequestString(metaInfo)
	if err != nil {
		panic(err)
	}

	sendTrackerRequest(&trackerRequestURL, client)

	// fmt.Println(len(client.peers))
	// for _, peer := range client.peers {
	// 	fmt.Println(peer)
	// }

}
