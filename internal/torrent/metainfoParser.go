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
	"fmt"
	bencode "gotorrent/internal/bencode"
	"os"
)

// since announce-list is optional this function returns an err if the announce list key is not found in the metaInfoDict
func checkAnnounceListPresenceAndType(metaInfoDict map[string]any) ([][]string, bool, error) {

	outerVal, ok := metaInfoDict["announce-list"]
	if !ok {
		return nil, false, nil
	}

	outerList, ok := outerVal.([]any)
	if !ok {
		return nil, false, fmt.Errorf("announce-list type does not match the required type")
	}

	var result [][]string

	for _, innerVal := range outerList {
		innerList, ok := innerVal.([]any)
		if !ok {
			return nil, false, fmt.Errorf("announce-list type does not meet the required type")
		}
		var innerResult []string
		for _, items := range innerList {
			item, ok := items.(string)
			if !ok {
				return nil, false, fmt.Errorf("announce-list type does not meet the required type")
			}
			innerResult = append(innerResult, item)
		}
		result = append(result, innerResult)
	}

	return result, true, nil
}

func CreateMetaInfoFromFile(filePath string) (*MetaInfo, error) {
	metaInfo := NewMetaInfo()

	metaInfoFile, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("Error opening the meta-info file")
	}
	defer metaInfoFile.Close()

	metaInfoDict, err := bencode.Decode(metaInfoFile)
	if err != nil {
		return nil, fmt.Errorf("Error Decoding the meta-info file")
	}

	//NOTE:Standard .torrent files use space-separated keys
	infoRaw, ok := metaInfoDict["info"]
	if !ok {
		panic(fmt.Errorf("info not in decoded meta-info file"))
	}
	info, ok := infoRaw.(map[string]any)
	if !ok {
		panic(fmt.Errorf("info not of desired type : map[string]any"))
	}
	pieceLength, err := checkFieldPresenceAndType[int64](info, "piece length")
	if err != nil {
		return nil, err
	}
	metaInfo.SetPieceLength(pieceLength)

	pieces, err := checkFieldPresenceAndType[string](info, "pieces")
	if err != nil {
		return nil, err
	}
	metaInfo.SetPieces(pieces)

	private, present, err := checkOptionalFieldPresenceAndType[int64](info, "private")
	if err != nil {
		return nil, err
	}
	if present {
		metaInfo.SetPrivate(private)
	}

	name, err := checkFieldPresenceAndType[string](info, "name")
	if err != nil {
		return nil, err
	}
	metaInfo.SetName(name)

	length, err := checkFieldPresenceAndType[int64](info, "length")
	if err != nil {
		return nil, err
	}
	metaInfo.SetLength(length)

	announce, err := checkFieldPresenceAndType[string](metaInfoDict, "announce")
	if err != nil {
		return nil, err
	}
	metaInfo.SetAnnounce(announce)

	announceList, present, err := checkAnnounceListPresenceAndType(metaInfoDict)
	if err != nil {
		return nil, err
	}
	if present {
		metaInfo.SetAnnounceList(announceList)
	}

	creationDate, present, err := checkOptionalFieldPresenceAndType[int64](metaInfoDict, "creation date")
	if err != nil {
		return nil, err
	}
	if present {
		metaInfo.SetCreationDate(creationDate)
	}

	comment, present, err := checkOptionalFieldPresenceAndType[string](metaInfoDict, "comment")
	if err != nil {
		return nil, err
	}
	if present {
		metaInfo.SetComment(comment)
	}

	createdBy, present, err := checkOptionalFieldPresenceAndType[string](metaInfoDict, "created by")
	if err != nil {
		return nil, err
	}
	if present {
		metaInfo.SetCreatedBy(createdBy)
	}

	encoding, present, err := checkOptionalFieldPresenceAndType[string](metaInfoDict, "encoding")
	if err != nil {
		return nil, err
	}
	if present {
		metaInfo.SetEncoding(encoding)
	}

	return metaInfo, nil
}
