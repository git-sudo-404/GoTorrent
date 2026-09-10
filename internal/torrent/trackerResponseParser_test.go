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
	"bytes"
	"encoding/binary"
	"fmt"
	"gotorrent/internal/bencode"
	"net"
	"testing"
)

func TestParseTrackerResponse(t *testing.T) {

	trackerResponseDict := map[string]any{}
	trackerResponseDict["interval"] = int64(10)
	trackerResponseDict["min interval"] = int64(5)
	trackerResponseDict["tracker id"] = "tracker_id_1"
	trackerResponseDict["complete"] = int64(2)
	trackerResponseDict["incomplete"] = int64(3)

	mockPeerIps := new(bytes.Buffer)
	mockPeerIps.Write(net.ParseIP("192.2.1.11").To4())
	var port uint16 = 8001
	binary.Write(mockPeerIps, binary.BigEndian, port)

	mockPeerIps.Write(net.ParseIP("192.2.1.12").To4())
	port = 8001
	binary.Write(mockPeerIps, binary.BigEndian, port)

	mockPeerIps.Write(net.ParseIP("192.2.1.13").To4())
	port = 8001
	binary.Write(mockPeerIps, binary.BigEndian, port)

	mockPeerIps.Write(net.ParseIP("192.2.1.14").To4())
	port = 8001
	binary.Write(mockPeerIps, binary.BigEndian, port)

	mockPeerIps.Write(net.ParseIP("192.2.1.15").To4())
	port = 8001
	binary.Write(mockPeerIps, binary.BigEndian, port)

	mockPeerBytes := mockPeerIps.String()

	trackerResponseDict["peers"] = mockPeerBytes

	encoder := bencode.CreateNewEncoder()
	encoder.Encode(trackerResponseDict)
	bencodedTrackerResponse := encoder.String()

	mockTrackerResponse := new(bytes.Buffer)
	mockTrackerResponse.WriteString(bencodedTrackerResponse)

	tr, err := ParseTrackerResponse(mockTrackerResponse)
	if err != nil {
		panic(err)
	}

	if tr.interval != int64(10) {
		t.Errorf("\nGOT  : %d\nWANT : %d", tr.interval, 10)
	}
	if *tr.minInterval != int64(5) {
		t.Errorf("\nGOT  : %d\nWANT : %d", *tr.minInterval, 5)
	}
	if tr.trackerId != "tracker_id_1" {
		t.Errorf("\nGOT  : %s\nWANT : tracker_id_1", tr.trackerId)
	}
	if tr.complete != int64(2) {
		t.Errorf("\nGOT  : %d\nWANT : %d", tr.complete, 2)
	}
	if tr.incomplete != int64(3) {
		t.Errorf("\nGOT  : %d\nWANT : %d", tr.incomplete, 3)
	}

	for i := 1; i <= 5; i++ {
		wantedIp := net.ParseIP(fmt.Sprintf("192.2.1.1%d", i))
		if !tr.peers[i-1].ip.Equal(wantedIp) {
			t.Errorf("expected IP %s, got %s", wantedIp, tr.peers[i-1].ip)
		}
		if tr.peers[i-1].port != int64(8001) {
			t.Errorf("Port mismatch in the tr.peers , wanted 8001")
		}
	}

}
