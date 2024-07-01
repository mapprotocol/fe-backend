package utils

import (
	"encoding/json"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/ethereum/go-ethereum/common"
	"strings"
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

func IsValidBitcoinAddress(address string, network *chaincfg.Params) bool {
	if _, err := btcutil.DecodeAddress(address, network); err != nil {
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
