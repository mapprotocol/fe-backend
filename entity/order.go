package entity

type CreateOrderRequest struct {
	SrcChain string `json:"srcChain"`
	SrcToken string `json:"srcToken"`
	Sender   string `json:"sender"`
	Amount   string `json:"amount"`
	DstChain string `json:"dstChain"`
	DstToken string `json:"dstToken"`
	Receiver string `json:"receiver"`
	Action   uint8  `json:"action"`
	Hash     string `json:"hash"`
	Slippage uint64 `json:"slippage"`
}

type UpdateOrderRequest struct {
	OrderID  uint64 `json:"orderId"`
	InTxHash string `json:"inTxHash"`
}

type OrderListRequest struct {
	Sender string `form:"sender"`
	Page   int    `form:"page"`
	Size   int    `form:"size"`
}

type OrderDetailRequest struct {
	Sender  string `form:"sender"`
	OrderID uint64 `form:"orderId"`
}

// response

type CreateOrderResponse struct {
	OrderID uint64 `json:"order_id"`
	Relayer string `json:"relayer"`
}

type OrderListResponse struct {
	OrderID   uint64 `json:"orderId"`
	SrcChain  string `json:"srcChain"`
	SrcToken  string `json:"srcToken"`
	Sender    string `json:"sender"`
	InAmount  string `json:"inAmount"`
	DstChain  string `json:"dstChain"`
	DstToken  string `json:"dstToken"`
	Receiver  string `json:"receiver"`
	OutAmount string `json:"outAmount"`
	Action    uint8  `json:"action"`
	Stage     uint8  `json:"stage"`
	Status    uint8  `json:"status"`
	CreatedAt int64  `json:"createdAt"`
}

type OrderDetailResponse struct {
	OrderID   uint64 `json:"orderId"`
	SrcChain  string `json:"srcChain"`
	SrcToken  string `json:"srcToken"`
	Sender    string `json:"sender"`
	InAmount  string `json:"inAmount"`
	DstChain  string `json:"dstChain"`
	DstToken  string `json:"dstToken"`
	Receiver  string `json:"receiver"`
	OutAmount string `json:"outAmount"`
	Action    uint8  `json:"action"`
	Stage     uint8  `json:"stage"`
	Status    uint8  `json:"status"`
	CreatedAt int64  `json:"createdAt"`
}
