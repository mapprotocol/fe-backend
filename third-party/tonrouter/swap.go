package tonrouter

import (
	"encoding/json"
	"fmt"
	uhttp "github.com/mapprotocol/fe-backend/utils/http"
	"github.com/mapprotocol/fe-backend/utils/reqerror"
	"github.com/spf13/viper"
	"strconv"
)

const SuccessCode = 0

const (
	PathBridgeSwap = "/bridger/swap"
)

var Domain string

type BridgeSwapResponse struct {
	Errno   int    `json:"errno"`
	Message string `json:"message"`
	Data    struct {
		TxParams *TxParams `json:"txParams"`
	} `json:"data"`
}

type TxParams struct {
	To    string `json:"to"`
	Value string `json:"value"`
	Data  string `json:"data"`
}

type BridgeSwapRequest struct {
	Amount          string `json:"amount"`
	Slippage        uint64 `json:"slippage"`
	TokenOutAddress string `json:"tokenOutAddress"`
	Receiver        string `json:"receiver"`
	OrderID         uint64 `json:"orderId"`
}

func init() {
	Domain = viper.GetStringMapString("endpoints")["filter"]
}

func BridgeSwap(request *BridgeSwapRequest) (*TxParams, error) {
	params, err := uhttp.URLEncode(request) // todo
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%s%s?%s", Domain, PathBridgeSwap, params)
	ret, err := uhttp.Get(url, nil, nil)
	if err != nil {
		return nil, reqerror.NewExternalRequestError(
			url,
			reqerror.WithError(err),
		)
	}
	response := BridgeSwapResponse{}
	if err := json.Unmarshal(ret, &response); err != nil {
		return nil, err
	}
	if response.Errno != SuccessCode {
		return nil, reqerror.NewExternalRequestError(
			url,
			reqerror.WithCode(strconv.Itoa(response.Errno)),
			reqerror.WithMessage(response.Message),
		)
	}
	return response.Data.TxParams, nil
}
