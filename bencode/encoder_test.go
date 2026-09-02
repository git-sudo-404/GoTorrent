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

import "testing"

func TestWriteInt(t *testing.T) {
	e := CreateNewEncoder()
	var num int64
	num = 10
	e.encodeInt(num)

	got := e.getEncodedString()
	want := "i10e"

	if string(got) != string(want) {
		t.Errorf("\bGOT  : %s\nWANT : %s", got, want)
	}
}

func TestWriteUint(t *testing.T) {
	e := CreateNewEncoder()
	var num uint64
	num = 10
	e.encodeUint(num)

	got := e.getEncodedString()
	want := "i10e"

	if string(got) != string(want) {
		t.Errorf("\nGOT  : %s\nWANT : %s", got, want)
	}
}

func TestWriteString(t *testing.T) {
	e := CreateNewEncoder()
	var str string
	str = "hi"
	e.encodeString(str)

	got := e.getEncodedString()
	want := "2:hi"

	if string(got) != string(want) {
		t.Errorf("\nGOT  : %s\nWANT : %s", got, want)
	}
}

func TestWriteList(t *testing.T) {
	e := CreateNewEncoder()
	list := []any{int64(10), "hi", []any{uint64(20)}}
	e.encodeList(list)

	got := e.getEncodedString()
	want := "li10e2:hili20eee"

	if string(got) != string(want) {
		t.Errorf("\nGOT  : %s\nWANT : %s", got, want)
	}
}

func TestWriteDict(t *testing.T) {
	e := CreateNewEncoder()
	var dict map[string]any
	dict = map[string]any{}
	dict["key1"] = uint64(10)
	dict["key2"] = int64(20)
	dict["key3"] = "hello"
	dict["key4"] = []any{int64(10), "hi", []any{uint64(20)}}
	dict["key5"] = map[string]any{"key6": []any{int64(30)}}
	e.encodeDict(dict)

	got := e.getEncodedString()
	want := "d4:key1i10e4:key2i20e4:key35:hello4:key4li10e2:hili20eee4:key5d4:key6li30eeee"

	if string(got) != string(want) {
		t.Errorf("\nGOT  : %s\nWANT : %s", got, want)
	}
}
