package tonrouter

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mapprotocol/fe-backend/resource/log"
	"github.com/mapprotocol/fe-backend/utils"
	uhttp "github.com/mapprotocol/fe-backend/utils/http"
	"github.com/mapprotocol/fe-backend/utils/reqerror"
	"math/big"
	"strconv"
)

type BalanceResponse struct {
	Errno   int    `json:"errno"`
	Message string `json:"message"`
	Data    struct {
		Balance   uint   `json:"balance"`
		Decimals  uint   `json:"decimals"`
		Formatted string `json:"formatted"`
	} `json:"data"`
}

func Balance() (*big.Float, error) {
	url := fmt.Sprintf("%s%s", endpoint, PathBalance)
	log.Logger().Debugf("ton balance url: %s", url)
	ret, err := uhttp.Get(url, nil, nil)
	if err != nil {
		return nil, reqerror.NewExternalRequestError(
			url,
			reqerror.WithError(err),
		)
	}

	response := BalanceResponse{}
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
	if utils.IsEmpty(response.Data.Formatted) {
		return nil, reqerror.NewExternalRequestError(
			url,
			reqerror.WithError(errors.New("balance is empty")),
		)
	}

	balance, ok := new(big.Float).SetString(response.Data.Formatted)
	if !ok {
		return nil, reqerror.NewExternalRequestError(
			url,
			reqerror.WithError(fmt.Errorf("invalid token amount out: %s", response.Data.Formatted)),
		)
	}
	return balance, nil
}
