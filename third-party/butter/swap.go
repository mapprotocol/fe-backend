package butter

import (
	"encoding/json"
	"fmt"
	"github.com/mapprotocol/fe-backend/resource/log"
	uhttp "github.com/mapprotocol/fe-backend/utils/http"
	"github.com/mapprotocol/fe-backend/utils/reqerror"
	"github.com/spf13/viper"
	"strconv"
)

const SuccessCode = 0

const (
	PathRouteAndSwap   = "/routeAndSwap"
	PathRoute          = "/route"
	PathEvmCrossInSwap = "/evmCrossInSwap"
)

var Domain string
var entrance string

type RouterAndSwapRequest struct {
	FromChainID     string `json:"fromChainId"`
	ToChainID       string `json:"toChainId"`
	Amount          string `json:"amount"`
	TokenInAddress  string `json:"tokenInAddress"`
	TokenOutAddress string `json:"tokenOutAddress"`
	Type            string `json:"type"`
	Slippage        uint64 `json:"slippage"`
	From            string `json:"from"`
	Receiver        string `json:"receiver"`
	Referrer        string `json:"referrer"`
	RateOrNativeFee string `json:"rateOrNativeFee"`
}

type RouterRequest struct {
	FromChainID     string `json:"fromChainId"`
	ToChainID       string `json:"toChainId"`
	Amount          string `json:"amount"`
	TokenInAddress  string `json:"tokenInAddress"`
	TokenOutAddress string `json:"tokenOutAddress"`
	Type            string `json:"type"`
	Slippage        uint64 `json:"slippage"`
	Entrance        string `json:"entrance"`
}

type EvmCrossInSwapRequest struct {
	Hash         string `json:"hash"`
	SrcChainId   string `json:"srcChainId"`
	From         string `json:"from"`
	Router       string `json:"router"`   // 签名地址
	Receiver     string `json:"receiver"` // 接受者
	MinAmountOut string `json:"minAmountOut"`
	OrderIdHex   string `json:"orderIdHex"`
	Fee          string `json:"fee"`
	FeeReceiver  string `json:"feeReceiver"`
}

type RouterAndSwapResponse struct {
	Errno   int    `json:"errno"`
	Message string `json:"message"`
	Data    []struct {
		Route struct {
			MinAmountOut struct {
				Amount string `json:"amount"`
				Symbol string `json:"symbol"`
			} `json:"minAmountOut"`
		}
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

type RouteData struct {
	Errno   int    `json:"errno"`
	Message string `json:"message"`
	Data    []struct {
		Diff      string `json:"diff"`
		BridgeFee struct {
			Amount string `json:"amount"`
		} `json:"bridgeFee"`
		TradeType int `json:"tradeType"`
		GasFee    struct {
			Amount string `json:"amount"`
			Symbol string `json:"symbol"`
		} `json:"gasFee"`
		SwapFee struct {
			NativeFee string `json:"nativeFee"`
			TokenFee  string `json:"tokenFee"`
		} `json:"swapFee"`
		FeeConfig struct {
			FeeType         int    `json:"feeType"`
			Referrer        string `json:"referrer"`
			RateOrNativeFee int    `json:"rateOrNativeFee"`
		} `json:"feeConfig"`
		GasEstimated       string `json:"gasEstimated"`
		GasEstimatedTarget string `json:"gasEstimatedTarget"`
		TimeEstimated      int    `json:"timeEstimated"`
		Hash               string `json:"hash"`
		Entrance           string `json:"entrance"`
		Timestamp          int64  `json:"timestamp"`
		HasLiquidity       bool   `json:"hasLiquidity"`
		SrcChain           struct {
			ChainID string `json:"chainId"`
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
		} `json:"srcChain"`
		Contract     string `json:"contract"`
		MinAmountOut struct {
			Amount string `json:"amount"`
			Symbol string `json:"symbol"`
		} `json:"minAmountOut"`
	} `json:"data"`
}

type EvmCrossInSwapData struct {
	Errno   int    `json:"errno"`
	Message string `json:"message"`
	Data    []struct {
		To      string `json:"to"`
		Data    string `json:"data"`
		Value   string `json:"value"`
		ChainID string `json:"chainId"`
		Method  string `json:"method"`
		Args    []struct {
			Type  string `json:"type"`
			Value struct {
				OrderID  string `json:"orderId"`
				Receiver string `json:"receiver"`
				Token    string `json:"token"`
				Amount   struct {
					Type string `json:"type"`
					Hex  string `json:"hex"`
				} `json:"amount"`
				FromChain   int64  `json:"fromChain"`
				ToChain     string `json:"toChain"`
				Fee         string `json:"fee"`
				FeeReceiver string `json:"feeReceiver"`
				From        string `json:"from"`
				ButterData  string `json:"butterData"`
			} `json:"value"`
		} `json:"args"`
	} `json:"data"`
}

func Init() {
	Domain = viper.GetStringMapString("endpoints")["butter"]
	entrance = viper.GetStringMapString("butter")["entrance"]
}

func RouteAndSwap(req *RouterAndSwapRequest) (*TxData, error) {
	params := fmt.Sprintf(
		"fromChainId=%s&toChainId=%s&amount=%s&tokenInAddress=%s&tokenOutAddress=%s&type=%s&slippage=%d&entrance=%s&from=%s&receiver=%s",
		req.FromChainID, req.ToChainID, req.Amount, req.TokenInAddress, req.TokenOutAddress, req.Type, req.Slippage, entrance, req.From, req.Receiver,
	)
	url := fmt.Sprintf("%s%s?%s", Domain, PathRouteAndSwap, params)
	log.Logger().Debug(fmt.Sprintf("route and swap url: %s", url))
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

func RouteAndSwapSol(req *RouterAndSwapRequest) (*RouterAndSwapResponse, error) {
	params := fmt.Sprintf(
		"fromChainId=%s&toChainId=%s&amount=%s&tokenInAddress=%s&tokenOutAddress=%s&type=%s&slippage=%d&from=%s&receiver=%s",
		req.FromChainID, req.ToChainID, req.Amount, req.TokenInAddress, req.TokenOutAddress, req.Type, req.Slippage, req.From, req.Receiver,
	)
	if req.Referrer != "" {
		params = fmt.Sprintf("%s&referrer=%s", params, req.Referrer)
	}
	if req.RateOrNativeFee != "" {
		params = fmt.Sprintf("%s&rateOrNativeFee=%s", params, req.RateOrNativeFee)
	}
	url := fmt.Sprintf("%s%s?%s", Domain, PathRouteAndSwap, params)
	//fmt.Println("url ------------ ", url)
	log.Logger().Debug(fmt.Sprintf("route and swap url: %s", url))
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
	return &response, nil
}

func Route(req *RouterRequest) (*RouteData, error) {
	params := fmt.Sprintf(
		"fromChainId=%s&toChainId=%s&amount=%s&tokenInAddress=%s&tokenOutAddress=%s&type=%s&slippage=%d&entrance=%s",
		req.FromChainID, req.ToChainID, req.Amount, req.TokenInAddress, req.TokenOutAddress, req.Type, req.Slippage, entrance,
	)
	url := fmt.Sprintf("%s%s?%s", Domain, PathRoute, params)
	log.Logger().Debug(fmt.Sprintf("route and swap url: %s", url))
	ret, err := uhttp.Get(url, nil, nil)
	if err != nil {
		return nil, reqerror.NewExternalRequestError(
			url,
			reqerror.WithError(err),
		)
	}
	response := RouteData{}
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
			reqerror.WithMessage("route back zero"),
		)
	}
	return &response, nil
}

func EvmCrossInSwap(req *EvmCrossInSwapRequest) (*EvmCrossInSwapData, error) {
	params := fmt.Sprintf(
		"hash=%s&srcChainId=%s&from=%s&router=%s&receiver=%s&minAmountOut=%s&orderIdHex=%s&fee=%s&feeReceiver=%s",
		req.Hash, req.SrcChainId, req.From, req.Router, req.Receiver, req.MinAmountOut, req.OrderIdHex, req.Fee, req.FeeReceiver,
	)
	url := fmt.Sprintf("%s%s?%s", Domain, PathEvmCrossInSwap, params)
	log.Logger().Debug(fmt.Sprintf("route and swap url: %s", url))
	ret, err := uhttp.Get(url, nil, nil)
	if err != nil {
		return nil, reqerror.NewExternalRequestError(
			url,
			reqerror.WithError(err),
		)
	}
	response := EvmCrossInSwapData{}
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
			reqerror.WithMessage("evmCrossInSwap back zero"),
		)
	}
	return &response, nil
}
