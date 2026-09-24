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
	"io"
	"strconv"
)

type decoder struct {
	*bufio.Reader
}

func (d *decoder) decodeInt() (any, error) {

	d.ReadByte() // consume the 'i'
	numStr, err := d.ReadString('e')
	if err != nil {
		panic(err)
	}
	numStr = numStr[:len(numStr)-1]
	if value, err := strconv.ParseInt(numStr, 10, 64); err == nil {
		return value, nil
	} else if value, err := strconv.ParseUint(numStr, 10, 64); err == nil {
		return value, nil
	} else {
		return nil, err
	}
}

func (d *decoder) decodeString() (string, error) {
	lenNumStr, err := d.ReadString(':')
	if err != nil {
		panic(err)
	}
	lenNumStr = lenNumStr[:len(lenNumStr)-1]
	lenNum, err := strconv.ParseInt(lenNumStr, 10, 64)
	if err != nil {
		panic(err)
	}

	buf := make([]byte, lenNum)
	if _, err := io.ReadFull(d, buf); err != nil {
		panic(err)
	}
	return string(buf), nil
}

func (d *decoder) decodeList() ([]any, error) {
	var result []any
	d.ReadByte() // consume the 'l'

	for {

		peekBytes, err := d.Peek(1)
		if err != nil {
			panic(err)
		}

		if rune(peekBytes[0]) == 'e' {
			d.ReadByte()
			break
		}

		switch rune(peekBytes[0]) {
		case 'i':
			decodedInt, err := d.decodeInt()
			if err != nil {
				panic(err)
			}
			result = append(result, decodedInt)
		case 'l':
			decodedList, err := d.decodeList()
			if err != nil {
				panic(err)
			}
			result = append(result, decodedList)
		case 'd':
			decodedDict, err := d.decodeDict()
			if err != nil {
				panic(err)
			}
			result = append(result, decodedDict)
		default:
			decodedString, err := d.decodeString()
			if err != nil {
				panic(err)
			}
			result = append(result, decodedString)
		}
	}

	return result, nil
}

func (d *decoder) decodeDict() (map[string]any, error) {
	d.ReadByte() // consume the 'd'
	var dict map[string]any
	dict = map[string]any{}

	for {

		peekBytes, err := d.Peek(1)
		if err != nil {
			panic(err)
		}

		if rune(peekBytes[0]) == 'e' {
			d.ReadByte()
			break
		}

		key, err := d.decodeString()
		if err != nil {
			panic(err)
		}

		peekBytes, err = d.Peek(1)
		if err != nil {
			panic(err)
		}

		switch rune(peekBytes[0]) {
		case 'i':
			value, err := d.decodeInt()
			if err != nil {
				panic(err)
			}
			dict[key] = value
		case 'l':
			value, err := d.decodeList()
			if err != nil {
				panic(err)
			}
			dict[key] = value
		case 'd':
			value, err := d.decodeDict()
			if err != nil {
				panic(err)
			}
			dict[key] = value
		default:
			value, err := d.decodeString()
			if err != nil {
				panic(err)
			}
			dict[key] = value
		}

	}
	return dict, nil
}

//NOTE: we didn't implement the decoder in the exact same way as the encoder using the CreateNewDecoder for a specific purpose.
// -> It made sense to reuse the same bytes.Buffer in encoder , since it was a needed GC Optimization
// -> But then the reusing of bufio.Reader is not safe and might introduce some buggy code due to race conditions if not handled properly and adds additional overhead
// -> So the tradeOff is 4KB reallocation for every decode Vs implementation overhead which has the risk of race conditions
// -> It's obvious choice to the choose the first

func Decode(ior io.Reader) (map[string]any, error) {
	br := bufio.NewReader(ior)
	decoder := &decoder{br}
	return decoder.decodeDict()
}

//NOTE: The below code does not actually belong to bencode, but then it is needed to get the rawInfo from the meta-info file which will be hashed and sent along the tracker request.
// Whenever the below functions are used, they may destroy the bytes Buffer of the decoder .
// So, in the application layer use a seperate decoder structs for these purposes

func (d *decoder) appendRawInt(rawInfo *[]byte) error {
	buf, err := d.ReadBytes('e')
	if err != nil {
		return err
	}
	*rawInfo = append(*rawInfo, buf...)
	return nil
}

func (d *decoder) appendRawString(rawInfo *[]byte) error {
	stringLenBytes, err := d.ReadBytes(':')
	if err != nil {
		return err
	}
	*rawInfo = append(*rawInfo, stringLenBytes...)
	stringLen, err := strconv.ParseInt(string(stringLenBytes[:len(stringLenBytes)-1]), 10, 64)
	if err != nil {
		return err
	}
	buf := make([]byte, stringLen)
	if _, err := io.ReadFull(d, buf); err != nil {
		return err
	}
	*rawInfo = append(*rawInfo, buf...)
	return nil
}

func (d *decoder) appendRawList(rawInfo *[]byte) error {
	lbyte, err := d.ReadByte()
	*rawInfo = append(*rawInfo, lbyte)
	if err != nil {
		return err
	}
	for {

		peakByte, err := d.Peek(1)
		if err != nil {
			return err
		}

		if rune(peakByte[0]) == 'e' {
			lastByte, err := d.ReadByte()
			if err != nil {
				return err
			}
			*rawInfo = append(*rawInfo, lastByte)
			break
		}

		switch rune(peakByte[0]) {
		case 'i':
			err := d.appendRawInt(rawInfo)
			if err != nil {
				return err
			}
		case 'l':
			err := d.appendRawList(rawInfo)
			if err != nil {
				return err
			}
		case 'd':
			if err := d.appendRawDict(rawInfo); err != nil {
				return err
			}
		default: // string
			err := d.appendRawString(rawInfo)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (d *decoder) appendRawDict(rawInfo *[]byte) error {
	dByte, err := d.ReadByte()
	if err != nil {
		return err
	}
	*rawInfo = append(*rawInfo, dByte)

	for {

		peakByte, err := d.Peek(1)
		if err != nil {
			return err
		}

		if rune(peakByte[0]) == 'e' {
			lastByte, err := d.ReadByte()
			if err != nil {
				return err
			}
			*rawInfo = append(*rawInfo, lastByte)
			break
		}

		switch rune(peakByte[0]) {
		case 'e':

		case 'i':
			if err := d.appendRawInt(rawInfo); err != nil {
				return err
			}
		case 'l':
			if err := d.appendRawList(rawInfo); err != nil {
				return err
			}
		case 'd':
			if err := d.appendRawDict(rawInfo); err != nil {
				return err
			}
		default:
			if err := d.appendRawString(rawInfo); err != nil {
				return err
			}
		}

	}
	return nil
}

func (d *decoder) getRawInfoBytes() ([]byte, error) {
	var rawInfo []byte
	d.ReadByte() // consume the 'd'
	for {

		peekBytes, err := d.Peek(1)
		if err != nil {
			panic(err)
		}

		if rune(peekBytes[0]) == 'e' {
			d.ReadByte()
			break
		}

		key, err := d.decodeString()
		if err != nil {
			panic(err)
		}

		if key == "info" {

			if err := d.appendRawDict(&rawInfo); err != nil {
				return nil, err
			}

			break
		} else {

			peekBytes, err = d.Peek(1)
			if err != nil {
				panic(err)
			}

			switch rune(peekBytes[0]) {
			case 'i':
				_, err := d.decodeInt()
				if err != nil {
					panic(err)
				}
			case 'l':
				_, err := d.decodeList()
				if err != nil {
					panic(err)
				}
			case 'd':
				_, err := d.decodeDict()
				if err != nil {
					panic(err)
				}
			default:
				_, err := d.decodeString()
				if err != nil {
					panic(err)
				}
			}

		}
	}
	return rawInfo, nil
}

func RawInfoBytes(ior io.Reader) ([]byte, error) {
	br := bufio.NewReader(ior)
	decoder := &decoder{br}
	return decoder.getRawInfoBytes()
}
