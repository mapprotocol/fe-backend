package ceffu

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"sort"
)

func Sign(data string, secret string) (string, error) {
	privateKey, err := ParseRSAPrivateKey(secret)
	if err != nil {
		return "", err
	}

	hashed := sha512.Sum512([]byte(data))
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA512, hashed[:])
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(signature), nil
}

func Verify(publicKey *rsa.PublicKey, data map[string]interface{}, signBase64 string) (bool, error) {
	delete(data, "encoded")
	delete(data, "sign")

	dataBytes, err := json.Marshal(sortKeys(data))
	if err != nil {
		return false, err
	}

	decodeSign, err := base64.StdEncoding.DecodeString(signBase64)
	if err != nil {
		return false, err
	}

	hashed := sha256.Sum256(dataBytes)

	if err := rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hashed[:], decodeSign); err != nil {
		return false, err
	}
	return true, nil
}

func Decode(privateKey *rsa.PrivateKey, data string) ([]byte, error) {
	decodedData, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return []byte{}, err
	}
	decrypted, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, decodedData)
	if err != nil {
		return []byte{}, err
	}
	return decrypted, nil
}

func ParseRSAPrivateKey(privateKeyBase64 string) (*rsa.PrivateKey, error) {
	decodedPrivateKey, err := base64.StdEncoding.DecodeString(privateKeyBase64)
	if err != nil {
		return nil, err
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(decodedPrivateKey)
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}

func ParseRSAPublicKey(publicKeyBase64 string) (*rsa.PublicKey, error) {
	decodedPublicKey, err := base64.StdEncoding.DecodeString(publicKeyBase64)
	if err != nil {
		return nil, err
	}
	//publicKey, err := x509.ParsePKIXPublicKey(decodedPublicKey)
	publicKey, err := x509.ParsePKCS1PublicKey(decodedPublicKey)
	if err != nil {
		return nil, err
	}
	return publicKey, nil
}

func sortKeys(data interface{}) interface{} {
	switch v := data.(type) {
	case map[string]interface{}:
		result := make(map[string]interface{})
		keys := make([]string, 0, len(v))
		for key := range v {
			keys = append(keys, key)
		}
		sort.Strings(keys)

		for _, key := range keys {
			result[key] = sortKeys(v[key])
		}
		return result
	case []interface{}:
		sortedList := make([]interface{}, len(v))
		for i, item := range v {
			sortedList[i] = sortKeys(item)
		}
		return sortedList
	default:
		return data
	}
}
