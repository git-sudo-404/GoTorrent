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

package torrent

import (
	"bufio"
	"gotorrent/internal/bencode"
	"os"
	"path/filepath"
	"testing"
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
	if _, err := metaInfoFileWriter.Write(encoder.Bytes()); err != nil {
		return err
	}
	return metaInfoFileWriter.Flush()
}

func CreateTestMetaInfo() (*MetaInfo, error) {
	metaInfoDict := NewMetaInfoDict()
	WriteMetaInfoToFile(metaInfoDict)
	defer os.Remove(GetTestMetaInfoFilePath())
	metaInfo, err := CreateMetaInfoFromFile(GetTestMetaInfoFilePath())
	if err != nil {
		return nil, err
	}
	return metaInfo, nil
}

func TestCreateMetaInfoFromFile(t *testing.T) {

	expectedDict := NewMetaInfoDict()

	WriteMetaInfoToFile(expectedDict)
	defer func() {
		os.Remove(GetTestMetaInfoFilePath())
	}()
	metaInfo, err := CreateMetaInfoFromFile(GetTestMetaInfoFilePath())
	if err != nil {
		panic(err)
	}

	infoDict := expectedDict["info"].(map[string]any)

	// Required Info fields
	if got, err := metaInfo.GetPieceLength(); err != nil || got != infoDict["piece length"].(int64) {
		t.Errorf("GetPieceLength() = %v, err = %v; want %v", got, err, infoDict["piece length"])
	}

	if got, err := metaInfo.GetPieces(); err != nil || got != infoDict["pieces"].(string) {
		t.Errorf("GetPieces() = %v, err = %v; want %v", got, err, infoDict["pieces"])
	}

	if got, err := metaInfo.GetName(); err != nil || got != infoDict["name"].(string) {
		t.Errorf("GetName() = %v, err = %v; want %v", got, err, infoDict["name"])
	}

	if got, err := metaInfo.GetLength(); err != nil || got != infoDict["length"].(int64) {
		t.Errorf("GetLength() = %v, err = %v; want %v", got, err, infoDict["length"])
	}

	// Optional Info fields
	if got, err := metaInfo.GetPrivate(); err != nil || got != infoDict["private"].(int64) {
		t.Errorf("GetPrivate() = %v, err = %v; want %v", got, err, infoDict["private"])
	}

	// Required Root fields
	if got, err := metaInfo.GetAnnounce(); err != nil || got != expectedDict["announce"].(string) {
		t.Errorf("GetAnnounce() = %v, err = %v; want %v", got, err, expectedDict["announce"])
	}

	// Optional Root fields
	if got, err := metaInfo.GetCreationDate(); err != nil || got != expectedDict["creation date"].(int64) {
		t.Errorf("GetCreationDate() = %v, err = %v; want %v", got, err, expectedDict["creation date"])
	}

	if got, err := metaInfo.GetComment(); err != nil || got != expectedDict["comment"].(string) {
		t.Errorf("GetComment() = %v, err = %v; want %v", got, err, expectedDict["comment"])
	}

	if got, err := metaInfo.GetCreatedBy(); err != nil || got != expectedDict["created by"].(string) {
		t.Errorf("GetCreatedBy() = %v, err = %v; want %v", got, err, expectedDict["created by"])
	}

	if got, err := metaInfo.GetEncoding(); err != nil || got != expectedDict["encoding"].(string) {
		t.Errorf("GetEncoding() = %v, err = %v; want %v", got, err, expectedDict["encoding"])
	}

	// Nested slices (announce-list)
	gotAnnounceList, err := metaInfo.GetAnnounceList()
	if err != nil {
		t.Errorf("GetAnnounceList() returned unexpected error: %v", err)
	} else {
		expectedAnnounceList := [][]string{
			{"http://tracker1.example.com:8080/announce"},
			{"http://tracker2.example.com:8080/announce"},
		}

		if len(gotAnnounceList) != len(expectedAnnounceList) {
			t.Errorf("GetAnnounceList() outer length = %d; want %d", len(gotAnnounceList), len(expectedAnnounceList))
		} else {
			for i := range expectedAnnounceList {
				if len(gotAnnounceList[i]) != len(expectedAnnounceList[i]) {
					t.Errorf("GetAnnounceList()[%d] inner length = %d; want %d", i, len(gotAnnounceList[i]), len(expectedAnnounceList[i]))
					continue
				}
				for j := range expectedAnnounceList[i] {
					if gotAnnounceList[i][j] != expectedAnnounceList[i][j] {
						t.Errorf("GetAnnounceList()[%d][%d] = %s; want %s", i, j, gotAnnounceList[i][j], expectedAnnounceList[i][j])
					}
				}
			}
		}
	}
}

func TestCreateMetaInfoFromFile_MissingOptionalFields(t *testing.T) {
	// Construct a minimal meta-info dictionary omitting all optional fields
	minimalDict := map[string]any{
		"info": map[string]any{
			"piece length": int64(20),
			"pieces":       "12345678901234567890",
			"name":         "minimal.txt",
			"length":       int64(1024),
		},
		"announce": "http://tracker.example.com:8080/announce",
	}

	if err := WriteMetaInfoToFile(minimalDict); err != nil {
		t.Fatalf("failed to write minimal test file: %v", err)
	}
	defer func() {
		os.Remove(GetTestMetaInfoFilePath())
	}()

	metaInfo, err := CreateMetaInfoFromFile(GetTestMetaInfoFilePath())
	if err != nil {
		t.Fatalf("CreateMetaInfoFromFile failed on minimal dict: %v", err)
	}

	// 1. Check Optional Info field (private)
	if val, err := metaInfo.GetPrivate(); err == nil || val != 0 {
		t.Errorf("GetPrivate() = %v, err = %v; want 0 and non-nil error when omitted", val, err)
	}

	// 2. Check Optional Root fields
	if val, err := metaInfo.GetAnnounceList(); err == nil || val != nil {
		t.Errorf("GetAnnounceList() = %v, err = %v; want nil slice and non-nil error when omitted", val, err)
	}

	if val, err := metaInfo.GetCreationDate(); err == nil || val != 0 {
		t.Errorf("GetCreationDate() = %v, err = %v; want 0 and non-nil error when omitted", val, err)
	}

	if val, err := metaInfo.GetComment(); err == nil || val != "" {
		t.Errorf("GetComment() = %v, err = %v; want empty string and non-nil error when omitted", val, err)
	}

	if val, err := metaInfo.GetCreatedBy(); err == nil || val != "" {
		t.Errorf("GetCreatedBy() = %v, err = %v; want empty string and non-nil error when omitted", val, err)
	}

	if val, err := metaInfo.GetEncoding(); err == nil || val != "" {
		t.Errorf("GetEncoding() = %v, err = %v; want empty string and non-nil error when omitted", val, err)
	}
}

//TODO: Test the meta-info parsing with the real .torrent file in the testdata
