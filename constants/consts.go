package constants

const (
	ChainIDOfChainPool = "137"
	ChainIDOfEthereum  = "1"
)

const (
	USDTDecimalOfChainPool = 1e6
	USDTDecimalOfEthereum  = 1e6
)

const (
	WBTCDecimalOfChainPool = 1e8
	WBTCDecimalOfEthereum  = 1e8
)

const USDTDecimalNumberOfChainPool = 6
const BTCDecimalNumberOfChainPool = 8

const BTCDecimal = 1e8

const (
	WBTCOfChainPool = "0x1bfd67037b42cf73acf2047067bd4f2c47d9bfd6"
	WBTCOfEthereum  = "0x2260FAC5E5542a773Aa44fBCfeDf7C193bc2C599"
	USDTOfChainPool = "0xc2132d05d31c914a87c6611c10748aeb04b58e8f"
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
	ExchangeNameFlushExchange = "MAP FE"
)

const USDTSymbol = "USDT"
const WBTCSymbol = "WBTC"

const NativeSymbolOfChainPool = "MATIC"

const LocalRouteGasFee = "0.03"

// LocalRouteHash is the hash of local route for chain pool
// This route indicates that the `deliver` method of the chain pool needs to be called to exchange
// the same token on the same chain.
const LocalRouteHash = "0x0000000000000000000000000000000000000000000000000000000000022776"

// LocalRouteBitcoinHash is the hash of local route for bitcoin chain
// This route indicates that the transaction was sent from Bitcoin
const LocalRouteBitcoinHash = "0x0000000000000000000000000000000000000000000000313230561203979757"
