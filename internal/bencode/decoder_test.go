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

func TestAppendRawInt(t *testing.T) {
	encoded := "i123e"

	d := &decoder{
		bufio.NewReader(bytes.NewBufferString(encoded)),
	}

	var rawInfo []byte

	err := d.appendRawInt(&rawInfo)
	if err != nil {
		t.Fatal(err)
	}

	want := []byte(encoded)

	if !bytes.Equal(rawInfo, want) {
		t.Errorf("\nGOT  : %q\nWANT : %q", rawInfo, want)
	}
}

func TestAppendRawString(t *testing.T) {
	encoded := "5:hello"

	d := &decoder{
		bufio.NewReader(bytes.NewBufferString(encoded)),
	}

	var rawInfo []byte

	err := d.appendRawString(&rawInfo)
	if err != nil {
		t.Fatal(err)
	}

	want := []byte(encoded)

	if !bytes.Equal(rawInfo, want) {
		t.Errorf("\nGOT  : %q\nWANT : %q", rawInfo, want)
	}
}

func TestAppendRawList(t *testing.T) {
	encoded := "li10e2:hili20eee"

	d := &decoder{
		bufio.NewReader(bytes.NewBufferString(encoded)),
	}

	var rawInfo []byte

	err := d.appendRawList(&rawInfo)
	if err != nil {
		t.Fatal(err)
	}

	want := []byte(encoded)

	if !bytes.Equal(rawInfo, want) {
		t.Errorf("\nGOT  : %q\nWANT : %q", rawInfo, want)
	}
}

func TestAppendRawDict(t *testing.T) {
	encoded := "d4:name4:test6:lengthi10ee"

	d := &decoder{
		bufio.NewReader(bytes.NewBufferString(encoded)),
	}

	var rawInfo []byte

	err := d.appendRawDict(&rawInfo)
	if err != nil {
		t.Fatal(err)
	}

	want := []byte(encoded)

	if !bytes.Equal(rawInfo, want) {
		t.Errorf("\nGOT  : %q\nWANT : %q", rawInfo, want)
	}
}

func TestAppendRawDictNested(t *testing.T) {
	encoded := "d4:name4:test4:infod3:fooi20eee"

	d := &decoder{
		bufio.NewReader(bytes.NewBufferString(encoded)),
	}

	var rawInfo []byte

	err := d.appendRawDict(&rawInfo)
	if err != nil {
		t.Fatal(err)
	}

	want := []byte(encoded)

	if !bytes.Equal(rawInfo, want) {
		t.Errorf("\nGOT  : %q\nWANT : %q", rawInfo, want)
	}
}

func TestGetRawInfoBytes(t *testing.T) {
	encoded := "d8:announce3:foo4:infod4:name4:test6:lengthi10eee"

	d := &decoder{
		bufio.NewReader(bytes.NewBufferString(encoded)),
	}

	got, err := d.getRawInfoBytes()
	if err != nil {
		t.Fatal(err)
	}

	want := []byte("d4:name4:test6:lengthi10ee")

	if !bytes.Equal(got, want) {
		t.Errorf("\nGOT  : %q\nWANT : %q", got, want)
	}
}

func TestGetRawInfoBytesInfoFirst(t *testing.T) {
	encoded := "d4:infod4:name4:test6:lengthi10eee"

	d := &decoder{
		bufio.NewReader(bytes.NewBufferString(encoded)),
	}

	got, err := d.getRawInfoBytes()
	if err != nil {
		t.Fatal(err)
	}

	want := []byte("d4:name4:test6:lengthi10ee")

	if !bytes.Equal(got, want) {
		t.Errorf("\nGOT  : %q\nWANT : %q", got, want)
	}
}
