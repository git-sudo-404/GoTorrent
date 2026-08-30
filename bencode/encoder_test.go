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

package bencode

import (
	"testing"
)

func TestEncodeInteger(t *testing.T) {
	got := encodeInteger(13)
	want := []byte("i13e")
	if string(got) != string(want) {
		t.Errorf("Got %s | Want %s", got, want)
	} else {
		t.Log("[PASS] : encodeInteger()")
	}
}

func TestEncodeString(t *testing.T) {
	got := encodeString("hi")
	want := []byte("2:hi")
	if string(got) != string(want) {
		t.Errorf("Got %s | Want %s", got, want)
	} else {
		t.Log("[PASS] : encodeString()")
	}
}

func TestEncodeList(t *testing.T) {
	got := encodeList([]any{13, "hi"})
	want := "li13e2:hie"
	if string(got) != string(want) {
		t.Errorf("Got %s | Want %s", got, want)
	} else {
		t.Log("[PASS] : encodeList()")
	}
}

func TestEncodeDict(t *testing.T) {
	got := encodeDict(map[string]any{
		"key1": 13,
		"key2": "hi",
		"key3": []any{12, "hi"},
		"key4": map[string]any{
			"k5": 1,
		},
	})
	want := []byte("d4:key1i13e4:key22:hi4:key3li12e2:hie4:key4d2:k5i1eee")
	if string(got) != string(want) {
		t.Errorf("Got %s | Want %s", got, want)
	} else {
		t.Log("[PASS] : encodeDict()")
	}
}
