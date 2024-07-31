package tonrouter

import (
	"encoding/json"
	"fmt"
	uhttp "github.com/mapprotocol/fe-backend/utils/http"
	"github.com/mapprotocol/fe-backend/utils/reqerror"
	"strconv"
)

type BridgeSwapRequest struct {
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	Hash     string `json:"hash"`
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
	params := fmt.Sprintf("sender=%s&receiver=%s&hash=%s", request.Sender, request.Receiver, request.Hash)
	url := fmt.Sprintf("%s%s?%s", endpoint, PathBridgeSwap, params)
	fmt.Println("============================== url: ", url)
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
	return response.Data, nil
}
