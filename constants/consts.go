package constants

const (
	ChainIDOfChainPool = "22776"
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

const USDTDecimalNumberOfChainPool = 18

const (
	WBTCOfChainPool = "0xb877e3562a660c7861117c2f1361a26abaf19beb"
	WBTCOfEthereum  = "0x2260FAC5E5542a773Aa44fBCfeDf7C193bc2C599"
	USDTOfChainPool = "0x33daba9618a75a7aff103e53afe530fbacf4a3dd"
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

const BridgeFeeSymbolOfTON = "USDT"
const BridgeFeeSymbolOfBTC = "BTC"

const NativeSymbolOfChainPool = "MAPO"

const LocalRouteGasFee = "0.121116"
const LocalRouteHash = "0x0000000000000000000000000000000000000000000000000000000000022776"
