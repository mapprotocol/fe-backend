package utils

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
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

func Base64ToHex(base64Str string) (string, error) {
	decodedBytes, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(decodedBytes), nil
}
