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
	"fmt"
	"gotorrent/bencode"
	"os"
)

func checkFieldPresenceAndType[T any](metaInfoDict map[string]any, field string) (T, error) {
	var zero T
	value, ok := metaInfoDict[field]
	if !ok {
		return zero, fmt.Errorf("Meta-Info missing needed field")
	}
	if v, ok := value.(T); !ok {
		return zero, fmt.Errorf("Undesired field value type")
	} else {
		return v, nil
	}
}

func checkOptionalFieldPresenceAndType[T any](metaInfoDict map[string]any, field string) (T, bool, error) {
	var zero T
	value, ok := metaInfoDict[field]
	if !ok {
		return zero, false, nil
	}
	if v, ok := value.(T); !ok {
		return zero, false, fmt.Errorf("Undesired field value type")
	} else {
		return v, true, nil
	}
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
	pieceLength, err := checkFieldPresenceAndType[int64](metaInfoDict, "piece length")
	if err != nil {
		return nil, err
	}
	metaInfo.setPieceLength(pieceLength)

	pieces, err := checkFieldPresenceAndType[string](metaInfoDict, "pieces")
	if err != nil {
		return nil, err
	}
	metaInfo.setPieces(pieces)

	private, present, err := checkOptionalFieldPresenceAndType[int64](metaInfoDict, "private")
	if err != nil {
		return nil, err
	}
	if present {
		metaInfo.setPrivate(private)
	}

	name, err := checkFieldPresenceAndType[string](metaInfoDict, "name")
	if err != nil {
		return nil, err
	}
	metaInfo.setName(name)

	length, err := checkFieldPresenceAndType[int64](metaInfoDict, "length")
	if err != nil {
		return nil, err
	}
	metaInfo.setLength(length)

	announce, err := checkFieldPresenceAndType[string](metaInfoDict, "announce")
	if err != nil {
		return nil, err
	}
	metaInfo.setAnnounce(announce)

	announceList, present, err := checkOptionalFieldPresenceAndType[[][]string](metaInfoDict, "announce-list")
	if err != nil {
		return nil, err
	}
	if present {
		metaInfo.setAnnounceList(announceList)
	}

	creationDate, present, err := checkOptionalFieldPresenceAndType[int64](metaInfoDict, "creation date")
	if err != nil {
		return nil, err
	}
	if present {
		metaInfo.setCreationDate(creationDate)
	}

	comment, present, err := checkOptionalFieldPresenceAndType[string](metaInfoDict, "comment")
	if err != nil {
		return nil, err
	}
	if present {
		metaInfo.setComment(comment)
	}

	createdBy, present, err := checkOptionalFieldPresenceAndType[string](metaInfoDict, "created by")
	if err != nil {
		return nil, err
	}
	if present {
		metaInfo.setCreatedBy(createdBy)
	}

	encoding, present, err := checkOptionalFieldPresenceAndType[string](metaInfoDict, "encoding")
	if err != nil {
		return nil, err
	}
	if present {
		metaInfo.setEncoding(encoding)
	}

	return metaInfo, nil
}
