package constants

const (
	ChainIDOfChainPool = "137" // todo
	ChainIDOfEthereum  = "1"
)

const (
	USDTDecimalOfChainPool = 1e6
	USDTDecimalOfEthereum  = 1e6
)

const WBTCOfChainPoll = "0x0d500B1d8E8eF31E21C99d1Db9A6444d3ADf1270" // todo

const (
	USDTOfChainPoll = "0xc2132D05D31c914a87C6611C10748AEb04B58e8F" // todo
	USDTOfEthereum  = "0xdac17f958d2ee523a2206206994597c13d831ec7" // todo
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

const BridgeFeeSymbol = "USDT"
