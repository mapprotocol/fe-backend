package err

import "errors"

var (
	LowBalanceInHotWallet1 = errors.New("not enough balance in the hot-wallet1")
	LowBalanceInHotWallet2 = errors.New("not enough balance in the hot-wallet2")
	LowFeeInHotWalletFee1  = errors.New("not enough fees in the hot-wallet-fee1")
	LowFeeInHotWalletFee2  = errors.New("not enough fees in the hot-wallet-fee2")
	LowFeeInHotWalletFee3  = errors.New("not enough fees in the hot-wallet-fee3")
)
