package constants

const (
	ChainIDOfChainPool = "56"
	ChainIDOfEthereum  = "1"
)

const (
	USDTDecimalOfChainPool = 1e18
	USDTDecimalOfEthereum  = 1e6
)

const (
	WBTCDecimalOfChainPool = 1e18
	WBTCDecimalOfEthereum  = 1e8
)

const BTCDecimal = 1e8

const (
	WBTCOfChainPool = "0x7130d2A12B9BCbFAe4f2634d864A1Ee1Ce3Ead9c"
	WBTCOfEthereum  = "0x2260FAC5E5542a773Aa44fBCfeDf7C193bc2C599"
	USDTOfChainPool = "0x55d398326f99059fF775485246999027B3197955"
	USDTOfEthereum  = "0xdac17f958d2ee523a2206206994597c13d831ec7"
	USDTOfTON       = "EQCxE6mUtQJKFnGfaROTKOt1lZbDiiX1kCixRv7Nw2Id_sDs"
)

const (
	TONChainID      = "1360104473493505"
	BTCChainID      = "313230561203979757" // common.BytesToAddress([]byte("BTCChainID")).Big().String()[:18]
	BTCTokenAddress = "0x0000000000000000000000000000000000000000"
)

const (
	SlippageMin = 300
	SlippageMax = 5000
)

const (
	ExchangeNameButter        = "Butter"
	ExchangeNameFlushExchange = "mapoX"
)

const USDTSymbol = "USDT"
const WBTCSymbol = "BTCB"

// LocalRouteBitcoinHash is the hash of local route for bitcoin chain
// This route indicates that the transaction was sent from Bitcoin
const LocalRouteBitcoinHash = "0x0000000000000000000000000000000000000000000000313230561203979757"
