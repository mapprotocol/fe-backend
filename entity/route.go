package entity

import "github.com/mapprotocol/fe-backend/third-party/butter"

type RouteRequest struct {
	FromChainID     string `form:"fromChainId"`
	ToChainID       string `form:"toChainId"`
	Amount          string `form:"amount"`
	TokenInAddress  string `form:"tokenInAddress"`
	TokenOutAddress string `form:"tokenOutAddress"`
	Kind            string `form:"type"`
	Slippage        string `form:"slippage"`
}

type SwapRequest struct {
	Hash     string `form:"hash"`
	Slippage string `form:"slippage"`
	From     string `form:"from"`
	Receiver string `form:"receiver"`
}

type RouteResponse struct {
	Route       []Route                     `json:"route"`
	ButterRoute []*butter.RouteResponseData `json:"butterRoute"`
}

type SwapResponse struct {
	To      string `json:"to"`
	Data    string `json:"data"`
	Value   string `json:"value"`
	ChainId string `json:"chainId"`
}

type Route struct {
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
}
