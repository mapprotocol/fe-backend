package ceffu

import (
	"encoding/json"
	"fmt"
	"github.com/mapprotocol/ceffu-fe-backend/utils/reqerror"
	"net/http"
	"strings"
	"time"

	uhttp "github.com/mapprotocol/ceffu-fe-backend/utils/http"
)

type CreatSubWalletRequest struct {
	ParentWalletID string `json:"parentWalletId"`           // parent wallet id
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

type DepositHistoryRequest struct {
	WalletID   int64  `json:"walletId"`             // Prime wallet id or sub wallet id
	CoinSymbol string `json:"coinSymbol,omitempty"` // Coin symbol (in capital letters); All symbols if not specific
	Network    string `json:"network,omitempty"`    // Network symbol; All networks if not specific
	StartTime  int64  `json:"startTime"`            // Start time(timestamp in milliseconds)
	EndTime    int64  `json:"endTime"`              // End time(timestamp in milliseconds)
	PageLimit  int64  `json:"pageLimit"`            // Page limit
	PageNo     int64  `json:"pageNo"`               // Page no
	Timestamp  int64  `json:"timestamp"`            // Current Timestamp in millisecond
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

type DepositHistoryResponse struct {
	Data struct {
		Data      []*Transaction `json:"data"`
		TotalPage int            `json:"totalPage"`
		PageNo    int            `json:"pageNo"`
		PageLimit int            `json:"pageLimit"`
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

func getURLWithParams(path, params string) string {
	return fmt.Sprintf("%s%s?%s", Domain, path, params)
}

func CreateSubWallet(parentWalletID, walletName string) (walletId int64, walletType uint32, err error) {
	ts := time.Now().Unix()
	request := CreatSubWalletRequest{
		ParentWalletID: parentWalletID,
		WalletName:     walletName,
		AutoCollection: 1,
		RequestID:      ts,
		Timestamp:      ts,
	}
	data, err := json.Marshal(request)
	if err != nil {
		return 0, 0, err
	}
	dataStr := string(data)

	signature, err := Sign(dataStr, "")
	if err != nil {
		return 0, 0, err
	}
	headers := http.Header{
		"open-apikey": []string{""},
		"signature":   []string{signature},
	}

	ret, err := uhttp.Post(getURL(PathCreateSubWallet), headers, strings.NewReader(dataStr))
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
	request := DepositAddressRequest{
		CoinSymbol: symbol,
		Network:    network,
		Timestamp:  time.Now().Unix() * 1000,
		WalletID:   walletID,
	}
	params, err := uhttp.URLEncode(request)
	if err != nil {
		return "", err
	}
	url := getURLWithParams(PathGetDepositAddress, params)

	signature, err := Sign(params, "")
	if err != nil {
		return "", err
	}
	headers := http.Header{
		"open-apikey": []string{""},
		"signature":   []string{signature},
	}

	ret, err := uhttp.Get(url, headers, nil)
	if err != nil {
		return "", reqerror.NewExternalRequestError(
			url,
			reqerror.WithError(err),
		)
	}
	response := DepositAddressResponse{}
	if err := json.Unmarshal(ret, &response); err != nil {
		return "", err
	}
	if response.Code != SuccessCode {
		return "", reqerror.NewExternalRequestError(
			url,
			reqerror.WithCode(response.Code),
			reqerror.WithMessage(response.Message),
		)
	}
	return response.Data.WalletAddress, nil
}

// DepositHistory This method allows to get deposit history of the requested Wallet Id, coinSymbol and network.
// If PrimeWallet ID provided, returns sub wallet deposit history under the Prime Wallet.
// If SubWallet ID provided, returns specified sub wallet deposit history.
//
// Notes:
// walletId must be provided.
// Please notice the default startTime and endTime to make sure that time interval is within 0-30 days.
//
// reference: https://apidoc.ceffu.io/apidoc/shared-c9ece2c6-3ab4-4667-bb7d-c527fb3dbf78/api-3471585
func DepositHistory(walletID int64, symbol, network string, startTime, endTime int64, pageNo, pageLimit int64) ([]*Transaction, error) {
	request := DepositHistoryRequest{
		WalletID:   walletID,
		CoinSymbol: symbol,
		Network:    network,
		StartTime:  startTime,
		EndTime:    endTime,
		PageLimit:  pageLimit,
		PageNo:     pageNo,
		Timestamp:  time.Now().Unix() * 1000,
	}
	params, err := uhttp.URLEncode(request)
	if err != nil {
		return nil, err
	}
	url := getURLWithParams(PathDepositHistory, params)

	signature, err := Sign(params, "")
	if err != nil {
		return nil, err
	}
	headers := http.Header{
		"open-apikey": []string{""},
		"signature":   []string{signature},
	}

	ret, err := uhttp.Get(url, headers, nil)
	if err != nil {
		return nil, reqerror.NewExternalRequestError(
			url,
			reqerror.WithError(err),
		)
	}
	response := DepositHistoryResponse{}
	if err := json.Unmarshal(ret, &response); err != nil {
		return nil, err
	}
	if response.Code != SuccessCode {
		return nil, reqerror.NewExternalRequestError(
			url,
			reqerror.WithCode(response.Code),
			reqerror.WithMessage(response.Message),
		)
	}

	return response.Data.Data, nil
}

// Transfer This method allows to transfer asset between Sub Wallet and Prime Wallet Restriction:
// Only applicable to Prime wallet structure.
//
// reference: https://apidoc.ceffu.io/apidoc/shared-c9ece2c6-3ab4-4667-bb7d-c527fb3dbf78/api-3471348
func Transfer(symbol string, amount float64, fromWalletID, toWalletID int64) (*TransferResponseData, error) {
	timestamp := time.Now().Unix() * 1000
	request := TransferRequest{
		CoinSymbol:   symbol,
		Amount:       amount,
		FromWalletID: fromWalletID,
		ToWalletID:   toWalletID,
		RequestID:    timestamp,
		Timestamp:    timestamp,
	}
	data, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	dataStr := string(data)

	signature, err := Sign(dataStr, "")
	if err != nil {
		return nil, err
	}
	headers := http.Header{
		"open-apikey": []string{""},
		"signature":   []string{signature},
	}

	ret, err := uhttp.Post(getURL(PathTransfer), headers, strings.NewReader(dataStr))
	if err != nil {
		return nil, err
	}
	response := TransferResponse{}
	if err := json.Unmarshal(ret, &response); err != nil {
		return nil, err
	}
	if response.Code != SuccessCode {
		// todo encapsulated external error type
		return nil, fmt.Errorf("code: %s, message: %s", response.Code, response.Message)
	}
	return &response.Data, nil
}
