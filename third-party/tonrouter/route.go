package tonrouter

import (
	"encoding/json"
	"errors"
	"fmt"
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
	PathRoute       = "/v2/route"
	PathBridgeRoute = "/v2/route/bridge"
	PathGetRoute    = "/v2/route/hash"
	PathBridgeSwap  = "/v2/swap/bridge"
	PathBalance     = "/v2/jetton/router/balance"
)

var endpoint string

type RouteRequest struct {
	TokenInAddress  string `json:"tokenInAddress"`
	TokenOutAddress string `json:"tokenOutAddress"`
	Amount          string `json:"amount"`
	Slippage        uint64 `json:"slippage"`
}

type RouteResponse struct {
	Errno   int          `json:"errno"`
	Message string       `json:"message"`
	Data    []*RouteData `json:"data"`
}

type BridgeRouteRequest struct {
	ToChainID       string `json:"toChainId"`
	TokenInAddress  string `json:"tokenInAddress"`
	TokenOutAddress string `json:"tokenOutAddress"`
	Amount          string `json:"amount"`
	TonSlippage     uint64 `json:"tonSlippage"`
	Slippage        uint64 `json:"slippage"`
}

type BridgeRouteResponse struct {
	Errno   int          `json:"errno"`
	Message string       `json:"message"`
	Data    []*RouteData `json:"data"`
}

type GetRouteResponse struct {
	Errno   int    `json:"errno"`
	Message string `json:"message"`
	Data    struct {
		SrcChain struct {
			TokenAmountOut string `json:"tokenAmountOut"`
		} `json:"srcChain"`
	} `json:"data"`
}

type RouteData struct {
	Diff      string `json:"diff"`
	Hash      string `json:"hash"`
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

func Init() {
	cfg := viper.GetStringMapString("ton")
	endpoint = cfg["endpoint"]

	if utils.IsEmpty(endpoint) {
		panic("ton router endpoint is empty")
	}
}

func Route(request *RouteRequest) (*RouteData, error) {
	params := fmt.Sprintf(
		"tokenInAddress=%s&tokenOutAddress=%s&amount=%s&slippage=%d",
		request.TokenInAddress, request.TokenOutAddress, request.Amount, request.Slippage,
	)
	url := fmt.Sprintf("%s%s?%s", endpoint, PathRoute, params)
	log.Logger().Debugf("ton route url: %s", url)
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
		)
	}
	return response.Data[0], nil
}

func BridgeRoute(request *BridgeRouteRequest) (*RouteData, error) {
	params := fmt.Sprintf(
		"toChainId=%s&tokenInAddress=%s&tokenOutAddress=%s&amount=%s&tonSlippage=%d&slippage=%d",
		request.ToChainID, request.TokenInAddress, request.TokenOutAddress, request.Amount, request.TonSlippage, request.Slippage,
	)
	url := fmt.Sprintf("%s%s?%s", endpoint, PathBridgeRoute, params)
	log.Logger().Debugf("ton bridge route url: %s", url)
	ret, err := uhttp.Get(url, nil, nil)
	if err != nil {
		return nil, reqerror.NewExternalRequestError(
			url,
			reqerror.WithError(err),
		)
	}
	response := BridgeRouteResponse{}
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
		)
	}
	return response.Data[0], nil
}

func GetRouteAmountOut(hash string) (*big.Float, error) {
	url := fmt.Sprintf("%s%s/%s", endpoint, PathGetRoute, hash)
	fmt.Println("============================== url: ", url)
	log.Logger().Debugf("ton get route amount out url: %s", url)
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
	if utils.IsEmpty(response.Data.SrcChain.TokenAmountOut) {
		return nil, reqerror.NewExternalRequestError(
			url,
			reqerror.WithError(errors.New("token amount out is empty")),
		)
	}

	amountOut, ok := new(big.Float).SetString(response.Data.SrcChain.TokenAmountOut)
	if !ok {
		return nil, reqerror.NewExternalRequestError(
			url,
			reqerror.WithError(fmt.Errorf("invalid token amount out: %s", response.Data.SrcChain.TokenAmountOut)),
		)
	}
	return amountOut, nil
}
