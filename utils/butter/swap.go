package butter

import (
	"encoding/json"
	"fmt"
	uhttp "github.com/mapprotocol/ceffu-fe-backend/utils/http"
	"github.com/mapprotocol/ceffu-fe-backend/utils/reqerror"
	"strconv"
)

const Domain = "https://bs-router-test.chainservice.io"

const SuccessCode = 0

const PathRouteAndSwap = "/routeAndSwap"

type RouterAndSwapRequest struct {
	FromChainID     string `json:"fromChainId"`
	ToChainID       string `json:"toChainId"`
	Amount          string `json:"amount"`
	TokenInAddress  string `json:"tokenInAddress"`
	TokenOutAddress string `json:"tokenOutAddress"`
	Kind            string `json:"type"`
	Slippage        string `json:"slippage"`
	Entrance        string `json:"entrance"`
	From            string `json:"from"`
	Receiver        string `json:"receiver"`
}

type RouterAndSwapResponse struct {
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

func RouterAndSwap(request *RouterAndSwapRequest) (*TxData, error) {
	params, err := uhttp.URLEncode(request)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%s/%s?%s", Domain, PathRouteAndSwap, params)
	fmt.Println("============================== url:", url)
	ret, err := uhttp.Get(url, nil, nil)
	if err != nil {
		return nil, reqerror.NewExternalRequestError(
			url,
			reqerror.WithError(err),
		)
	}
	response := RouterAndSwapResponse{}
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
	if response.Data[0].TxParam.Errno != SuccessCode {
		return nil, reqerror.NewExternalRequestError(
			url,
			reqerror.WithCode(strconv.Itoa(response.Errno)),
			reqerror.WithMessage(response.Message),
		)
	}
	return response.Data[0].TxParam.Data[0], nil
}
