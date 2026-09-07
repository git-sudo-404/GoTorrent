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

package urlencoder

import "bytes"

type URLEncoder struct {
	bytes.Buffer
}

const hexTable = "0123456789ABCDEF"

func (e *URLEncoder) EncodeString(s string) {
	for i, _ := range s {
		c := s[i]
		if (c >= '0' && c <= '9') ||
			(c >= 'a' && c <= 'z') ||
			(c >= 'A' && c <= 'Z') ||
			c == '-' ||
			c == '_' ||
			c == '.' ||
			c == '~' {
			e.WriteByte(c)
		} else {
			e.WriteByte('%')
			e.WriteByte(hexTable[c>>4])
			e.WriteByte(hexTable[c&0x0f])
		}
	}
}

func (e *URLEncoder) EncodeBytes(b []byte) {
	for _, c := range b {
		if (c >= '0' && c <= '9') ||
			(c >= 'a' && c <= 'z') ||
			(c >= 'A' && c <= 'Z') ||
			c == '-' ||
			c == '_' ||
			c == '.' ||
			c == '~' {
			e.WriteByte(c)
		} else {
			e.WriteByte('%')
			e.WriteByte(hexTable[c>>4])
			e.WriteByte(hexTable[c&0x0f])
		}

	}
}

func NewURLEncoder() *URLEncoder {
	return &URLEncoder{}
}
