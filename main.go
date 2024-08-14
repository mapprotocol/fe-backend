package main

import (
	"context"
	"fmt"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/mapprotocol/fe-backend/config"
	"github.com/mapprotocol/fe-backend/logic"
	"github.com/mapprotocol/fe-backend/resource/db"
	"github.com/mapprotocol/fe-backend/resource/log"
	"github.com/mapprotocol/fe-backend/utils/alarm"
	"github.com/spf13/viper"
	"os"
)

const (
	VersionMajor = 0          // Major version component of the current release
	VersionMinor = 0          // Minor version component of the current release
	VersionPatch = 1          // Patch version component of the current release
	VersionMeta  = "unstable" // Version metadata to append to the version string
)

func version() string {
	return fmt.Sprintf("%d.%d.%d-%s", VersionMajor, VersionMinor, VersionPatch, VersionMeta)
}

func toConfig() (*logic.CollectCfg, error) {
	cfg := &logic.CollectCfg{}
	cfg.Testnet = viper.GetBool("testnet")

	network := &chaincfg.MainNetParams
	if cfg.Testnet {
		network = &chaincfg.TestNet3Params
	}
	cfg.StrHotWallet1Priv = viper.GetString("hotwalletPriv1")
	cfg.StrHotWallet2Priv = viper.GetString("hotwalletPriv2")

	cfg.StrHotWalletFee1Privkey = viper.GetString("hotwalletFeePriv1")
	cfg.StrHotWalletFee2Privkey = viper.GetString("hotwalletFeePriv2")
	cfg.StrFee3Privkey = viper.GetString("hotwalletFeePriv3")

	strHotAddr1 := viper.GetString("hotwalletAddress1")
	hotAddr1, err := btcutil.DecodeAddress(strHotAddr1, network)
	if err != nil {
		log.Logger().WithField("error", err).Error("decode hot1 address failed")
		return cfg, err
	}
	cfg.HotWallet1 = hotAddr1

	strHotAddr2 := viper.GetString("hotwalletAddress2")
	hotAddr2, err := btcutil.DecodeAddress(strHotAddr2, network)
	if err != nil {
		log.Logger().WithField("error", err).Error("decode hot2 address failed")
		return cfg, err
	}
	cfg.HotWallet2 = hotAddr2

	strFeeAddr1 := viper.GetString("hotwalletFeeAddress1")
	feeAddr1, err := btcutil.DecodeAddress(strFeeAddr1, network)
	if err != nil {
		log.Logger().WithField("error", err).Error("decode fee1 address failed")
		return cfg, err
	}
	cfg.HotWalletFee1 = feeAddr1

	strFeeAddr2 := viper.GetString("hotwalletFeeAddress2")
	feeAddr2, err := btcutil.DecodeAddress(strFeeAddr2, network)
	if err != nil {
		log.Logger().WithField("error", err).Error("decode fee2 address failed")
		return cfg, err
	}
	cfg.HotWalletFee2 = feeAddr2

	strFeeAddr3 := viper.GetString("hotwalletFeeAddress3")
	feeAddr3, err := btcutil.DecodeAddress(strFeeAddr3, network)
	if err != nil {
		log.Logger().WithField("error", err).Error("decode fee3 address failed")
		return cfg, err
	}
	cfg.HotWalletFee3 = feeAddr3

	amount0 := viper.GetFloat64("minHotWallet2Amount")
	amount1, err := btcutil.NewAmount(amount0)
	if err != nil {
		return cfg, err
	}
	cfg.HotWallet2Line = int64(amount1)

	amount0 = viper.GetFloat64("maxTranferAmount")
	amount1, err = btcutil.NewAmount(amount0)
	if err != nil {
		return cfg, err
	}
	cfg.MaxTransferAmount = int64(amount1)

	return cfg, nil
}
func main() {
	args := os.Args
	if len(args) >= 2 && args[1] == "version" {
		fmt.Println("version:", version())
		return
	}

	alarm.ValidateEnv()
	// init config
	config.Init()
	// init log
	log.Init(viper.GetString("env"), viper.GetString("logDir"))
	// init db
	dbConf := viper.GetStringMapString("database")
	db.Init(dbConf["user"], dbConf["password"], dbConf["host"], dbConf["port"], dbConf["name"])

	cfg, err := toConfig()
	if err != nil {
		log.Logger().WithField("error", err).Error("read config failed")
		return
	}
	if err := logic.Run(cfg); err != nil {
		log.Logger().WithField("error", err).Error("collect failed")
		alarm.Slack(context.Background(), fmt.Sprintf("collect failed: %s", err.Error()))
		return
	}
}
