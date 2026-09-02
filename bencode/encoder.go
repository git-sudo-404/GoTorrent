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
	"fmt"
	"sort"
	"strconv"
)

type Encoder struct {
	bytes.Buffer
}

func (e *Encoder) encodeInt(num int64) {
	e.WriteByte('i')
	e.WriteString(strconv.FormatInt(num, 10))
	e.WriteByte('e')
}

func (e *Encoder) encodeUint(num uint64) {
	e.WriteByte('i')
	e.WriteString(strconv.FormatUint(num, 10))
	e.WriteByte('e')
}

func (e *Encoder) encodeString(str string) {
	e.WriteString(strconv.Itoa(len(str)))
	e.WriteByte(':')
	e.WriteString(str)
}

func (e *Encoder) encodeList(list []any) {
	e.WriteByte('l')
	for _, item := range list {
		switch v := item.(type) {
		case int64:
			e.encodeInt(v)
		case uint64:
			e.encodeUint(v)
		case string:
			e.encodeString(v)
		case []any:
			e.encodeList(v)
		}
	}
	e.WriteByte('e')
}

func (e *Encoder) encodeDict(dict map[string]any) {
	e.WriteByte('d')

	var keys []string
	for key, _ := range dict {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	for _, key := range keys {
		e.encodeString(key)
		switch v := dict[key].(type) {
		case int64:
			e.encodeInt(v)
		case uint64:
			e.encodeUint(v)
		case string:
			e.encodeString(v)
		case []any:
			e.encodeList(v)
		case map[string]any:
			e.encodeDict(v)
		}
	}

	e.WriteByte('e')
}

func CreateNewEncoder() *Encoder {
	return &Encoder{
		bytes.Buffer{},
	}
}

func (e *Encoder) getEncodedString() string {
	return e.String()
}

func (e *Encoder) getEncodedBytes() []byte {
	return e.Bytes()
}

func (e *Encoder) Encode(data any) error {
	switch v := data.(type) {
	case int64:
		e.encodeInt(v)
	case uint64:
		e.encodeUint(v)
	case string:
		e.encodeString(v)
	case []any:
		e.encodeList(v)
	case map[string]any:
		e.encodeDict(v)
	default:
		return fmt.Errorf("bencode: unsupported data type: %T", data)
	}
	return nil
}
