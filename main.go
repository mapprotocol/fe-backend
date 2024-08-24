package main

import (
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"github.com/mapprotocol/fe-backend/logic"
	"github.com/mapprotocol/fe-backend/third-party/butter"
	"github.com/mapprotocol/fe-backend/third-party/tonrouter"
	"github.com/spf13/viper"

	"github.com/mapprotocol/fe-backend/config"
	"github.com/mapprotocol/fe-backend/resource/db"
	"github.com/mapprotocol/fe-backend/resource/log"
	"github.com/mapprotocol/fe-backend/router"
)

func main() {
	// init config
	config.Init()
	// init log
	log.Init(viper.GetString("env"), viper.GetString("logDir"))
	// init db
	dbConf := viper.GetStringMapString("database")
	db.Init(dbConf["user"], dbConf["password"], dbConf["host"], dbConf["port"], dbConf["name"])

	butter.Init()
	tonrouter.Init()
	logic.Init()

	bitcoinConf := viper.GetStringMapString("bitcoin")
	logic.InitMempoolClient(bitcoinConf["network"], bitcoinConf["vault"])

	engine := gin.Default()
	router.Register(engine)
	_ = endless.ListenAndServe(viper.GetString("address"), engine)
}
