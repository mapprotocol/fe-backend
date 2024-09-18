package main

import (
	blog "log"
	"os"
	"os/signal"
	"syscall"

	_func "github.com/mapprotocol/fe-backend/utils/func"

	"github.com/spf13/viper"

	"github.com/mapprotocol/fe-backend/config"
	"github.com/mapprotocol/fe-backend/resource/db"
	"github.com/mapprotocol/fe-backend/resource/log"
	"github.com/mapprotocol/fe-backend/resource/tonclient"
	"github.com/mapprotocol/fe-backend/task"
	"github.com/mapprotocol/fe-backend/third-party/butter"
	"github.com/mapprotocol/fe-backend/third-party/filter"
	"github.com/mapprotocol/fe-backend/third-party/tonrouter"
	"github.com/mapprotocol/fe-backend/utils/tx"
)

func main() {
	// init config
	config.Init()
	// init log
	log.Init(viper.GetString("env"), viper.GetString("logDir"))
	// init db
	dbConfig := viper.GetStringMapString("database")
	db.Init(dbConfig["user"], dbConfig["password"], dbConfig["host"], dbConfig["port"], dbConfig["name"])

	task.InitMempoolClient(viper.GetString("network"))

	tx.InitTransactor(viper.GetStringMapString("chainPool")["senderprivatekey"])

	tonConfig := viper.GetStringMapString("ton")
	tonclient.Init(tonConfig["words"], tonConfig["password"])

	filter.Init()
	butter.Init()
	tonrouter.Init()
	task.Init()

	//runTask()
	//runBTCTask()
	//runTONTask()
	runSolTask()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer signal.Stop(sigs)
	select {
	case <-sigs:
		blog.Println("Interrupt received, shutting down now.")
	}
}

func runTask() {
	_func.Go(task.HandlePendingOrdersOfFirstStageFromEVM, "HandlePendingOrdersOfFirstStageFromEVM")
}

func runBTCTask() {
	_func.Go(task.HandlePendingOrdersOfFirstStageFromBTCToEVM, "HandlePendingOrdersOfFirstStageFromBTCToEVM")
	_func.Go(task.HandleConfirmedOrdersOfFirstStageFromBTCToEVM, "HandleConfirmedOrdersOfFirstStageFromBTCToEVM")
	_func.Go(task.HandlePendingOrdersOfSecondStageFromBTCToEVM, "HandlePendingOrdersOfSecondStageFromBTCToEVM")
}

func runTONTask() {
	_func.Go(task.HandlePendingOrdersOfFirstStageFromTONToEVM, "HandlePendingOrdersOfFirstStageFromTONToEVM")
	_func.Go(task.HandleConfirmedOrdersOfFirstStageFromTONToEVM, "HandleConfirmedOrdersOfFirstStageFromTONToEVM")
	_func.Go(task.HandlePendingOrdersOfSecondStageFromTONToEVM, "HandlePendingOrdersOfSecondStageFromTONToEVM")
	_func.Go(task.HandleConfirmedOrdersOfFirstStageFromEVMToTON, "HandleConfirmedOrdersOfFirstStageFromEVMToTON")
	_func.Go(task.HandlePendingOrdersOfSecondSStageFromEVMToTON, "HandlePendingOrdersOfSecondSStageFromEVMToTON")
}

func runSolTask() {
	//_func.Go(task.FilterEventToSol, "FilterEventToSol")
	_func.Go(task.HandlerEvm2Sol, "HandlerEvm2Sol")
	_func.Go(task.FilterSol2Evm, "FilterSol2Evm")
	_func.Go(task.HandleSol2EvmButter, "HandleSol2EvmButter")
	_func.Go(task.HandleSol2EvmFinish, "HandleSol2EvmFinish")
}
