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
	"bufio"
	"bytes"
	"reflect"
	"testing"
)

func TestDecodeInt(t *testing.T) {
	encodedBuffer := bytes.NewBufferString("i13e")
	br := bufio.NewReader(encodedBuffer)
	decoder := &decoder{br}

	got, err := decoder.decodeInt()
	want := int64(13)

	if err != nil {
		panic(err)
	}

	switch v := got.(type) {
	case int64:
		if v != want {
			t.Errorf("\nv  : %d\nWANT : %d", v, want)
		}
	case uint64:
		if v != uint64(want) {
			t.Errorf("\nv  : %d\nWANT : %d", v, want)
		}
	default:
		t.Errorf("Integer decoded into a type that's not desireable")
	}
}

func TestDecodeString(t *testing.T) {
	encodedBuffer := bytes.NewBufferString("5:hello")
	br := bufio.NewReader(encodedBuffer)
	decoder := &decoder{br}
	got, err := decoder.decodeString()
	if err != nil {
		panic(err)
	}
	want := "hello"

	if got != want {
		t.Errorf("\nv  : %s\nWANT : %s", got, want)
	}
}

func TestDecodeList(t *testing.T) {
	encodedBuffer := bytes.NewBufferString("li10e2:hili20eed4:key16:value1ee")
	br := bufio.NewReader(encodedBuffer)
	decoder := &decoder{br}

	got, err := decoder.decodeList()
	if err != nil {
		panic(err)
	}
	want := []any{int64(10), "hi", []any{int64(20)}, map[string]any{"key1": "value1"}}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("\nGOT  : %#v\nWANT : %#v", got, want)
	}
}

func TestDecodeDict(t *testing.T) {
	encodedBuffer := bytes.NewBufferString("d4:key1i10e4:key2i20e4:key35:hello4:key4li10e2:hili20eee4:key5d4:key6li30eeee")
	br := bufio.NewReader(encodedBuffer)
	decoder := &decoder{br}

	got, err := decoder.decodeDict()
	if err != nil {
		panic(err)
	}

	var want map[string]any
	want = map[string]any{}
	want["key1"] = int64(10)
	want["key2"] = int64(20)
	want["key3"] = "hello"
	want["key4"] = []any{int64(10), "hi", []any{int64(20)}}
	want["key5"] = map[string]any{"key6": []any{int64(30)}}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("\nGOT  : %#v\nWANT : %#v", got, want)
	}

}
