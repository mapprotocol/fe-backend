package butter

import (
	"encoding/json"
	"fmt"
	"github.com/mapprotocol/fe-backend/utils"
	uhttp "github.com/mapprotocol/fe-backend/utils/http"
	"github.com/mapprotocol/fe-backend/utils/reqerror"
	"github.com/spf13/viper"
	"strconv"
)

const SuccessCode = 0

const PathRoute = "/route"

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
		Bridge         string `json:"bridge"`
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
		Bridge         string `json:"bridge"`
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
		Bridge         string `json:"bridge"`
	} `json:"dstChain"`
	MinAmountOut struct {
		Amount string `json:"amount"`
		Symbol string `json:"symbol"`
	} `json:"minAmountOut"`
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
func Route(request *RouteRequest) ([]*RouteResponseData, error) { // todo check butter special code(code: 2003, message: No Route Found)
	params := fmt.Sprintf(
		"fromChainId=%s&toChainId=%s&tokenInAddress=%s&tokenOutAddress=%s&amount=%s&type=%s&slippage=%d&entrance=%s",
		request.FromChainID, request.ToChainID, request.TokenInAddress, request.TokenOutAddress, request.Amount, request.Type, request.Slippage, entrance,
	)
	url := fmt.Sprintf("%s%s?%s", endpoint, PathRoute, params)
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
	return response.Data, nil
}
