package tonrouter

import (
	"encoding/json"
	"fmt"
	"github.com/mapprotocol/fe-backend/resource/log"
	uhttp "github.com/mapprotocol/fe-backend/utils/http"
	"github.com/mapprotocol/fe-backend/utils/reqerror"
	"strconv"
)

type BridgeSwapRequest struct {
	Sender       string `json:"sender"`
	Receiver     string `json:"receiver"`
	FeeCollector string `json:"feeCollector"`
	FeeRatio     string `json:"feeRatio"`
	Hash         string `json:"hash"`
}

type BridgeSwapResponse struct {
	Errno   int       `json:"errno"`
	Message string    `json:"message"`
	Data    *TxParams `json:"data"`
}

type TxParams struct {
	To           string `json:"to"`
	Value        string `json:"value"`
	MinAmountOut string `json:"minAmountOut"`
	Data         string `json:"data"`
}

func BridgeSwap(request *BridgeSwapRequest) (*TxParams, error) {
	params := fmt.Sprintf(
		"sender=%s&receiver=%s&feeCollector=%s&feeRatio=%s&hash=%s",
		request.Sender, request.Receiver, request.FeeCollector, request.FeeRatio, request.Hash,
	)
	url := fmt.Sprintf("%s%s?%s", endpoint, PathBridgeSwap, params)
	log.Logger().Debugf("ton swap url: %s", url)
	ret, err := uhttp.Get(url, nil, nil)
	if err != nil {
		return nil, reqerror.NewExternalRequestError(
			url,
			reqerror.WithError(err),
		)
	}
	response := BridgeSwapResponse{}
	if err := json.Unmarshal(ret, &response); err != nil {
		return nil, reqerror.NewExternalRequestError(
			url,
			reqerror.WithMessage(string(ret)),
			reqerror.WithError(err),
		)
	}
	if response.Errno != SuccessCode {
		return nil, reqerror.NewExternalRequestError(
			url,
			reqerror.WithCode(strconv.Itoa(response.Errno)),
			reqerror.WithMessage(response.Message),
			reqerror.WithPublicError(response.Message),
		)
	}
	return response.Data, nil
}
