// Copyright 2021 Compass Systems
// SPDX-License-Identifier: LGPL-3.0-only

package keystore

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"os"
	"path/filepath"
	"testing"
)

func Test01(t *testing.T) {
	fpath := "./hotwallet2.json"
	keyP, err := GetWalletKey(fpath)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(hexutil.Encode(keyP))
}

func createTestFile(t *testing.T) (*os.File, string) {
	filename := "./test_key"

	fp, err := filepath.Abs(filename)
	if err != nil {
		t.Fatal(err)
	}

	file, err := os.Create(fp)
	if err != nil {
		t.Fatal(err)
	}

	return file, fp
}
