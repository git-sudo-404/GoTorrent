//go:build testrun

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

package testutils

import (
	"bufio"
	"gotorrent/bencode"
	metainfo "gotorrent/meta-info"
	"os"
	"path/filepath"
)

func GetTestMetaInfoFilePath() string {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	metaInfoFilePath := filepath.Join(wd, "meta-info_test.torrent")
	return metaInfoFilePath
}

func NewMetaInfoDict() map[string]any {
	// Create a deterministic 60-byte binary pattern (3 x 20-byte SHA1 hashes)
	testPieces := string([]byte{
		0xaa, 0xf4, 0xc6, 0x1d, 0xdc, 0xa0, 0x7a, 0x7f, 0x2a, 0x08, 0x25, 0x13, 0x6c, 0xe3, 0x0c, 0x83, 0x4d, 0x0a, 0xc6, 0x3f,
		0xc1, 0xb2, 0xa1, 0x18, 0x72, 0xeb, 0x66, 0x67, 0x52, 0xf3, 0x3d, 0x87, 0x04, 0xd3, 0x8f, 0x52, 0x9a, 0x33, 0xdc, 0xc4,
		0x5f, 0xa3, 046, 0xe6, 0xd4, 0xd3, 0x89, 0x31, 0xa2, 0xa9, 0x2d, 0x1d, 0x01, 0xa8, 0x74, 0x6f, 0xd6, 0xae, 0x22, 0x32,
	})

	return map[string]any{
		"info": map[string]any{
			"piece length": int64(20),
			"pieces":       testPieces,
			"private":      int64(1),
			"name":         "TestFile.zip",
			"length":       int64(1024),
		},
		"announce": "http://tracker.example.com:8080/announce",
		"announce-list": []any{
			[]any{"http://tracker1.example.com:8080/announce"},
			[]any{"http://tracker2.example.com:8080/announce"},
		},
		"creation date": int64(12341234),
		"comment":       "This is a testing meta-info file",
		"created by":    "parser_test.go",
		"encoding":      "utf-8",
	}
}

func WriteMetaInfoToFile(metaInfoDict map[string]any) error {

	metaInfoFile, err := os.Create(GetTestMetaInfoFilePath())
	if err != nil {
		panic(err)
	}
	defer func() {
		metaInfoFile.Close()
	}()

	encoder := bencode.CreateNewEncoder()
	if err := encoder.Encode(metaInfoDict); err != nil {
		panic(err)
	}

	metaInfoFileWriter := bufio.NewWriter(metaInfoFile)
	if _, err := metaInfoFileWriter.Write([]byte(encoder.String())); err != nil {
		return err
	}
	return metaInfoFileWriter.Flush()
}

func CreateTestMetaInfo() (*metainfo.MetaInfo, error) {
	metaInfoDict := NewMetaInfoDict()
	WriteMetaInfoToFile(metaInfoDict)
	defer os.Remove(GetTestMetaInfoFilePath())
	metaInfo, err := metainfo.CreateMetaInfoFromFile(GetTestMetaInfoFilePath())
	if err != nil {
		return nil, err
	}
	return metaInfo, nil
}
