package bencode

import (
	"strconv"
)

//TODO:
// 	1.encodeString()
//	3.encodeList()
//	4.encodeDict()

func encodeInteger(num int) []byte {
	var result []byte
	result = append(result, 'i')
	result = append(result, strconv.Itoa(num)...)
	result = append(result, 'e')
	return result
}

func encodeString(s string) []byte {
	var result []byte
	result = append(result, strconv.Itoa(len(s))...)
	result = append(result, s...)
	return result
}
