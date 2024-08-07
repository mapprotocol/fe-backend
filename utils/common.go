package utils

import (
	"encoding/binary"
	"encoding/json"
	"strings"
)

func IsEmpty(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

func IsDuplicateError(err string) bool {
	return strings.Contains(err, "Duplicate entry")
}

func JSON(v interface{}) string {
	bs, _ := json.Marshal(v)
	return string(bs)
}

func Uint64ToByte32(num uint64) [32]byte {
	var result [32]byte
	binary.BigEndian.PutUint64(result[:], num)
	return result
}

func BytesToUint64(b []byte) uint64 {
	return binary.BigEndian.Uint64(b)
}

func TrimHexPrefix(s string) string {
	if len(s) >= 2 && (s[0:2] == "0x" || s[0:2] == "0X") {
		return s[2:]
	}
	return s
}
