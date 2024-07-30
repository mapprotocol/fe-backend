package butter

import (
	"encoding/json"
	"fmt"
	uhttp "github.com/mapprotocol/fe-backend/utils/http"
	"github.com/mapprotocol/fe-backend/utils/reqerror"
	"github.com/spf13/viper"
	"strconv"
)

const SuccessCode = 0

const PathRouteAndSwap = "/routeAndSwap"

var Domain string

type RouterAndSwapRequest struct {
	FromChainID     string `json:"fromChainId"`
	ToChainID       string `json:"toChainId"`
	Amount          string `json:"amount"`
	TokenInAddress  string `json:"tokenInAddress"`
	TokenOutAddress string `json:"tokenOutAddress"`
	Type            string `json:"type"`
	Slippage        uint64 `json:"slippage"`
	Entrance        string `json:"entrance"`
	From            string `json:"from"`
	Receiver        string `json:"receiver"`
}

type RouterAndSwapResponse struct {
	Errno   int    `json:"errno"`
	Message string `json:"message"`
	Data    []struct {
		TxParam struct {
			Errno   int       `json:"errno"`
			Message string    `json:"message"`
			Data    []*TxData `json:"data"`
		} `json:"txParam"`
	} `json:"data"`
}

type TxData struct {
	To      string `json:"to"`
	Data    string `json:"data"`
	Value   string `json:"value"`
	ChainId string `json:"chainId"`
}

func init() {
	Domain = viper.GetStringMapString("endpoints")["butter"]
}

func RouterAndSwap(request *RouterAndSwapRequest) (*TxData, error) {
	params, err := uhttp.URLEncode(request) // todo
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%s%s?%s", Domain, PathRouteAndSwap, params)
	ret, err := uhttp.Get(url, nil, nil)
	if err != nil {
		return nil, reqerror.NewExternalRequestError(
			url,
			reqerror.WithError(err),
		)
	}
	response := RouterAndSwapResponse{}
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
	if response.Data[0].TxParam.Errno != SuccessCode {
		return nil, reqerror.NewExternalRequestError(
			url,
			reqerror.WithCode(strconv.Itoa(response.Errno)),
			reqerror.WithMessage(response.Message),
		)
	}
	return response.Data[0].TxParam.Data[0], nil
}
