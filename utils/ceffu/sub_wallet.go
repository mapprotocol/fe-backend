package ceffu

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	uhttp "github.com/mapprotocol/ceffu-fe-backend/utils/http"
)

type CreatSubWalletRequest struct {
	ParentWalletID int64  `json:"parentWalletId"`           // parent wallet id
	WalletName     string `json:"walletName,omitempty"`     // Sub Wallet name (Max 20 characters)
	AutoCollection int64  `json:"autoCollection,omitempty"` // Enable auto sweeping to parent wallet; ; 0: Not enable (Default Value), Suitable for API user who required Custody to maintain; asset ledger of each subaccount; ; 1: Enable, Suitable for API user who will maintain asset ledger of each subaccount at; their end.
	RequestID      int64  `json:"requestId"`                // Request identity
	Timestamp      int64  `json:"timestamp"`                // Current Timestamp
}

type DepositAddressRequest struct {
	CoinSymbol string `json:"coinSymbol"` // Coin Symbol (in capital letters); Required for Prime wallet; Not required for Qualified; wallet
	Network    string `json:"network"`    // Network symbol
	Timestamp  int64  `json:"timestamp"`  // Current Timestamp in millisecond
	WalletID   int64  `json:"walletId"`   // Sub Wallet id
}

type TransferRequest struct {
	CoinSymbol   string  `json:"coinSymbol"`   // Coin symbol
	Amount       float64 `json:"amount"`       // Transfer amount
	FromWalletID int64   `json:"fromWalletId"` // From wallet ID
	ToWalletID   int64   `json:"toWalletId"`   // To wallet ID
	RequestID    int64   `json:"requestId"`    // Client request identifier, Client provided Unique Identifier. (Max 70 characters)
	Timestamp    int64   `json:"timestamp"`    // Current timestamp in millisecond
}

type CreatSubWalletResponse struct {
	Data struct {
		WalletId          int64  `json:"walletId"`
		WalletIdStr       string `json:"walletIdStr"`
		WalletName        string `json:"walletName"`
		WalletType        uint32 `json:"walletType"`
		ParentWalletId    int64  `json:"parentWalletId"`
		ParentWalletIdStr string `json:"parentWalletIdStr"`
	} `json:"data"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

type DepositAddressResponse struct {
	Data struct {
		WalletAddress string `json:"walletAddress"`
		Memo          string `json:"memo"`
	} `json:"data"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

type TransferResponseData struct {
	OrderViewId string `json:"orderViewId"` // Transfer transaction Id
	Status      int32  `json:"status"`      // Status: 10: Pending, 20: Processing, 30: Send success, 99: Failed
	Direction   int32  `json:"direction"`   // Transfer direction: 10: prime wallet->sub wallet, 20: sub wallet->prime wallet, 30: sub wallet-> sub wallet, 40: prime wallet → prime wallet
}

type TransferResponse struct {
	Data    TransferResponseData `json:"data"`
	Code    string               `json:"code"`
	Message string               `json:"message"`
}

func getURL(path string) string {
	return fmt.Sprintf("%s%s", Domain, path)
}

func CreateSubWallet(parentWalletID int64, walletName string) (walletId int64, walletType uint32, err error) {
	headers := http.Header{
		"open-apikey":  []string{""},
		"signature":    []string{""},
		"Content-Type": []string{"application/json"},
	}
	request := CreatSubWalletRequest{
		AutoCollection: 1,
		ParentWalletID: parentWalletID,
		RequestID:      time.Now().Unix() * 1000,
		Timestamp:      time.Now().Unix() * 1000,
		WalletName:     walletName,
	}
	data, err := json.Marshal(request)
	if err != nil {
		return 0, 0, err
	}
	body := strings.NewReader(string(data))

	ret, err := uhttp.Post(getURL(PathCreateSubWallet), headers, body)
	if err != nil {
		return 0, 0, err
	}
	response := CreatSubWalletResponse{}
	if err := json.Unmarshal(ret, &response); err != nil {
		return 0, 0, err
	}
	if response.Code != SuccessCode {
		// todo encapsulated external error type
		return 0, 0, fmt.Errorf("code: %s, message: %s", response.Code, response.Message)
	}
	return response.Data.WalletId, response.Data.WalletType, nil
}

func GetDepositAddress(network, symbol string, walletID int64) (string, error) {
	headers := http.Header{
		"open-apikey": []string{""},
		"signature":   []string{""},
	}

	request := DepositAddressRequest{
		CoinSymbol: symbol,
		Network:    network,
		Timestamp:  time.Now().Unix() * 1000,
		WalletID:   walletID,
	}
	data, err := json.Marshal(request)
	if err != nil {
		return "", err
	}
	body := strings.NewReader(string(data))

	ret, err := uhttp.Get(getURL(PathGetDepositAddress), headers, body)
	if err != nil {
		return "", err
	}
	response := DepositAddressResponse{}
	if err := json.Unmarshal(ret, &response); err != nil {
		return "", err
	}
	if response.Code != SuccessCode {
		// todo encapsulated external error type
		return "", fmt.Errorf("code: %s, message: %s", response.Code, response.Message)
	}
	return response.Data.WalletAddress, nil
}

// Transfer This method allows to transfer asset between Sub Wallet and Prime Wallet Restriction:
// Only applicable to Prime wallet structure.
//
// reference: https://apidoc.ceffu.io/apidoc/shared-c9ece2c6-3ab4-4667-bb7d-c527fb3dbf78/api-3471348
func Transfer(symbol string, amount float64, fromWalletID, toWalletID int64) (*TransferResponseData, error) {
	headers := http.Header{
		"open-apikey": []string{""},
		"signature":   []string{""},
	}

	request := TransferRequest{
		CoinSymbol:   symbol,
		Amount:       amount,
		FromWalletID: fromWalletID,
		ToWalletID:   toWalletID,
		RequestID:    time.Now().Unix() * 1000,
		Timestamp:    time.Now().Unix() * 1000,
	}
	data, err := json.Marshal(request)
	if err != nil {
		return &TransferResponseData{}, err
	}
	body := strings.NewReader(string(data))

	ret, err := uhttp.Get(getURL(PathTransfer), headers, body)
	if err != nil {
		return &TransferResponseData{}, err
	}
	response := TransferResponse{}
	if err := json.Unmarshal(ret, &response); err != nil {
		return &TransferResponseData{}, err
	}
	if response.Code != SuccessCode {
		// todo encapsulated external error type
		return &TransferResponseData{}, fmt.Errorf("code: %s, message: %s", response.Code, response.Message)
	}
	return &response.Data, nil
}
