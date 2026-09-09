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

//NOTE:
// This struct and everything is just used for the internal purposes , the gotorrent still uses the .torrent bencoded dictionary
//TODO:
// Current Client implementation need not cover the creation of the .torrent file , but then if this would get wrapped in to as a
// distributed p2p file sharing app then the implementation of the .torrent file creation is needed

package torrent

import (
	"fmt"
	"gotorrent/bencode"
	"os"
)

//TODO:
//   1. Add multi-file torrent

type info struct {
	// fields common to both single file & mulit file torrent
	pieceLength int64  //  number of bytes in each piece (integer)
	pieces      string // consisting of the concatenation of all 20-byte SHA1 hash values, one per piece (byte string, i.e. not urlencoded)
	private     *int64 // (optional) this field is an integer. If it is set to "1", the client MUST publish its presence
	//to get other peers ONLY via the trackers explicitly described in the metainfo file. If this field is set to "0" or is not present,
	//the client may obtain peer from other means, e.g. PEX peer exchange, dht. Here, "private" may be read as "no external peer source".

	// single-file mode torrent info
	name   string // the filename. This is purely advisory.
	length int64  // length of the file in bytes (integer)
}

type MetaInfo struct {
	//info dictionary (here, single file torrent)
	info
	announce     string      // The announce URL of the tracker (string)
	announceList *[][]string //(optional) this is an extention to the official specification, offering backwards-compatibility. (list of lists of strings).
	creationDate *int64      //  (optional) the creation time of the torrent, in standard UNIX epoch format (integer, seconds since 1-Jan-1970 00:00:00 UTC)
	comment      *string     //   (optional) free-form textual comments of the author (string)
	createdBy    *string     //  (optional) name and version of the program used to create the .torrent (string)
	encoding     *string     // (optional) the string encoding format used to generate the pieces part of the info dictionary in the .torrent metafile (string)
}

func NewMetaInfo() *MetaInfo {
	return &MetaInfo{
		info: info{
			pieceLength: 0,
			pieces:      "",
			name:        "",
			length:      0,
		},
		announce: "",
	}
}

func (mi *MetaInfo) SetPieceLength(pieceLength int64) *MetaInfo {
	mi.pieceLength = pieceLength
	return mi
}

func (mi *MetaInfo) SetPieces(pieces string) *MetaInfo {
	mi.pieces = pieces
	return mi
}

func (mi *MetaInfo) SetPrivate(private int64) *MetaInfo {
	mi.private = &private
	return mi
}

func (mi *MetaInfo) SetName(name string) *MetaInfo {
	mi.name = name
	return mi
}

func (mi *MetaInfo) SetLength(length int64) *MetaInfo {
	mi.length = length
	return mi
}

func (mi *MetaInfo) SetAnnounce(announce string) *MetaInfo {
	mi.announce = announce
	return mi
}

func (mi *MetaInfo) SetAnnounceList(announceList [][]string) *MetaInfo {
	mi.announceList = &announceList
	return mi
}

func (mi *MetaInfo) SetCreationDate(creationDate int64) *MetaInfo {
	mi.creationDate = &creationDate
	return mi
}

func (mi *MetaInfo) SetComment(comment string) *MetaInfo {
	mi.comment = &comment
	return mi
}

func (mi *MetaInfo) SetCreatedBy(createdBy string) *MetaInfo {
	mi.createdBy = &createdBy
	return mi
}

func (mi *MetaInfo) SetEncoding(encoding string) *MetaInfo {
	mi.encoding = &encoding
	return mi
}

func (mi *MetaInfo) GetPieceLength() (int64, error) {
	return mi.pieceLength, nil
}

func (mi *MetaInfo) GetPieces() (string, error) {
	return mi.pieces, nil
}

func (mi *MetaInfo) GetPrivate() (int64, error) {
	if mi.private == nil {
		return 0, fmt.Errorf("Optional Field Not Set , returning 0 instead")
	}
	return *mi.private, nil
}

func (mi *MetaInfo) GetName() (string, error) {
	return mi.name, nil
}

func (mi *MetaInfo) GetLength() (int64, error) {
	return mi.length, nil
}

func (mi *MetaInfo) GetAnnounce() (string, error) {
	return mi.announce, nil
}

func (mi *MetaInfo) GetAnnounceList() ([][]string, error) {
	if mi.announceList == nil {
		return nil, fmt.Errorf("Optional Field Not Set , returning nil instead")
	}
	return *mi.announceList, nil
}

func (mi *MetaInfo) GetCreationDate() (int64, error) {
	if mi.creationDate == nil {
		return 0, fmt.Errorf("Optional Field Not Set , returning 0 instead")
	}
	return *mi.creationDate, nil
}

func (mi *MetaInfo) GetComment() (string, error) {
	if mi.comment == nil {
		return "", fmt.Errorf("Optional Field Not Set , returning empty string instead")
	}
	return *mi.comment, nil
}

func (mi *MetaInfo) GetCreatedBy() (string, error) {
	if mi.createdBy == nil {
		return "", fmt.Errorf("Optional Field Not Set , returning empty string instead")
	}
	return *mi.createdBy, nil
}

func (mi *MetaInfo) GetEncoding() (string, error) {
	if mi.encoding == nil {
		return "", fmt.Errorf("Optional Field Not Set , returning empty string instead")
	}
	return *mi.encoding, nil
}

func (mi *MetaInfo) GetInfoDict() (map[string]any, error) {
	var info map[string]any
	info = map[string]any{}

	pieceLength, err := mi.GetPieceLength()
	if err != nil {
		return nil, err
	}
	info["piece length"] = pieceLength

	pieces, err := mi.GetPieces()
	if err != nil {
		return nil, err
	}
	info["pieces"] = pieces

	private, err := mi.GetPrivate()
	if err == nil {
		info["private"] = private
	}

	name, err := mi.GetName()
	if err != nil {
		return nil, err
	}
	info["name"] = name

	length, err := mi.GetLength()
	if err != nil {
		return nil, err
	}
	info["length"] = length

	return info, nil
}

func checkFieldPresenceAndType[T any](metaInfoDict map[string]any, field string) (T, error) {
	var zero T
	value, ok := metaInfoDict[field]
	if !ok {
		return zero, fmt.Errorf("Meta-Info missing needed field : %s", field)
	}
	if v, ok := value.(T); !ok {
		return zero, fmt.Errorf("Undesired field value type  , Field : %s ", field)
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
		return zero, false, fmt.Errorf("Undesired field value type , Field : %s ", field)
	} else {
		return v, true, nil
	}
}

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
