package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	uhttp "github.com/mapprotocol/ceffu-fe-backend/utils/http"
	"net/http"
	"net/url"
	"reflect"
	"sort"
	"strings"
	"time"
)

const Domain = ""

type CreatePrimeWalletRequest struct {
	RequestID  int64  `json:"requestId"` // unique Identifier
	Timestamp  int64  `json:"timestamp,omitempty"`
	WalletName string `json:"walletName"` // wallet name
	WalletType int64  `json:"walletType"` // Wallet type
}

type CreatePrimeWalletRequestResponse struct {
	Data struct {
		WalletId    int64  `json:"walletId"`
		WalletIdStr string `json:"walletIdStr"`
		WalletName  string `json:"walletName"`
		WalletType  int    `json:"walletType"`
	} `json:"data"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

type CreatSubWalletRequest struct {
	ParentWalletID int64  `json:"parentWalletId"`           // parent wallet id
	WalletName     string `json:"walletName,omitempty"`     // Sub Wallet name (Max 20 characters)
	AutoCollection int64  `json:"autoCollection,omitempty"` // Enable auto sweeping to parent wallet; ; 0: Not enable (Default Value), Suitable for API user who required Custody to maintain; asset ledger of each subaccount; ; 1: Enable, Suitable for API user who will maintain asset ledger of each subaccount at; their end.
	RequestID      int64  `json:"requestId"`                // Request identity
	Timestamp      int64  `json:"timestamp"`                // Current Timestamp
}

func CreatePrimeWallet() {
	url := Domain + "/open-api/v1/wallet/create"

	headers := http.Header{
		"open-apikey":  []string{""},
		"signature":    []string{""},
		"Content-Type": []string{"application/json"},
	}
	request := CreatePrimeWalletRequest{
		RequestID:  time.Now().Unix() * 1000,
		Timestamp:  time.Now().Unix() * 1000,
		WalletName: "neoiss-test",
		WalletType: 1,
	}
	data, err := json.Marshal(request)
	if err != nil {
		panic(err)
	}
	body := strings.NewReader(string(data))

	ret, err := uhttp.Post(url, headers, nil, body)
	if err != nil {
		panic(err)
	}
	response := CreatePrimeWalletRequestResponse{}
	if err := json.Unmarshal(ret, &response); err != nil {
		panic(err)
	}
	fmt.Println(response)
}

func CreateSubWallet() {
	url := Domain + "/open-api/v1/subwallet/create"

	headers := http.Header{
		"open-apikey":  []string{""},
		"signature":    []string{""},
		"Content-Type": []string{"application/json"},
	}
	request := CreatSubWalletRequest{
		AutoCollection: 1,
		ParentWalletID: 0,
		RequestID:      time.Now().Unix() * 1000,
		Timestamp:      time.Now().Unix() * 1000,
		WalletName:     "test",
	}
	data, err := json.Marshal(request)
	if err != nil {
		panic(err)
	}
	body := strings.NewReader(string(data))

	ret, err := uhttp.Post(url, headers, nil, body)
	if err != nil {
		panic(err)
	}
	response := CreatePrimeWalletRequestResponse{}
	if err := json.Unmarshal(ret, &response); err != nil {
		panic(err)
	}
	fmt.Println(response)
}

//func main() {
//
//	url := "/open-api/v1/wallet/create"
//	method := "POST"
//
//	payload := strings.NewReader(``)
//
//	client := &http.Client{}
//	req, err := http.NewRequest(method, url, payload)
//
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	req.Header.Add("open-apikey", "")
//	req.Header.Add("signature", "")
//	req.Header.Add("User-Agent", "Apifox/1.0.0 (https://apifox.com)")
//	req.Header.Add("Content-Type", "application/json")
//
//	res, err := client.Do(req)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	defer res.Body.Close()
//
//	body, err := ioutil.ReadAll(res.Body)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	fmt.Println(string(body))
//}

func stringify(data map[string]interface{}) (string, error) {
	sortedData := sortKeys(data)
	bytes, err := json.Marshal(sortedData)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
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

func decode(privateKey *rsa.PrivateKey, data string) (string, error) {
	decodedData, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}
	decrypted, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, decodedData) // todo
	if err != nil {
		return "", err
	}
	return string(decrypted), nil
}

func verify(publicKey *rsa.PublicKey, data map[string]interface{}) (bool, error) {
	sign, ok := data["sign"]
	if !ok {
		return false, errors.New("missing required sign")
	}

	delete(data, "encoded")
	delete(data, "sign")

	bytes, err := json.Marshal(sortKeys(data))
	if err != nil {
		return false, err
	}

	hashed := sha256.Sum256(bytes)
	decodeSign, err := base64.StdEncoding.DecodeString(sign.(string))
	if err != nil {
		return false, err
	}
	if err := rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hashed[:], decodeSign); err != nil {
		return false, err
	}
	return true, nil
}

func sign(data string, secret string) (string, error) {
	privateKeyBytes, err := base64.StdEncoding.DecodeString(secret)
	if err != nil {
		return "", err
	}

	privateKey, err := x509.ParsePKCS8PrivateKey(privateKeyBytes)
	if err != nil {
		return "", err
	}

	hashed := sha512.Sum512([]byte(data))
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey.(*rsa.PrivateKey), crypto.SHA512, hashed[:])
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(signature), nil
}

func parseRSAPrivateKey(privateKeyBase64 string) (*rsa.PrivateKey, error) {
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

func parseRSAPublicKey(publicKeyBase64 string) (*rsa.PublicKey, error) {
	decodedPublicKey, err := base64.StdEncoding.DecodeString(publicKeyBase64)
	if err != nil {
		fmt.Println("============================== 1")
		return nil, err
	}
	publicKey, err := x509.ParsePKIXPublicKey(decodedPublicKey)
	//publicKey, err := x509.ParsePKCS1PublicKey(decodedPublicKey)
	if err != nil {
		fmt.Println("============================== 2")
		return nil, err
	}
	return publicKey.(*rsa.PublicKey), nil
}

// var data = "timestamp=123456"

var secret = "MIIEvwIBADANBgkqhkiG9w0BAQEFAASCBKkwggSlAgEAAoIBAQCqGo7QFuucyIXgXMRHArVP8W+eSclsDcj+7vVJDampjuv/PdxwUzaftsliRc0wXGsyREgO83pU4/bfLmB/njsA+dJzFevlEGec4OhT279e48NUmPoGoNsX5aYq2NAuTS3oYTJh2J9SIFi8cidj1fprixYUCU5yH2FvPRdH6A0yYOFb84sUjALhwI2q8DMXV45ybgp82EO2PYzFjolkfepUHSqQ7VqsBVWsDZ08TEr72bi4mO+BulRcCKu6ywxw2KiBR1/i2ar/npvUvZNkREa68oz0KUUH1uVM76b8iMeoSIQTKAeB/EKNUS1Rup0cWuR1p1dYqRkspYPb28OhmhUpAgMBAAECggEADeWBiUp2ESbomP27Izn7af6Fad8JT4SIyRroewFcvPdqHD4Hhj2mFsIuDZM6Qhsqvr6JTH9jnQ/KmU0GoSZiF6BRKwm9bcc7T7un/0HSjoP47y5YLrZxb7BZNOLljwLLH1LhdNDnoyP1W9/Pi/5tKOAB+70O5Y/eu+G3xy4T9euGIQn2bAcbhQq48nDCUb6Ro+7gQhJfVxip7lNGWvOlGem9vaAMP7rKon4gynK+tuBTkhDQqmPACWPI++x6utsH8GiOBxRGyDgcbYqAlVeGZwRL3Q2tecxKK80xQFb+8nonZBXZUGTFXKFnC3mKRcXN+TqN3kfIcqAELUJR1QlrIQKBgQDKdNmpHJM1Fj9edJEMBY+LfkA6XPrjybH0toCvbBAA5parm2m05KVRZF8gRmrVOsV+tXkk8PH9lfcWzlLTk2803xyBPlHPnpyxFcqXoZmF1vETyJk4kqFZwS3Seq7G0NJqJpY/iBoRoZKyFjrfC3KnxCQsDUshs54M6mkqEm3PpQKBgQDXF0tA1WBkZ3bz3dU3+MeQFDCsN9Zr/6vFoHLmMz3OxLL67NJzxnbzSs0HNoxjy0L30WDhRIvI/xKFKKVscOJKD7w0oY04BSkYF0QquM/es3z8Bbciw7ttg5A+ZDHTQGjpJEMbWncZCsl+iZmri/qDslkKKq2v06nmk819wwo4NQKBgQCUKfPpGWp6HW/+1lwYajFlKt4iWE2cSs2bg7ylpPYJUrfNmw0/P13lNQmQ+zfQGRTT6EdiS5sttISCAjkHcgyeqvGXfF4vDasqxgHf+nn6QxVnHxVTG6xNnVzFftdN1SFIYjjvAdHiOVa8Uhx/g0dDk/3M52Wmomb2mM6h5Z7LqQKBgQDGqwBaKPw4oQxRIaPQaBxD6zIt0AFgja2mA5Y9JDVBp5M9i8KzJyw1efC4adzwTA1WAvH+ACcxBtCfZ7Sr3fRVvgTzhAiBJtsXIl5XK47sv1KBIfJOzQVwmOWBi2AuJL8CIPlO6Zc57SnBk+z5c3h3biMp7dOxpMq4a+qQ77afxQKBgQC0cKJLhONJTz74JNBNzHm6DRspyWHsTT92eKeptrFSTACvpijl4f4V0SR6ntT3M32a8pgXH5GMA1nfAvU5YGenzi97TX9u1z5pPuGtSDoVsXvp+VaH+HNNMrtnxAggaPlHhaA5nZhk+4fjI182ncsEF2g5LhL1bSYLSmPkrMtAAw=="

var params = map[string]interface{}{
	"encoded":   "FT6Ml+rDAvxxd0bNy/4DmqujN6y5HCYRF7+1vSKf5qNSAtd9qopNn9uetpfAPKcovZgcTpTHtiOSPY3XhTLntdbdJcbbYyIt/j95YIPriWgOi2J0XTwpvrSPtjr5L+PkMiD23bXH5qockcWdODn+Dp+WhUtO/I7SCEROzz8vO9sUxUf7dHx63ECCXDPkL/fNdlkwYirhokAqSG5HniJMzZYM53qHHZmPdj3i2DvDFSYEGMMhZOBrEEAJGCH+KFZQ11d51TGNDY9nqsgrXHtdSFKlszClNoQJj2UwaysKssRRIvjk8H1gQZ+u7w4MufBoGnCVKL3UeXTUNFjS1qjSuS63p5UlHF1J2AFvtQqGP+5ctetAtM1HLQNrmXzsLeTv83mnvUovsJsQ48+GkrWo6p0UnUgBd3MwspuUcomT03XbYeM9rNm9Cd/ht2SOGWp4In4mhTPHlLRxTnUkLgGpALJqwGZSVe7aKVlaMri7QGLJ2dBe8Rf47dLpGXilyi3b",
	"entityId":  "367070228881951744",
	"webhookId": "1783350890655723522",
	"sign":      "fE4dLqjZn0OorGuvurm0FPcpJof/rGDenej9EQGPBiH8/nH0+sPZJvds97ZrhKqfmjzba85Nf9SvfELVVfK3WQ==",
	"event":     "3",
	"timestamp": 1714129420857,
	"data": map[string]interface{}{
		"orderViewId":  "24073681492141542077760001",
		"transferType": 20,
		"direction":    20,
		"coinSymbol":   "BTC",
		"amount":       "10.11000000",
		"feeSymbol":    "BTC",
		"status":       40,
		"txTime":       1714129418000,
		"walletId":     370387791055810560,
	},
}

//varr params = map[string]interface{}{
//	"sign":      "Ih0+K/XdOw11936h9e+Lu4hiu2SnAMZC3grIr8LwKyrsXLKBkpk00Zvuu4MS2LFeJBB/opklHok1q8z5gcPaZNgOJiM3wg9Yo/hDup8I3t/YhYEGffKaEdXvi1aezCpo/Hu+YAmdlgy38yStXHNpjEQeT+QbfH1zg2dnPeD+3tOaTtpS6n4EAYWmvhCr5fy8xZry93wWgWq4zsInSEwNE2BGyhx0Xzy1moGpkPyVNpJBo/4Ngku23Ek1RL1/38zvdYzMlJuhgSrOUDf1wrFzZuwDsHohIrzLnCfq0XkFQXhR/gK+gJziG+B9iNefkQKykXfZ0DDbdzZE9Zj2v2pmHw==",
//	"timestamp": 1714129420857,
//}

var publicKeyBase64 = "MIICdwIBADANBgkqhkiG9w0BAQEFAASCAmEwggJdAgEAAoGBAL0FOAjeYpaZ1bqCsa3YD9OLnMCHgvefgYDDufZN0a+csNy+8fxlVjugfX82sYDzw4l8IacWGLRYarO74BPluPNV+Whwu/1+hRAAK0LFN2qqd0gMrA3mbe6JBccpzzOaazlG6KCd0ICPyzJGKafqZsb4SnpG3yEN1PHeaDUFvTNtAgMBAAECgYADgpJMz9xi0Y5/fSfRg56fngsWJC4RbRvZiUjtwvou2akLIFycBEG6r7tE3n4jV+N8rRpu4OEqkC3DEq0RkYWCrx1pniTe9bOlrUpVe79iHd9CE9XBGq4HpC1ABIcMNEei/PSciOCVzQxCRFpfWRa3hzwZ1y5TDlhaT5gv9bF2PQJBAN5z5C9VXMlJkKC2fE+nXqKFUgoFgFXMgOslQqmQa3W0Zj+CyD8vyj3LLS0fpetng2k81a3tzNiXqHbUkk3QP1sCQQDZhp/nsATtDwiJUG6DrzU04cY9lbzNmtb109sAweNrov6oZwETS8wQW3UWYV4IqOOl2XWhHJyjlsJebTwXVlrXAkEAyjWEhbZNurdBXaWkCG/2qTsRYQSxLMzRn25mU2ZxGDSdATxrtGxHpbYr4am0E/ErVh0zi3/vRi9Ntn7yYwNaowJBALZ50bBpH2jB4LZYC61aEDdBYqyM7SpJRyRHSYN0ItRLkncwmV1Xi2L5ZdqVaW24R+f76UpzFw/AS2MtHWiyX1cCQAy9F6S1utVJLbk3oShCZnE2Bocn/KSqIgsc9FvtjRN8VOFJ7lGwLduqgPeyr1nxMp2fjVF9KD+M+3z5bBbXh5I="

//func main() {
//	//privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
//	//if err != nil {
//	//	panic(err)
//	//}
//	////publicKey := &privateKey.PublicKey
//	//
//	//privateKeyBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
//	//if err != nil {
//	//	panic(err)
//	//}
//	//
//	//fmt.Println("============================== ", base64.StdEncoding.EncodeToString(privateKeyBytes))
//
//	//d, err := json.Marshal(data)
//	//if err != nil {
//	//	panic(err)
//	//}
//	//// 签名
//	//sig, err := sign(string(d), secret)
//	//if err != nil {
//	//	panic(err)
//	//}
//	//fmt.Println(sig)
//
//	//keys := sortKeys(params)
//	//jsonData, err := json.Marshal(keys)
//	//if err != nil {
//	//	panic(err)
//	//}
//	//fmt.Println(string(jsonData))
//
//	//privateKeyBytes, err := base64.StdEncoding.DecodeString(secret)
//	//if err != nil {
//	//	panic(err)
//	//}
//	//privateKey, err := x509.ParsePKCS8PrivateKey(privateKeyBytes)
//	//if err != nil {
//	//	panic(err)
//	//}
//	//
//	//pk := privateKey.(*rsa.PrivateKey).PublicKey
//
//	rsaPublicKey, err := parseRSAPublicKey(publicKeyBase64)
//	if err != nil {
//		panic(err)
//	}
//
//	delete(params, "encoded")
//	verified, err := verify(rsaPublicKey, params, params["sign"].(string))
//	if err != nil {
//		panic(err)
//	}
//
//	fmt.Println(verified)
//	//// 模拟数据
//	//data := make(map[string]interface{})
//	//data["message"] = "test"
//	//
//	//// 编码数据
//	//encoded := base64.StdEncoding.EncodeToString([]byte("encoded data"))
//	//data["encoded"] = encoded
//	//
//	//// 签名
//	//sign := "signature"
//	//data["sign"] = sign
//	//
//	//fmt.Println("decode:", decode(privateKey, data["encoded"].(string)))
//	//fmt.Println("verified:", verify(publicKey, data, data["sign"].(string)))
//}

func URLEncode(s interface{}) (string, error) {
	if s == nil {
		return "", errors.New("provided value is nil")
	}

	val := reflect.ValueOf(s)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return "", errors.New("provided value is not a struct")
	}

	typ := val.Type()
	urls := url.Values{}
	for i := 0; i < val.NumField(); i++ {
		if !val.Field(i).CanInterface() {
			continue
		}
		name := typ.Field(i).Name
		//urls.Add(typ.Field(i).Tag.Get("json"), fmt.Sprintf("%v", val.Field(i).Interface()))
		tag := typ.Field(i).Tag.Get("json")
		if tag != "" {
			index := strings.Index(tag, ",")
			if index == -1 {
				name = tag
			} else {
				name = tag[:index]
			}
		}
		urls.Set(name, fmt.Sprintf("%v", val.Field(i).Interface()))
	}
	return urls.Encode(), nil
}

type T struct {
	BizType    int `json:"bizType,---"`
	WalletType int `json:"walletType"`
	Timestamp  int `json:"timestamp"`
}

func main() {
	tt := T{
		BizType:    1,
		WalletType: 20,
		Timestamp:  1591999999,
	}

	_ = tt
	encode, err := URLEncode(tt)
	if err != nil {
		panic(err)
	}
	fmt.Println("============================== ", encode)
}
