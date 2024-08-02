package butter

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mapprotocol/fe-backend/resource/log"
	uhttp "github.com/mapprotocol/fe-backend/utils/http"
	"github.com/mapprotocol/fe-backend/utils/reqerror"
	"strconv"
)

const PathSwap = "/swap"

var (
	ErrNotFoundTxData = errors.New("not found tx data")
)

type SwapRequest struct {
	From     string `json:"from"`
	Receiver string `json:"receiver"`
	Hash     string `json:"hash"`
	Slippage uint64 `json:"slippage"`
	CallData string `json:"callData"`
}

type SwapResponse struct {
	Errno   int       `json:"errno"`
	Message string    `json:"message"`
	Data    []*TxData `json:"data"`
}

type TxData struct {
	To      string `json:"to"`
	Data    string `json:"data"`
	Value   string `json:"value"`
	ChainId string `json:"chainId"`
}

func Swap(request *SwapRequest) (*TxData, error) { // todo checkout code
	params := fmt.Sprintf(
		"from=%s&receiver=%s&hash=%s&slippage=%d&callData=%s",
		request.From, request.Receiver, request.Hash, request.Slippage, request.CallData,
	)
	url := fmt.Sprintf("%s%s?%s", endpoint, PathSwap, params)
	log.Logger().Debugf("butter swap url: %s", url)
	ret, err := uhttp.Get(url, nil, nil)
	if err != nil {
		return nil, reqerror.NewExternalRequestError(
			url,
			reqerror.WithError(err),
		)
	}
	response := SwapResponse{}
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
	if len(response.Data) == 0 {
		return nil, reqerror.NewExternalRequestError(
			url,
			reqerror.WithCode(strconv.Itoa(response.Errno)),
			reqerror.WithMessage(response.Message),
			reqerror.WithError(ErrNotFoundTxData),
			reqerror.WithPublicError(ErrNotFoundTxData.Error()),
		)
	}
	return response.Data[0], nil
}
