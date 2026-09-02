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
	"bytes"
	"sort"
	"strconv"
)

type encoder struct {
	bytes.Buffer
}

func (e *encoder) writeInt(num int64) {
	e.WriteByte('i')
	e.WriteString(strconv.FormatInt(num, 10))
	e.WriteByte('e')
}

func (e *encoder) writeUint(num uint64) {
	e.WriteByte('i')
	e.WriteString(strconv.FormatUint(num, 10))
	e.WriteByte('e')
}

func (e *encoder) writeString(str string) {
	e.WriteString(strconv.Itoa(len(str)))
	e.WriteByte(':')
	e.WriteString(str)
}

func (e *encoder) writeList(list []any) {
	e.WriteByte('l')
	for _, item := range list {
		switch v := item.(type) {
		case int64:
			e.writeInt(v)
		case uint64:
			e.writeUint(v)
		case string:
			e.writeString(v)
		case []any:
			e.writeList(v)
		}
	}
	e.WriteByte('e')
}

func (e *encoder) writeDict(dict map[string]any) {
	e.WriteByte('d')

	var keys []string
	for key, _ := range dict {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	for _, key := range keys {
		e.writeString(key)
		switch v := dict[key].(type) {
		case int64:
			e.writeInt(v)
		case uint64:
			e.writeUint(v)
		case string:
			e.writeString(v)
		case []any:
			e.writeList(v)
		case map[string]any:
			e.writeDict(v)
		}
	}

	e.WriteByte('e')
}

func CreateNewEncoder() *encoder {
	return &encoder{
		bytes.Buffer{},
	}
}

func (e *encoder) getEncodedString() string {
	return e.String()
}
