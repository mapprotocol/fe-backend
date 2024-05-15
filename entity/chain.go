package entity

type SupportedChainsRequest struct {
	Page int `form:"page"`
	Size int `form:"size"`
}

// response

type SupportedChainsResponse struct {
	ChainID   uint64 `json:"chain_id"`
	ChainName string `json:"chain_name"`
	ChainIcon string `json:"chain_icon"`
}
