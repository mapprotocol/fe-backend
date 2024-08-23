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
	WBTCOfChainPool = "0xb877e3562a660c7861117c2f1361a26abaf19beb"
	WBTCOfEthereum  = "0x2260FAC5E5542a773Aa44fBCfeDf7C193bc2C599"
)

const (
	WBTCDecimalOfChainPool uint64 = 1e18
	WBTCDecimalOfEthereum  uint64 = 1e8
)
