package params

const (
	BTCChainID      = "313230561203979757" // common.BytesToAddress([]byte("BTCChainID")).Big().String()[:18]
	BTCTokenAddress = "0x0000000000000000000000000000000000000000"
)

const (
	BTCDecimal = 1e8
)

// todo move to chain_pool.go
const (
	WBTCOfChainPool = "0x0d500B1d8E8eF31E21C99d1Db9A6444d3ADf1270" // todo
	WBTCOfEthereum  = ""                                           // todo
)

const (
	WBTCDecimalOfChainPool = 1e18
	WBTCDecimalOfEthereum  = 1e8
)
