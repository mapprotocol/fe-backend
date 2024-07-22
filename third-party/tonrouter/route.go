package tonrouter

import (
	"encoding/json"
	"fmt"
	uhttp "github.com/mapprotocol/fe-backend/utils/http"
	"github.com/mapprotocol/fe-backend/utils/reqerror"
	"strconv"
)

const Domain = "https://ton-router-test.chainservice.io"

const SuccessCode = 0

const (
	PathRoute        = "/route"
	PathRouteAndSwap = "/routeAndSwap"
)

type RouteRequest struct {
	FromChainID     string `json:"fromChainId"`
	ToChainID       string `json:"toChainId"`
	Amount          string `json:"amount"`
	TokenInAddress  string `json:"tokenInAddress"`
	TokenOutAddress string `json:"tokenOutAddress"`
	Slippage        string `json:"slippage"`
}

type RouteResponse struct {
	Errno   int    `json:"errno"`
	Message string `json:"message"`
	Data    struct {
		Route *RouteData `json:"route"`
	} `json:"data"`
}

type RouteData struct {
	Diff      string `json:"diff"`
	BridgeFee struct {
		Amount string `json:"amount"`
		Symbol string `json:"symbol"`
	} `json:"bridgeFee"`
	GasFee struct {
		Amount string `json:"amount"`
		Symbol string `json:"symbol"`
	} `json:"gasFee"`
	MinAmountOut struct {
		Amount string `json:"amount"`
		Symbol string `json:"symbol"`
	} `json:"minAmountOut"`
	SrcChain struct {
		ChainId        string `json:"chainId"`
		TokenAmountIn  string `json:"tokenAmountIn"`
		TokenAmountOut string `json:"tokenAmountOut"`
		Route          []struct {
			DexName string `json:"dexName"`
			Path    []struct {
				Fee     string `json:"fee"`
				Id      string `json:"id"`
				TokenIn struct {
					Type     string      `json:"type"`
					Address  string      `json:"address"`
					Name     string      `json:"name"`
					Symbol   string      `json:"symbol"`
					Image    string      `json:"image"`
					Decimals int         `json:"decimals"`
					Aliased  bool        `json:"aliased"`
					Price    string      `json:"price"`
					Source   interface{} `json:"source"`
				} `json:"tokenIn"`
				TokenOut struct {
					Type     string      `json:"type"`
					Address  string      `json:"address"`
					Name     string      `json:"name"`
					Symbol   string      `json:"symbol"`
					Image    string      `json:"image"`
					Decimals int         `json:"decimals"`
					Aliased  bool        `json:"aliased"`
					Price    string      `json:"price"`
					Source   interface{} `json:"source"`
				} `json:"tokenOut"`
			} `json:"path"`
		} `json:"route"`
	} `json:"srcChain"`
	Timestamp int64 `json:"timestamp"`
	TradeType int   `json:"tradeType"`
}

func Route(request *RouteRequest) (*RouteData, error) {
	params, err := uhttp.URLEncode(request) // todo
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%s%s?%s", Domain, PathRoute, params)
	fmt.Println("============================== route url: ", url)
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
	return response.Data.Route, nil
}
