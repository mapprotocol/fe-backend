package entity

type OrderListRequest struct {
	Sender string `form:"sender"`
	Page   int    `form:"page"`
	Size   int    `form:"size"`
}

type OrderDetailRequest struct {
	OrderID uint64 `form:"order_id"`
}

// response

type OrderListResponse struct {
	OrderID        uint64 `json:"order_id"`
	SrcChain       uint64 `json:"src_chain"`
	SrcToken       string `json:"src_token"`
	Sender         string `json:"sender"`
	Amount         string `json:"amount"`
	DstChain       uint64 `json:"dst_chain"`
	DstToken       string `json:"dst_token"`
	Receiver       string `json:"receiver"`
	DepositAddress string `json:"deposit_address"`
	TxHash         string `json:"tx_hash"`
	Action         uint8  `json:"action"`
	Status         uint8  `json:"status"`
	CreatedAt      int64  `json:"created_at"`
}

type OrderDetailResponse struct {
	OrderID        uint64 `json:"order_id"`
	SrcChain       uint64 `json:"src_chain"`
	SrcToken       string `json:"src_token"`
	Sender         string `json:"sender"`
	Amount         string `json:"amount"`
	DstChain       uint64 `json:"dst_chain"`
	DstToken       string `json:"dst_token"`
	Receiver       string `json:"receiver"`
	DepositAddress string `json:"deposit_address"`
	TxHash         string `json:"tx_hash"`
	Action         uint8  `json:"action"`
	Status         uint8  `json:"status"`
	CreatedAt      int64  `json:"created_at"`
}
