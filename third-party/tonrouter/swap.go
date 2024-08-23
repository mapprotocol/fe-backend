package tonrouter

import (
	"encoding/json"
	"fmt"
	"github.com/mapprotocol/fe-backend/resource/log"
	uhttp "github.com/mapprotocol/fe-backend/utils/http"
	"github.com/mapprotocol/fe-backend/utils/reqerror"
	"github.com/spf13/viper"
	"strconv"
)

const SuccessCode = 0

const (
	PathBridgeSwap = "/v2/bridger/swap"
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

func Init() {
	Domain = viper.GetStringMapString("endpoints")["tonrouter"]
}

func BridgeSwap(req *BridgeSwapRequest) (*TxParams, error) {
	params := fmt.Sprintf(
		"amount=%s&slippage=%d&tokenOutAddress=%s&receiver=%s&orderId=%d",
		req.Amount, req.Slippage, req.TokenOutAddress, req.Receiver, req.OrderID,
	)
	url := fmt.Sprintf("%s%s?%s", Domain, PathBridgeSwap, params)
	log.Logger().Debug(fmt.Sprintf("bridge swap url: %s", url))
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
