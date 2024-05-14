package main

import (
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"github.com/mapprotocol/ceffu-fe-backend/config"
	"github.com/mapprotocol/ceffu-fe-backend/resource/db"
	"github.com/mapprotocol/ceffu-fe-backend/resource/log"
	"github.com/mapprotocol/ceffu-fe-backend/router"
)

func main() {
	// init config
	config.Init()
	// init log
	log.Init(viper.GetString("env"), viper.GetString("logDir"))
	// init db
	dbConf := viper.GetStringMapString("database")
	db.Init(dbConf["user"], dbConf["password"], dbConf["host"], dbConf["port"], dbConf["name"])

	engine := gin.Default()
	router.Register(engine)
	_ = endless.ListenAndServe(viper.GetString("address"), engine)
}
