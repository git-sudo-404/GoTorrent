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

import "net"

type Peer struct { // no peerId when in compact mode = 1
	ip   net.IP // peer's IP address either IPv6 (hexed) or IPv4 (dotted quad) or DNS name (string)
	port int64  //  peer's port number (integer)
}

type TrackerResponse struct {
	interval    int64  // Interval in seconds that the client should wait between sending regular requests to the tracker
	minInterval *int64 // (optional) Minimum announce interval. If present clients must not reannounce more frequently than this.
	trackerId   string //  A string that the client should send back on its next announcements. If absent and a previous announce sent a tracker id, do not discard the old value; keep using it.
	complete    int64  //  number of peers with the entire file, i.e. seeders (integer)
	incomplete  int64  // number of non-seeder peers, aka "leechers" (integer)
	peers       []Peer
}

func (tr *TrackerResponse) SetPeers(peers []Peer) *TrackerResponse {
	tr.peers = peers
	return tr
}

func (tr *TrackerResponse) SetInterval(interval int64) *TrackerResponse {
	tr.interval = interval
	return tr
}

func (tr *TrackerResponse) SetMinInterval(minInterval int64) *TrackerResponse {
	tr.minInterval = &minInterval
	return tr
}

func (tr *TrackerResponse) SetTrackerId(trackerId string) *TrackerResponse {
	tr.trackerId = trackerId
	return tr
}

func (tr *TrackerResponse) SetComplete(complete int64) *TrackerResponse {
	tr.complete = complete
	return tr
}

func (tr *TrackerResponse) SetIncomplete(incomplete int64) *TrackerResponse {
	tr.incomplete = incomplete
	return tr
}
