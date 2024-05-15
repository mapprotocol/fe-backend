package entity

type SupportedTokensRequest struct {
	ChainID uint64 `form:"chain_id"`
	Symbol  string `form:"symbol"`
	Page    int    `form:"page"`
	Size    int    `form:"size"`
}

type DepositAddressRequest struct {
	ChainID     uint64 `json:"chain_id"`
	TokenSymbol string `json:"token_symbol"`
}

// response
type SupportedTokensResponse struct {
	ChainID  uint64 `json:"chain_id"`
	Name     string `json:"name"`
	Symbol   string `json:"symbol"`
	Decimals uint32 `json:"decimals"`
	Icon     string `json:"icon"`
}

type DepositAddressResponse struct {
	Address string `json:"address"`
}
