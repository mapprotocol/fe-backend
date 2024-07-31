package tonrouter

import (
	"encoding/json"
	"fmt"
	uhttp "github.com/mapprotocol/fe-backend/utils/http"
	"github.com/mapprotocol/fe-backend/utils/reqerror"
	"strconv"
)

type BalanceRequest struct {
	FromChainID     string `json:"fromChainId"`
	ToChainID       string `json:"toChainId"`
	Amount          string `json:"amount"`
	TokenInAddress  string `json:"tokenInAddress"`
	TokenOutAddress string `json:"tokenOutAddress"`
	Receiver        string `json:"receiver"`
	Slippage        uint64 `json:"slippage"`
}

type BalanceResponse struct {
	Errno   int    `json:"errno"`
	Message string `json:"message"`
	Data    struct {
		Balance   uint   `json:"balance"`
		Decimals  uint   `json:"decimals"`
		Formatted string `json:"formatted"`
	} `json:"data"`
}

func Balance() (string, error) {
	url := fmt.Sprintf("%s%s", endpoint, PathBalance)
	ret, err := uhttp.Get(url, nil, nil)
	if err != nil {
		return "0", reqerror.NewExternalRequestError(
			url,
			reqerror.WithError(err),
		)
	}
	response := BalanceResponse{}
	if err := json.Unmarshal(ret, &response); err != nil {
		return "0", err
	}
	if response.Errno != SuccessCode {
		return "0", reqerror.NewExternalRequestError(
			url,
			reqerror.WithCode(strconv.Itoa(response.Errno)),
			reqerror.WithMessage(response.Message),
		)
	}
	return response.Data.Formatted, nil
}
