package entity

type RouteRequest struct {
	FromChainID     string `form:"fromChainId"`
	ToChainID       string `form:"toChainId"`
	Amount          string `form:"amount"`
	TokenInAddress  string `form:"tokenInAddress"`
	TokenOutAddress string `form:"tokenOutAddress"`
	Type            string `form:"type"`
	Slippage        string `form:"slippage"`
	Action          uint8  `form:"action"`
}

type SwapRequest struct {
	SrcChain string `form:"srcChain"`
	SrcToken string `form:"srcToken"`
	Sender   string `form:"sender"`
	Amount   string `form:"amount"`
	DstChain string `form:"dstChain"`
	DstToken string `form:"dstToken"`
	Receiver string `form:"receiver"`
	Hash     string `form:"hash"`
	Slippage string `form:"slippage"`
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
