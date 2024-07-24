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

	testnet := viper.GetBool("testnet")
	network := &chaincfg.MainNetParams
	if testnet {
		network = &chaincfg.TestNet3Params
	}
	feeAddress, err := btcutil.DecodeAddress(viper.GetString("feeaddress"), network)
	if err != nil {
		log.Logger().WithField("error", err).Error("decode fee address failed")
		return
	}
	receiverAddress, err := btcutil.DecodeAddress(viper.GetString("receiver"), network)
	if err != nil {
		log.Logger().WithField("error", err).Error("decode fee address failed")
		return
	}

	cfg := &logic.CollectCfg{
		Testnet:                 testnet,
		StrHotWalletFee1Privkey: viper.GetString("feeprivatekey"),
		HotWalletFee1:           feeAddress,
		HotWallet1:              receiverAddress,
	}
	if err := logic.Run(cfg); err != nil {
		log.Logger().WithField("error", err).Error("collect failed")
		alarm.Slack(context.Background(), fmt.Sprintf("collect failed: %s", err.Error()))
		return
	}
}
