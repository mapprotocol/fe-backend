package butter

import (
	"encoding/json"
	"errors"
	"fmt"
	uhttp "github.com/mapprotocol/fe-backend/utils/http"
	"github.com/mapprotocol/fe-backend/utils/reqerror"
	"github.com/spf13/viper"
	"strconv"
)

const PathSwap = "/swap"
const PathRouteAndSwap = "/routeAndSwap"

var (
	ErrNotFoundTxData         = errors.New("not found tx data")
	ErrNotFoundRouteAndTxData = errors.New("not found route and tx data")
)

type SwapRequest struct {
	Hash     string `json:"hash"`
	Slippage uint64 `json:"slippage"`
	From     string `json:"from"`
	Receiver string `json:"receiver"`
	CallData string `json:"callData"`
}

type SwapResponse struct {
	Errno   int       `json:"errno"`
	Message string    `json:"message"`
	Data    []*TxData `json:"data"`
}

type RouteAndSwapRequest struct {
	FromChainID     string `json:"fromChainId"`
	ToChainID       string `json:"toChainId"`
	Amount          string `json:"amount"`
	TokenInAddress  string `json:"tokenInAddress"`
	TokenOutAddress string `json:"tokenOutAddress"`
	Kind            string `json:"type"`
	Slippage        uint64 `json:"slippage"`
	Entrance        string `json:"entrance"`
	From            string `json:"from"`
	Receiver        string `json:"receiver"`
}

type RouteAndSwapResponse struct {
	Errno   int    `json:"errno"`
	Message string `json:"message"`
	Data    []struct {
		Route struct {
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
							Decimals int    `json:"decimals"`
							Symbol   string `json:"symbol"`
							Icon     string `json:"icon"`
							Name     string `json:"name,omitempty"`
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
					AmountIn  string        `json:"amountIn"`
					AmountOut string        `json:"amountOut"`
					DexName   string        `json:"dexName"`
					Path      []interface{} `json:"path"`
					Extra     string        `json:"extra"`
				} `json:"route"`
				Bridge string `json:"bridge"`
			} `json:"dstChain"`
			MinAmountOut struct {
				Amount string `json:"amount"`
				Symbol string `json:"symbol"`
			} `json:"minAmountOut"`
		} `json:"route"`
		TxParam struct {
			Errno   int       `json:"errno"`
			Message string    `json:"message"`
			Data    []*TxData `json:"data"`
		} `json:"txParam"`
	} `json:"data"`
}

type TxData struct {
	To      string `json:"to"`
	Data    string `json:"data"`
	Value   string `json:"value"`
	ChainId string `json:"chainId"`
}

func Swap(request *SwapRequest) (*TxData, error) {
	params, err := uhttp.URLEncode(request) // todo
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%s%s?%s", Domain, PathSwap, params)
	fmt.Println("============================== swap url: ", url)
	ret, err := uhttp.Get(url, nil, nil)
	if err != nil {
		return nil, reqerror.NewExternalRequestError(
			url,
			reqerror.WithError(err),
		)
	}
	response := SwapResponse{}
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
	if len(response.Data) == 0 {
		return nil, reqerror.NewExternalRequestError(
			url,
			reqerror.WithCode(strconv.Itoa(response.Errno)),
			reqerror.WithMessage(response.Message),
			reqerror.WithError(ErrNotFoundTxData),
		)
	}
	return response.Data[0], nil
}

func RouteAndSwap(request *RouteAndSwapRequest) (*TxData, error) {
	request.Entrance = viper.GetStringMapString("butter")["entrance"]
	params, err := uhttp.URLEncode(request) // todo
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%s%s?%s", Domain, PathRouteAndSwap, params)
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
	if len(response.Data) == 0 {
		return nil, reqerror.NewExternalRequestError(
			url,
			reqerror.WithCode(strconv.Itoa(response.Errno)),
			reqerror.WithMessage(response.Message),
			reqerror.WithError(ErrNotFoundRouteAndTxData),
		)
	}
	if response.Data[0].TxParam.Errno != SuccessCode {
		return nil, reqerror.NewExternalRequestError(
			url,
			reqerror.WithCode(strconv.Itoa(response.Errno)),
			reqerror.WithMessage(response.Message),
		)
	}
	if len(response.Data[0].TxParam.Data) == 0 {
		return nil, reqerror.NewExternalRequestError(
			url,
			reqerror.WithCode(strconv.Itoa(response.Errno)),
			reqerror.WithMessage(response.Message),
			reqerror.WithError(ErrNotFoundTxData),
		)
	}
	return response.Data[0].TxParam.Data[0], nil
}
