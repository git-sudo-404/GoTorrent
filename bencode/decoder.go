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
	"io"
	"strconv"
)

type decoder struct {
	*bufio.Reader
}

func (d *decoder) decodeInt() int64 {
	d.ReadByte() // consume the 'i'
	numString, err := d.ReadString('e')
	if err != nil {
		panic(err)
	}
	num, err := strconv.ParseInt(numString, 10, 64)
	if err != nil {
		panic(err)
	}
	return num
}

func (d *decoder) decodeString() string {
	strlenString, err := d.ReadString(':')
	if err != nil {
		panic(err)
	}
	strlen, err := strconv.ParseInt(strlenString, 10, 64)
	if err != nil {
		panic(err)
	}

	var str string

	for i := 0; i < int(strlen); i++ {
		chr, err := d.ReadByte()
		if err != nil {
			panic(err)
		}
		str = str + string(chr)
	}
	return str
}

func (d *decoder) decodeList() []any {
	d.ReadByte() // consume the ';'
	peekVal, err := d.Peek(1)
	if err != nil {
		panic(err)
	}
	var list []any
	switch rune(peekVal[0]) {
	case 'i':
		list = append(list, d.decodeInt())
	case 'l':
		list = append(list, d.decodeList()...)
	case 'd':
		list = append(list, d.decodeDict())
	default:
		list = append(list, d.decodeString())
	}
	return list
}

func (d *decoder) decodeDict() map[string]any {
	var dict map[string]any
	dict = map[string]any{}
	d.ReadByte() // consume the 'd'
	for {
		key := d.decodeString()
		peekVal, err := d.Peek(1)
		if err != nil {
			panic(err)
		}
		switch rune(peekVal[0]) {
		case 'i':
			dict[key] = d.decodeInt()
		case 'd':
			dict[key] = d.decodeDict()
		case 'l':
			dict[key] = d.decodeList()
		default:
			dict[key] = d.decodeString()
		}
	}
}

//NOTE: we didn't implement the decoder in the exact same way as the encoder using the CreateNewDecoder for a specific purpose.
// -> It made sense to reuse the same bytes.Buffer in encoder , since it was a needed GC Optimization
// -> But then the reusing of bufio.Reader is not safe and might introduce some buggy code due to race conditions if not handled properly and adds additional overhead
// -> So the tradeOff is 4KB reallocation for every decode Vs implementation overhead which has the risk of race conditions
// -> It's obvious choice to the choose the first

func Decode(ior io.Reader) map[string]any {
	br := bufio.NewReader(ior)
	decoder := &decoder{br}
	return decoder.decodeDict()
}
