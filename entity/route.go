package entity

type RouteRequest struct {
	FromChainID     string `form:"fromChainId"`
	TokenInAddress  string `form:"tokenInAddress"`
	Amount          string `form:"amount"`
	ToChainID       string `form:"toChainId"`
	TokenOutAddress string `form:"tokenOutAddress"`
	FeeCollector    string `form:"feeCollector"`
	FeeRatio        string `form:"feeRatio"`
	Type            string `form:"type"`
	Slippage        string `form:"slippage"`
	Action          uint8  `form:"action"`
}

type RouteResponse struct {
	Hash        string `json:"hash"`
	TokenIn     Token  `json:"tokenIn"`
	TokenOut    Token  `json:"tokenOut"`
	AmountIn    string `json:"amountIn"`
	AmountOut   string `json:"amountOut"`
	Path        []Path `json:"path"`
	GasFee      Fee    `json:"gasFee"`
	BridgeFee   Fee    `json:"bridgeFee"`
	ProtocolFee Fee    `json:"protocolFee"`
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

type Token struct {
	ChainId  string `json:"chainId"`
	Address  string `json:"address"`
	Name     string `json:"name"`
	Decimals int    `json:"decimals"`
	Symbol   string `json:"symbol"`
	Icon     string `json:"icon"`
}

type Fee struct {
	Amount  string `json:"amount"`
	Symbol  string `json:"symbol"`
	ChainId string `json:"chainId"`
}

type Path struct {
	Name      string `json:"name"`
	AmountIn  string `json:"amountIn"`
	AmountOut string `json:"amountOut"`
	TokenIn   Token  `json:"tokenIn"`
	TokenOut  Token  `json:"tokenOut"`
}
