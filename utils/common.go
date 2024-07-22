package utils

import (
	"encoding/binary"
	"encoding/json"
	"strings"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/ethereum/go-ethereum/common"
)

func IsEmpty(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

func IsDuplicateError(err string) bool {
	return strings.Contains(err, "Duplicate entry")
}

func IsValidEvmAddress(address string) bool {
	return common.IsHexAddress(address)
}

func IsValidEvmHash(address string) bool {
	return IsHexHash(address)
}

func IsValidBitcoinAddress(address string, network *chaincfg.Params) bool {
	if _, err := btcutil.DecodeAddress(address, network); err != nil {
		return false
	}
	return true
}

func IsValidBitcoinHash(hash string) bool {
	if _, err := chainhash.NewHashFromStr(hash); err != nil {
		return false
	}
	return true
}

func JSON(v interface{}) string {
	bs, _ := json.Marshal(v)
	return string(bs)
}

func ValidatePage(page, size int) (int, int) {
	if page <= 0 {
		page = 1
	}

	switch {
	case size > 100:
		size = 100
	case size <= 0:
		size = 20
	}
	return page, size
}

func IsHexHash(s string) bool {
	if has0xPrefix(s) {
		s = s[2:]
	}
	return len(s) == 2*common.HashLength && isHex(s)
}

func Uint64ToByte32(num uint64) [32]byte {
	var result [32]byte
	binary.BigEndian.PutUint64(result[:], num)
	return result
}

// has0xPrefix validates str begins with '0x' or '0X'.
func has0xPrefix(str string) bool {
	return len(str) >= 2 && str[0] == '0' && (str[1] == 'x' || str[1] == 'X')
}

// isHexCharacter returns bool of c being a valid hexadecimal.
func isHexCharacter(c byte) bool {
	return ('0' <= c && c <= '9') || ('a' <= c && c <= 'f') || ('A' <= c && c <= 'F')
}

// isHex validates whether each byte is valid hexadecimal string.
func isHex(str string) bool {
	if len(str)%2 != 0 {
		return false
	}
	for _, c := range []byte(str) {
		if !isHexCharacter(c) {
			return false
		}
	}
	return true
}
