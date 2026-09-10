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
	"encoding/binary"
	"fmt"
	"gotorrent/internal/bencode"
	"io"
	"net"
)

func parsePeersCompactString(peersCompactString string) ([]Peer, error) {
	// first 4-bytes -> IP
	// next 2-bytes -> port

	data := []byte(peersCompactString)
	dataLen := len(data)
	if dataLen%6 != 0 {
		return nil, fmt.Errorf("The peers compact string length must be a multiple of 6 , but got of length : %d", dataLen)
	}

	numberOfPeers := dataLen / 6
	peers := make([]Peer, numberOfPeers)

	for i := 0; i < numberOfPeers; i++ {
		offset := i * 6

		ip := make(net.IP, 4)
		copy(ip, data[offset:offset+4])

		port := int64(binary.BigEndian.Uint16(data[offset+4 : offset+6]))

		peers[i] = Peer{
			ip:   ip,
			port: port,
		}
	}

	return peers, nil
}

func ParseTrackerResponse(r io.Reader) (*TrackerResponse, error) {

	trackerResponse := &TrackerResponse{}

	decodedResponseDict, err := bencode.Decode(r)
	if err != nil {
		return nil, err
	}
	if _, exists := decodedResponseDict["failure reason"]; exists {
		return nil, fmt.Errorf("Tracker Response Failed")
	}

	interval, err := checkFieldPresenceAndType[int64](decodedResponseDict, "interval")
	if err != nil {
		return nil, err
	}
	trackerResponse.SetInterval(interval)

	minInterval, minIntervalPresent, err := checkOptionalFieldPresenceAndType[int64](decodedResponseDict, "min interval")
	if err != nil {
		return nil, err
	}
	if minIntervalPresent {
		trackerResponse.SetMinInterval(minInterval)
	}

	trackerId, err := checkFieldPresenceAndType[string](decodedResponseDict, "tracker id")
	if err != nil {
		return nil, err
	}
	trackerResponse.SetTrackerId(trackerId)

	complete, err := checkFieldPresenceAndType[int64](decodedResponseDict, "complete")
	if err != nil {
		return nil, err
	}
	trackerResponse.SetComplete(complete)

	incomplete, err := checkFieldPresenceAndType[int64](decodedResponseDict, "incomplete")
	if err != nil {
		return nil, err
	}
	trackerResponse.SetIncomplete(incomplete)

	peersCompactString, err := checkFieldPresenceAndType[string](decodedResponseDict, "peers")
	if err != nil {
		return nil, err
	}

	peers, err := parsePeersCompactString(peersCompactString)
	if err != nil {
		return nil, err
	}
	trackerResponse.SetPeers(peers)

	return trackerResponse, nil
}
