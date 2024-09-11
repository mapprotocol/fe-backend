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
	WBTCOfChainPool = "0x1bfd67037b42cf73acf2047067bd4f2c47d9bfd6"
	WBTCOfEthereum  = "0x2260FAC5E5542a773Aa44fBCfeDf7C193bc2C599"
)

const (
	WBTCDecimalOfChainPool uint64 = 1e8
	WBTCDecimalOfEthereum  uint64 = 1e8
)

const (
	BTCDecimalNumber             = 8
	WBTCDecimalNumberOfChainPool = 8
	WBTCDecimalNumberOfEthereum  = 8
)
