package ceffu

const Domain = ""
const SuccessCode = "000000"

const (
	PathCreateSubWallet   = "/open-api/v1/subwallet/create"
	PathGetDepositAddress = "/open-api/v1/subwallet/deposit/address"
	PathDepositHistory    = "/open-api/v2/subwallet/deposit/history"
	PathTransfer          = "/open-api/v1/subwallet/transfer"
	PathWithdrawal        = "/open-api/v2/wallet/withdrawal"
	PathWithdrawalDetail  = "/open-api/v2/wallet/withdrawal/detail"
)
