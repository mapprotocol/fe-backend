package main

import (
	"github.com/mapprotocol/fe-backend/logic"
	"github.com/mapprotocol/fe-backend/utils/alarm"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"syscall"

	"github.com/mapprotocol/fe-backend/config"
	"github.com/mapprotocol/fe-backend/resource/db"
	"github.com/mapprotocol/fe-backend/resource/log"
)

func main() {
	alarm.ValidateEnv()

	// init config
	config.Init()
	// init log
	log.Init(viper.GetString("env"), viper.GetString("logDir"))
	// init db
	dbConf := viper.GetStringMapString("database")
	db.Init(dbConf["user"], dbConf["password"], dbConf["host"], dbConf["port"], dbConf["name"])

	cfg := &logic.CollectCfg{
		Testnet:       viper.GetBool("testnet"),
		StrFeePrivkey: viper.GetString("feeprivatekey"),
		//FeeAddress:    viper.GetString("feeaddress"), // todo
		//Receiver:      viper.GetString("receiver"), // todo
	}
	logic.RunCollect(cfg)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer signal.Stop(sigs)
	select {
	case <-sigs:
		log.Logger().Info("Signal received, shutting down now.")
	}
}
