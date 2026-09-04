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

package metainfo

import (
	"bufio"
	"gotorrent/bencode"
	"os"
	"path/filepath"
	"testing"
)

func getTestMetaInfoFilePath() string {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	metaInfoFilePath := filepath.Join(wd, "meta-info_test.torrent")
	return metaInfoFilePath
}

func newMetaInfoDict() map[string]any {
	metaInfoDict := map[string]any{}
	metaInfoDict["info"] = map[string]any{
		"piece length": int64(20),
		"pieces": "\xaa\xf4\xc6\x1d\xdc\xa0\x7a\x7f\x2a\x08\x25\x13\x6c\xe3\x0c\x83\x4d\x0a\xc6\x3f" +
			"\xc1\xb2\xa1\x18\x72\xeb\x66\x67\x52\xf3\x3d\x87\x04\xd3\x8f\x52\x9a\x33\xdc\xc4" +
			"\x5f\xa3\x46\xe6\xd4\xd3\x89\x31\xa2\xa9\x2d\x1d\x01\xa8\x74\x6f\xd6\xae\x22\x32",
		"private": int64(1),
		"name":    "TestFile.zip",
		"length":  int64(1024),
	}
	metaInfoDict["announce"] = "http://tracker.example.com:8080/announce"
	metaInfoDict["announce-list"] = []any{
		[]any{"http://tracker1.example.com:8080/announce"},
		[]any{"http://tracker2.example.com:8080/announce"},
	}
	metaInfoDict["creation date"] = int64(12341234)
	metaInfoDict["comment"] = "This is a testing meta-info file"
	metaInfoDict["created by"] = "parser_test.go"
	metaInfoDict["encoding"] = "utf-8"

	return metaInfoDict
}

func writeMetaInfoToFile(metaInfoDict map[string]any) error {

	metaInfoFile, err := os.Create(getTestMetaInfoFilePath())
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
	if _, err := metaInfoFileWriter.WriteString(encoder.String()); err != nil {
		return err
	}
	return metaInfoFileWriter.Flush()
}

func TestCreateMetaInfoFromFile(t *testing.T) {
	writeMetaInfoToFile(newMetaInfoDict())
	defer func() {
		os.Remove(getTestMetaInfoFilePath())
	}()
	metaInfo, err := CreateMetaInfoFromFile(getTestMetaInfoFilePath())
	if err != nil {
		panic(err)
	}

	print(metaInfo)

}
