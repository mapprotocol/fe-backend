package butter

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mapprotocol/fe-backend/constants"
	"github.com/mapprotocol/fe-backend/resource/log"
	"github.com/mapprotocol/fe-backend/utils"
	uhttp "github.com/mapprotocol/fe-backend/utils/http"
	"github.com/mapprotocol/fe-backend/utils/reqerror"
	"github.com/spf13/viper"
	"math/big"
	"strconv"
)

const SuccessCode = 0

const (
	PathRoute    = "/route"
	PathGetRoute = "/getRoute"
)

var (
	endpoint string
	entrance string
)

type RouteRequest struct {
	FromChainID     string `json:"fromChainId"`
	ToChainID       string `json:"toChainId"`
	TokenInAddress  string `json:"tokenInAddress"`
	TokenOutAddress string `json:"tokenOutAddress"`
	Amount          string `json:"amount"`
	Type            string `json:"type"`
	Slippage        uint64 `json:"slippage"`
}

type RouteResponse struct {
	Errno   int                  `json:"errno"`
	Message string               `json:"message"`
	Data    []*RouteResponseData `json:"data"`
}

type RouteResponseData struct {
	GasFee struct {
		Amount string `json:"amount"`
		Symbol string `json:"symbol"`
	} `json:"gasFee"`
	Hash     string `json:"hash"`
	SrcChain struct {
		ChainId string `json:"chainId"`
		TokenIn struct {
			Address  string `json:"address"`
			Name     string `json:"name"`
			Decimals int    `json:"decimals"`
			Symbol   string `json:"symbol"`
			Icon     string `json:"icon"`
		} `json:"tokenIn"`
		TokenOut struct {
			Address  string `json:"address"`
			Name     string `json:"name"`
			Decimals int    `json:"decimals"`
			Symbol   string `json:"symbol"`
			Icon     string `json:"icon"`
		} `json:"tokenOut"`
		TotalAmountIn  string `json:"totalAmountIn"`
		TotalAmountOut string `json:"totalAmountOut"`
		Bridge         string `json:"bridge"`
	} `json:"srcChain"`
	DstChain struct {
		ChainId string `json:"chainId"`
		TokenIn struct {
			Address  string `json:"address"`
			Name     string `json:"name"`
			Decimals int    `json:"decimals"`
			Symbol   string `json:"symbol"`
			Icon     string `json:"icon"`
		} `json:"tokenIn"`
		TokenOut struct {
			Address  string `json:"address"`
			Name     string `json:"name"`
			Decimals int    `json:"decimals"`
			Symbol   string `json:"symbol"`
			Icon     string `json:"icon"`
		} `json:"tokenOut"`
		TotalAmountIn  string `json:"totalAmountIn"`
		TotalAmountOut string `json:"totalAmountOut"`
		Bridge         string `json:"bridge"`
	} `json:"dstChain"`
}

type GetRouteResponse struct {
	Errno   int          `json:"errno"`
	Message string       `json:"message"`
	Data    GetRouteData `json:"data"`
}

type GetRouteData struct {
	SrcChain struct {
		TotalAmountOut string `json:"totalAmountOut"`
	} `json:"srcChain"`
	DstChain struct {
		TotalAmountOut string `json:"totalAmountOut"`
	} `json:"dstChain"`
}

func Init() {
	cfg := viper.GetStringMapString("butter")
	entrance = cfg["entrance"]
	endpoint = cfg["endpoint"]

	if utils.IsEmpty(entrance) {
		panic("butter entrance is empty")
	}
	if utils.IsEmpty(endpoint) {
		panic("butter endpoint is empty")
	}
}
func Route(request *RouteRequest) ([]*RouteResponseData, error) {
	if request.FromChainID == request.ToChainID && request.TokenInAddress == request.TokenOutAddress {
		return getLocalRoutes(request.Amount), nil
	}
	params := fmt.Sprintf(
		"fromChainId=%s&toChainId=%s&tokenInAddress=%s&tokenOutAddress=%s&amount=%s&type=%s&slippage=%d&entrance=%s",
		request.FromChainID, request.ToChainID, request.TokenInAddress, request.TokenOutAddress, request.Amount, request.Type, request.Slippage, entrance,
	)
	url := fmt.Sprintf("%s%s?%s", endpoint, PathRoute, params)
	log.Logger().Debugf("butter route url: %s", url)
	ret, err := uhttp.Get(url, nil, nil)
	if err != nil {
		return nil, reqerror.NewExternalRequestError(
			url,
			reqerror.WithError(err),
		)
	}
	response := RouteResponse{}
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
	if response.Data == nil {
		return nil, reqerror.NewExternalRequestError(
			url,
			reqerror.WithMessage(string(ret)),
			reqerror.WithError(errors.New("data is empty")),
		)
	}
	// For the route of the same chain exchange, butter only returns the data of src chain.
	//So, the data of src chain is copied to dst chain here to be compatible with subsequent operations.
	if request.FromChainID == request.ToChainID {
		for _, data := range response.Data {
			if data != nil {
				data.DstChain = data.SrcChain
			}
		}
	}
	return response.Data, nil
}

func GetRouteAmountOut(hash string) (*big.Float, error) {
	params := fmt.Sprintf("hash=%s", hash)
	url := fmt.Sprintf("%s%s?%s", endpoint, PathGetRoute, params)
	log.Logger().Debugf("butter get route amount url: %s", url)
	ret, err := uhttp.Get(url, nil, nil)
	if err != nil {
		return nil, reqerror.NewExternalRequestError(
			url,
			reqerror.WithError(err),
		)
	}

	response := GetRouteResponse{}
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

	totalAmountOut := response.Data.DstChain.TotalAmountOut
	if utils.IsEmpty(totalAmountOut) {
		totalAmountOut = response.Data.SrcChain.TotalAmountOut
	}
	if utils.IsEmpty(totalAmountOut) {
		return nil, reqerror.NewExternalRequestError(
			url,
			reqerror.WithError(errors.New("total amount out is empty")),
		)
	}

	amountOut, ok := new(big.Float).SetString(totalAmountOut)
	if !ok {
		return nil, reqerror.NewExternalRequestError(
			url,
			reqerror.WithError(fmt.Errorf("invalid total amount out: %s", totalAmountOut)),
		)
	}
	return amountOut, nil
}

func getLocalRoutes(amount string) []*RouteResponseData {
	chin := struct {
		ChainId string `json:"chainId"`
		TokenIn struct {
			Address  string `json:"address"`
			Name     string `json:"name"`
			Decimals int    `json:"decimals"`
			Symbol   string `json:"symbol"`
			Icon     string `json:"icon"`
		} `json:"tokenIn"`
		TokenOut struct {
			Address  string `json:"address"`
			Name     string `json:"name"`
			Decimals int    `json:"decimals"`
			Symbol   string `json:"symbol"`
			Icon     string `json:"icon"`
		} `json:"tokenOut"`
		TotalAmountIn  string `json:"totalAmountIn"`
		TotalAmountOut string `json:"totalAmountOut"`
		Bridge         string `json:"bridge"`
	}{
		ChainId: constants.ChainIDOfChainPool,
		TokenIn: struct {
			Address  string `json:"address"`
			Name     string `json:"name"`
			Decimals int    `json:"decimals"`
			Symbol   string `json:"symbol"`
			Icon     string `json:"icon"`
		}{
			Address:  constants.USDTOfChainPool,
			Name:     "Tether USD",
			Decimals: constants.USDTDecimalNumberOfChainPool,
			Symbol:   "USDT",
			Icon:     "https://files.mapprotocol.io/bridge/usdt.png",
		},
		TokenOut: struct {
			Address  string `json:"address"`
			Name     string `json:"name"`
			Decimals int    `json:"decimals"`
			Symbol   string `json:"symbol"`
			Icon     string `json:"icon"`
		}{
			Address:  constants.USDTOfChainPool,
			Name:     "Tether USD",
			Decimals: constants.USDTDecimalNumberOfChainPool,
			Symbol:   "USDT",
			Icon:     "https://files.mapprotocol.io/bridge/usdt.png",
		},
		TotalAmountIn:  amount,
		TotalAmountOut: amount,
		Bridge:         constants.ExchangeNameFlushExchange,
	}

	routes := []*RouteResponseData{
		{
			GasFee: struct {
				Amount string `json:"amount"`
				Symbol string `json:"symbol"`
			}{
				Amount: constants.LocalRouteGasFee,
				Symbol: constants.NativeSymbolOfChainPool,
			},
			Hash:     constants.LocalRouteHash,
			SrcChain: chin,
			DstChain: chin,
		},
	}
	return routes
}
