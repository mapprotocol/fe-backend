package entity

type SupportedTokensRequest struct {
	ChainID uint64 `form:"chain_id"`
	Symbol  string `form:"symbol"`
	Page    int    `form:"page"`
	Size    int    `form:"size"`
}

type CreateOrderRequest struct {
	SrcChain uint64 `json:"src_chain"`
	SrcToken string `json:"src_token"`
	Sender   string `json:"sender"`
	Amount   string `json:"amount"`
	DstChain uint64 `json:"dst_chain"`
	DstToken string `json:"dst_token"`
	Receiver string `json:"receiver"`
}

// response

type SupportedTokensResponse struct {
	ChainID  uint64 `json:"chain_id"`
	Name     string `json:"name"`
	Symbol   string `json:"symbol"`
	Decimals uint32 `json:"decimals"`
	Icon     string `json:"icon"`
}

type CreateOrderResponse struct {
	OrderID        uint64 `json:"order_id"`
	DepositAddress string `json:"deposit_address"`
}
