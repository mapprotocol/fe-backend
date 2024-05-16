package ceffu

import (
	"encoding/json"
	"fmt"
	"github.com/mapprotocol/ceffu-fe-backend/utils/reqerror"
	"net/http"
	"strconv"
	"strings"
	"time"

	uhttp "github.com/mapprotocol/ceffu-fe-backend/utils/http"
)

type CreatePrimeWalletRequest struct {
	WalletName string `json:"walletName"` // wallet name
	WalletType int64  `json:"walletType"` // Wallet type
	RequestID  int64  `json:"requestId"`  // unique Identifier
	Timestamp  int64  `json:"timestamp,omitempty"`
}

type WithdrawalRequest struct {
	Amount             string `json:"amount"`                       // withdrawal amount
	CoinSymbol         string `json:"coinSymbol"`                   // coin symbol
	Memo               string `json:"memo,omitempty"`               // memo/address tag
	Network            string `json:"network"`                      // network symbol
	WalletID           int64  `json:"walletId"`                     // wallet id
	WithdrawalAddress  string `json:"withdrawalAddress"`            // withdrawal address or to wallet id str  must have one
	ToWalletIDStr      string `json:"toWalletIdStr"`                // to wallet id str  or withdrawal address must have one
	CustomizeFeeAmount string `json:"customizeFeeAmount,omitempty"` // User-specified fee  , now support eth
	RequestID          int64  `json:"requestId"`                    // Unique Identifier
	Timestamp          int64  `json:"timestamp"`                    // Current Timestamp in millisecond
}

type WithdrawalDetailRequest struct {
	OrderViewID string `json:"orderViewId,omitempty"` // Withdrawal Transaction Id
	RequestID   string `json:"requestId,omitempty"`   // Client request identifier: Universal Unique identifier provided by the client side.
	Timestamp   int64  `json:"timestamp"`             // Current Timestamp in millisecond
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

type WithdrawalResponseData struct {
	OrderViewId  string `json:"orderViewId"`
	Status       int    `json:"status"`
	TransferType int    `json:"transferType"`
}

type WithdrawalResponse struct {
	Data    WithdrawalResponseData `json:"data"`
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
}

type WithdrawalDetailResponse struct {
	Code    string                        `json:"code"`    // response code, '000000' when successed, others represent there some error occured
	Data    *WithdrawalDetailResponseData `json:"data"`    // response data, maybe null
	Message string                        `json:"message"` // detail of response, when code != '000000', it's detail of error
}

type WithdrawalDetailResponseData struct {
	OrderViewID  string  `json:"orderViewId"`
	TxID         string  `json:"txId"` // transaction id (Only Applicable to on-chain transfer)
	TransferType int64   `json:"transferType"`
	Direction    int64   `json:"direction"`
	FromAddress  string  `json:"fromAddress"`
	ToAddress    string  `json:"toAddress"`
	Network      string  `json:"network"`
	CoinSymbol   string  `json:"coinSymbol"`
	Amount       string  `json:"amount"`
	FeeSymbol    string  `json:"feeSymbol"`
	FeeAmount    string  `json:"feeAmount"`
	Status       int64   `json:"status"`
	Memo         *string `json:"memo"`
	TxTime       string  `json:"txTime"`
	WalletStr    string  `json:"walletStr"`
	RequestID    *string `json:"requestId"` // universal unique identifier provided by the client side.
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

	ret, err := uhttp.Post(url, headers, body)
	if err != nil {
		panic(err)
	}
	response := CreatePrimeWalletRequestResponse{}
	if err := json.Unmarshal(ret, &response); err != nil {
		panic(err)
	}
	fmt.Println(response)
}

// Withdrawal This method enables the withdrawal of funds from the specified wallet to an external address or a Ceffu wallet.
// The withdrawal endpoint is applicable only to parent Qualified wallet ID or Cosign wallet or parent Prime wallet ID.
// To indicate the destination address, either 'withdrawalAddress' or 'ToWalletIdStr' must be provided.
// If the destination address is a Ceffu wallet address, the whitelisted address verification will be bypassed.
//
// IMPORTANT NOTES: The amount field in Withdrawal (v2) endpoint means withdrawal amount excluded network fee in v2,
// that is exact amount receiver will receive.
// Please use Get Withdrawal History v2 and Get Withdrawal Detail (v2) together with Withdrawal (v2).
//
// reference: https://apidoc.ceffu.io/apidoc/shared-c9ece2c6-3ab4-4667-bb7d-c527fb3dbf78/api-3471332
func Withdrawal(request *WithdrawalRequest) (*WithdrawalResponseData, error) {
	// todo padding
	headers := http.Header{
		"open-apikey": []string{""},
		"signature":   []string{""},
	}

	request.RequestID = time.Now().Unix() * 1000
	request.Timestamp = time.Now().Unix() * 1000

	data, err := json.Marshal(request)
	if err != nil {
		return &WithdrawalResponseData{}, err
	}
	body := strings.NewReader(string(data))

	ret, err := uhttp.Post(getURL(PathWithdrawal), headers, body)
	if err != nil {
		return &WithdrawalResponseData{}, err
	}
	response := WithdrawalResponse{}
	if err := json.Unmarshal(ret, &response); err != nil {
		return &WithdrawalResponseData{}, err
	}
	if response.Code != SuccessCode {
		// todo encapsulated external error type
		return &WithdrawalResponseData{}, fmt.Errorf("code: %s, message: %s", response.Code, response.Message)
	}
	return &response.Data, nil
}

func WithdrawalDetail(orderViewId string) (*WithdrawalDetailResponseData, error) {
	// todo padding
	headers := http.Header{
		"open-apikey": []string{""},
		"signature":   []string{""},
	}

	timestamp := time.Now().Unix() * 1000
	requestID := strconv.FormatInt(timestamp, 10)
	params := fmt.Sprintf("?orderViewId=%s&requestId=%s&timestamp=%d", orderViewId, requestID, timestamp)
	url := getURL(PathWithdrawalDetail) + params

	ret, err := uhttp.Get(url, headers, nil)
	if err != nil {
		return &WithdrawalDetailResponseData{},
			reqerror.NewExternalRequestError(
				url,
				reqerror.WithMethod(http.MethodGet),
				reqerror.WithError(err),
			)
	}
	response := WithdrawalDetailResponse{}
	if err := json.Unmarshal(ret, &response); err != nil {
		return &WithdrawalDetailResponseData{}, err
	}
	if response.Code != SuccessCode {
		return &WithdrawalDetailResponseData{},
			reqerror.NewExternalRequestError(
				getURL(PathWithdrawal),
				reqerror.WithCode(response.Code),
				reqerror.WithMessage(response.Message),
			)
	}
	return response.Data, nil
}
