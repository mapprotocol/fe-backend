// Copyright 2021 Compass Systems
// SPDX-License-Identifier: LGPL-3.0-only

/*
The keystore package is used to load keys from keystore files, both for live use and for testing.

# The Keystore

The keystore file is used as a file representation of a key. It contains 4 parts:
- The key type (secp256k1, sr25519)
- The PublicKey
- The Address
- The ciphertext

This keystore also requires a password to decrypt into a usable key.
The keystore library can be used to both encrypt keys into keystores, and decrypt keystore into keys.
For more information on how to encrypt and decrypt from the command line, reference the README: https://github.com/ChainSafe/ChainBridge

# The Keyring

The keyring provides predefined secp256k1 and srr25519 keys to use in testing.
These keys are automatically provided during runtime and stored in memory rather than being stored on disk.
There are 5 keys currenty supported: Alice, Bob, Charlie, Dave, and Eve.
*/
package keystore

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/term"
	"os"
	"strings"
	"syscall"
)

func getPassphrase() (string, error) {
	fmt.Print("Enter password: ")
	password, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		fmt.Printf("invalid input: %s\n", err)
		return "", err
	}
	fmt.Printf("\n")
	return strings.TrimSpace(string(password)), nil
}

func GetWalletKey(filepath string) ([]byte, error) {
	// Read key from file.
	keyjson, err := os.ReadFile(filepath)
	if err != nil {
		fmt.Println("Failed to read the keyfile at", filepath, err)
		return nil, err
	}
	// Decrypt key with passphrase.
	passphrase, err := getPassphrase()
	if err != nil {
		return nil, err
	}
	key, err := keystore.DecryptKey(keyjson, passphrase)
	if err != nil {
		fmt.Println("Error decrypting key", err)
		return nil, err
	}
	return crypto.FromECDSA(key.PrivateKey), nil
}
