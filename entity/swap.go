package entity

type SwapRequest struct {
	SrcChain     string `form:"srcChain"`
	SrcToken     string `form:"srcToken"`
	Sender       string `form:"sender"`
	Amount       string `form:"amount"`
	Decimal      uint8  `form:"decimal"`
	DstChain     string `form:"dstChain"`
	DstToken     string `form:"dstToken"`
	Receiver     string `form:"receiver"`
	FeeCollector string `form:"feeCollector"`
	FeeRatio     string `form:"feeRatio"`
	Hash         string `form:"hash"`
	Slippage     string `form:"slippage"`
}

type SwapResponse struct {
	To      string `json:"to"`
	Data    string `json:"data"`
	Value   string `json:"value"`
	ChainId string `json:"chainId"`
}
