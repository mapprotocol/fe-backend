package tonrouter

import (
	"encoding/json"
	"fmt"
	uhttp "github.com/mapprotocol/fe-backend/utils/http"
	"github.com/mapprotocol/fe-backend/utils/reqerror"
	"strconv"
)

type RouteAndSwapResponse struct {
	Errno   int    `json:"errno"`
	Message string `json:"message"`
	Data    struct {
		Route    *RouteData `json:"route"`
		TxParams *TxParams  `json:"txParams"`
	} `json:"data"`
}

type TxParams struct {
	To           string `json:"to"`
	Value        int    `json:"value"`
	MinAmountOut int    `json:"minAmountOut"`
	Data         string `json:"data"`
}

type RouteAndSwapRequest struct {
	FromChainID     string `json:"fromChainId"`
	ToChainID       string `json:"toChainId"`
	Amount          string `json:"amount"`
	TokenInAddress  string `json:"tokenInAddress"`
	TokenOutAddress string `json:"tokenOutAddress"`
	Receiver        string `json:"receiver"`
	Slippage        uint64 `json:"slippage"`
}

func RouteAndSwap(request *RouteAndSwapRequest) (*TxParams, error) {
	params, err := uhttp.URLEncode(request) // todo
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%s%s?%s", Domain, PathRouteAndSwap, params)
	fmt.Println("============================== url: ", url)
	ret, err := uhttp.Get(url, nil, nil)
	if err != nil {
		return nil, reqerror.NewExternalRequestError(
			url,
			reqerror.WithError(err),
		)
	}
	response := RouteAndSwapResponse{}
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
