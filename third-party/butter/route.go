package butter

import (
	"encoding/json"
	"fmt"
	uhttp "github.com/mapprotocol/fe-backend/utils/http"
	"github.com/mapprotocol/fe-backend/utils/reqerror"
	"strconv"
)

const Domain = "https://bs-router-test.chainservice.io"

const SuccessCode = 0

const PathRoute = "/route"

type RouteRequest struct {
	FromChainID     string `json:"fromChainId"`
	ToChainID       string `json:"toChainId"`
	Amount          string `json:"amount"`
	TokenInAddress  string `json:"tokenInAddress"`
	TokenOutAddress string `json:"tokenOutAddress"`
	Kind            string `json:"type"`
	Slippage        string `json:"slippage"`
	Entrance        string `json:"entrance"`
}

type RouteResponse struct {
	Errno   int                  `json:"errno"`
	Message string               `json:"message"`
	Data    []*RouteResponseData `json:"data"`
}

type RouteResponseData struct {
	Diff      string `json:"diff"`
	BridgeFee struct {
		Amount string `json:"amount"`
		Symbol string `json:"symbol"`
	} `json:"bridgeFee"`
	TradeType int `json:"tradeType"`
	GasFee    struct {
		Amount string `json:"amount"`
		Symbol string `json:"symbol"`
	} `json:"gasFee"`
	GasEstimated  string `json:"gasEstimated"`
	TimeEstimated int    `json:"timeEstimated"`
	Hash          string `json:"hash"`
	Timestamp     int64  `json:"timestamp"`
	SrcChain      struct {
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
		Route          []struct {
			AmountIn  string `json:"amountIn"`
			AmountOut string `json:"amountOut"`
			DexName   string `json:"dexName"`
			Path      []struct {
				Id      string `json:"id"`
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
				Fee string `json:"fee"`
			} `json:"path"`
			PriceImpact string `json:"priceImpact"`
		} `json:"route"`
		Bridge string `json:"bridge"`
	} `json:"srcChain"`
	BridgeChain struct {
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
		Route          []struct {
			AmountIn  string        `json:"amountIn"`
			AmountOut string        `json:"amountOut"`
			DexName   string        `json:"dexName"`
			Path      []interface{} `json:"path"`
		} `json:"route"`
		Bridge string `json:"bridge"`
	} `json:"bridgeChain"`
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
		Route          []struct {
			AmountIn  string `json:"amountIn"`
			AmountOut string `json:"amountOut"`
			DexName   string `json:"dexName"`
			Path      []struct {
				Id      string `json:"id"`
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
				Fee string `json:"fee"`
			} `json:"path"`
			PriceImpact string `json:"priceImpact"`
		} `json:"route"`
		Bridge string `json:"bridge"`
	} `json:"dstChain"`
	MinAmountOut struct {
		Amount string `json:"amount"`
		Symbol string `json:"symbol"`
	} `json:"minAmountOut"`
}

func Route(request *RouteRequest) ([]*RouteResponseData, error) {
	params, err := uhttp.URLEncode(request)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%s/%s?%s", Domain, PathRoute, params)
	ret, err := uhttp.Get(url, nil, nil)
	if err != nil {
		return nil, reqerror.NewExternalRequestError(
			url,
			reqerror.WithError(err),
		)
	}
	response := RouteResponse{}
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
