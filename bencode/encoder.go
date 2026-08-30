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
	"sort"
	"strconv"
)

func encodeInteger(num int) []byte { // <number> -> i<number>e
	var result []byte
	result = append(result, 'i')
	result = append(result, strconv.Itoa(num)...)
	result = append(result, 'e')
	return result
}

func encodeString(s string) []byte { // <string> -> len(<string>):<string>
	var result []byte
	result = append(result, strconv.Itoa(len(s))...)
	result = append(result, ':')
	result = append(result, s...)
	return result
}

func encodeList(list []any) []byte { // <List> -> I<List>e
	var result []byte
	result = append(result, 'l')
	for _, item := range list {
		switch v := item.(type) {
		case int:
			result = append(result, encodeInteger(v)...)
		case string:
			result = append(result, encodeString(v)...)
		}
	}
	result = append(result, 'e')
	return result
}

func encodeDict(dict map[string]any) []byte { // <Dict> -> d<Key><Val>e
	// Key -> only strings
	// Value -> can be any of these [int, string, list, dict]
	var result []byte
	result = append(result, 'd')

	keys := make([]string, 0, len(dict))
	for key, _ := range dict {
		keys = append(keys, key)
	}

	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	for _, key := range keys {
		result = append(result, encodeString(key)...)
		value := dict[key]
		switch v := value.(type) {
		case int:
			result = append(result, encodeInteger(v)...)
		case string:
			result = append(result, encodeString(v)...)
		case []any:
			result = append(result, encodeList(v)...)
		case map[string]any:
			result = append(result, encodeDict(v)...)
		}
	}
	result = append(result, 'e')
	return result
}
