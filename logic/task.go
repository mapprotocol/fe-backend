package logic

import (
	"github.com/mapprotocol/fe-backend/resource/log"
	"time"
)

const (
	FeeRateMultiple = 2
	FeeRateLimit    = 600
)

var globalFeeRate int64 = 20

func setGlobalFeeRate(feeRate int64) {
	globalFeeRate = feeRate
}

func GetGlobalFeeRate() int64 {
	if globalFeeRate == 0 {
		return 20
	}
	return globalFeeRate
}

func GetFeeRate() {
	for {
		fees, err := btcApiClient.RecommendedFees()
		if err != nil {
			log.Logger().WithField("error", err).Error("failed to get fee rate")
			time.Sleep(10 * time.Second)
			continue
		}

		feeRate := fees.FastestFee * FeeRateMultiple
		if feeRate > FeeRateLimit {
			feeRate = FeeRateLimit
		}

		setGlobalFeeRate(feeRate)

		log.Logger().WithField("fastestFee", fees.FastestFee).WithField("feeRate", globalFeeRate).Info("got fee rate")
		time.Sleep(30 * time.Minute)
		continue
	}
}
