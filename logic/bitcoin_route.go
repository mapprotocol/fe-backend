package logic

import "github.com/mapprotocol/fe-backend/constants"

type Route struct {
	Hash   string `json:"hash"`
	GasFee struct {
		Amount string `json:"amount"`
		Symbol string `json:"symbol"`
	} `json:"gasFee"`
	SrcChain struct {
		ChainId        string `json:"chainId"`
		TokenAmountIn  string `json:"tokenAmountIn"`
		TokenAmountOut string `json:"tokenAmountOut"`
		Route          []struct {
			DexName string `json:"dexName"`
			Path    []struct {
				TokenIn struct {
					Address  string `json:"address"`
					Name     string `json:"name"`
					Symbol   string `json:"symbol"`
					Image    string `json:"image"`
					Decimals int    `json:"decimals"`
				} `json:"tokenIn"`
				TokenOut struct {
					Address  string `json:"address"`
					Name     string `json:"name"`
					Symbol   string `json:"symbol"`
					Image    string `json:"image"`
					Decimals int    `json:"decimals"`
				} `json:"tokenOut"`
			} `json:"path"`
		} `json:"route"`
	} `json:"srcChain"`
}

func GetBitcoinLocalRoutes(amount string) *Route {
	token := struct {
		Address  string `json:"address"`
		Name     string `json:"name"`
		Symbol   string `json:"symbol"`
		Image    string `json:"image"`
		Decimals int    `json:"decimals"`
	}{
		Address:  constants.BTCTokenAddress,
		Name:     "Bitcoin",
		Symbol:   "BTC",
		Image:    "https://map-static-file.s3.amazonaws.com/mapSwap/merlin/0x0000000000000000000000000000000000000000.jpg", // todo
		Decimals: 8,
	}
	route := Route{
		Hash: constants.LocalRouteHash,
		GasFee: struct {
			Amount string `json:"amount"`
			Symbol string `json:"symbol"`
		}{
			Amount: "0.00007614",
			Symbol: "BTC",
		},
		SrcChain: struct {
			ChainId        string `json:"chainId"`
			TokenAmountIn  string `json:"tokenAmountIn"`
			TokenAmountOut string `json:"tokenAmountOut"`
			Route          []struct {
				DexName string `json:"dexName"`
				Path    []struct {
					TokenIn struct {
						Address  string `json:"address"`
						Name     string `json:"name"`
						Symbol   string `json:"symbol"`
						Image    string `json:"image"`
						Decimals int    `json:"decimals"`
					} `json:"tokenIn"`
					TokenOut struct {
						Address  string `json:"address"`
						Name     string `json:"name"`
						Symbol   string `json:"symbol"`
						Image    string `json:"image"`
						Decimals int    `json:"decimals"`
					} `json:"tokenOut"`
				} `json:"path"`
			} `json:"route"`
		}{
			ChainId:        constants.BTCChainID,
			TokenAmountIn:  amount,
			TokenAmountOut: amount,
			Route: []struct {
				DexName string `json:"dexName"`
				Path    []struct {
					TokenIn struct {
						Address  string `json:"address"`
						Name     string `json:"name"`
						Symbol   string `json:"symbol"`
						Image    string `json:"image"`
						Decimals int    `json:"decimals"`
					} `json:"tokenIn"`
					TokenOut struct {
						Address  string `json:"address"`
						Name     string `json:"name"`
						Symbol   string `json:"symbol"`
						Image    string `json:"image"`
						Decimals int    `json:"decimals"`
					} `json:"tokenOut"`
				} `json:"path"`
			}{
				{
					DexName: constants.ExchangeNameFlushExchange,
					Path: []struct {
						TokenIn struct {
							Address  string `json:"address"`
							Name     string `json:"name"`
							Symbol   string `json:"symbol"`
							Image    string `json:"image"`
							Decimals int    `json:"decimals"`
						} `json:"tokenIn"`
						TokenOut struct {
							Address  string `json:"address"`
							Name     string `json:"name"`
							Symbol   string `json:"symbol"`
							Image    string `json:"image"`
							Decimals int    `json:"decimals"`
						} `json:"tokenOut"`
					}{
						{
							TokenIn:  token,
							TokenOut: token,
						},
					},
				},
			},
		},
	}
	return &route
}
