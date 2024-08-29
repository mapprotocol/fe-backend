// Copyright 2021 Compass Systems
// SPDX-License-Identifier: LGPL-3.0-only

package keystore

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
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
func Test02(t *testing.T) {
	// Read key from file.
	fpath := "./hotwallet2.json"
	keyjson, err := os.ReadFile(fpath)
	if err != nil {
		fmt.Println("Failed to read the keyfile at", fpath, err)
		return
	}
	// Decrypt key with passphrase.
	passphrase := "abc-1234"

	key, err := keystore.DecryptKey(keyjson, passphrase)
	if err != nil {
		fmt.Println("Error decrypting key", err)
		return
	}
	fmt.Println(hexutil.Encode(crypto.FromECDSA(key.PrivateKey)))
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
