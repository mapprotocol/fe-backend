package main

import (
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"github.com/mapprotocol/fe-backend/config"
	"github.com/mapprotocol/fe-backend/resource/ceffu"
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

	// init ceffu client
	ceffuConf := viper.GetStringMapString("ceffu")
	ceffu.Init(ceffuConf["domain"], ceffuConf["key"], ceffuConf["secret"])

	engine := gin.Default()
	router.Register(engine)
	_ = endless.ListenAndServe(viper.GetString("address"), engine)
}
