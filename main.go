package main

import (
	blog "log"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"

	"github.com/spf13/viper"

	"github.com/mapprotocol/fe-backend/config"
	"github.com/mapprotocol/fe-backend/resource/db"
	"github.com/mapprotocol/fe-backend/resource/log"
	"github.com/mapprotocol/fe-backend/task"
	"github.com/mapprotocol/fe-backend/utils/tx"
)

func main() {
	// init config
	config.Init()
	// init log
	log.Init(viper.GetString("env"), viper.GetString("logDir"))
	// init db
	dbConf := viper.GetStringMapString("database")
	db.Init(dbConf["user"], dbConf["password"], dbConf["host"], dbConf["port"], dbConf["name"])

	task.InitMempoolClient(viper.GetString("network"))

	tx.InitTransactor(viper.GetString("senderprivatekey"))

	runTask()
	//runBTCTask()
	runTONTask()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer signal.Stop(sigs)
	select {
	case <-sigs:
		blog.Println("Interrupt received, shutting down now.")
	}
}

func runTask() {
	go func() {
		var err error
		defer func() {
			stack := string(debug.Stack())
			log.Logger().WithField("error", err).WithField("stack", stack).Error("failed to HandlePendingOrdersOfFirstStageFromEVM")

			if r := recover(); r != nil {
				log.Logger().WithField("error", r).Error("failed to recover HandlePendingOrdersOfFirstStageFromEVM")
			}
		}()

		err = task.HandlePendingOrdersOfFirstStageFromEVM()
	}()
}

func runBTCTask() {
	go func() {
		defer func() {
			stack := string(debug.Stack())
			log.Logger().WithField("stack", stack).Error("failed to HandlePendingOrdersOfFirstStageFromBTCToEVM")

			if r := recover(); r != nil {
				log.Logger().WithField("error", r).Error("failed to recover HandlePendingOrdersOfFirstStageFromBTCToEVM")
			}
		}()

		task.HandlePendingOrdersOfFirstStageFromBTCToEVM()
	}()

	go func() {
		defer func() {
			stack := string(debug.Stack())
			log.Logger().WithField("stack", stack).Error("failed to HandleConfirmedOrdersOfFirstStageFromBTCToEVM")

			if r := recover(); r != nil {
				log.Logger().WithField("error", r).Error("failed to recover HandleConfirmedOrdersOfFirstStageFromBTCToEVM")
			}
		}()

		task.HandleConfirmedOrdersOfFirstStageFromBTCToEVM()
	}()

	go func() {
		defer func() {
			stack := string(debug.Stack())
			log.Logger().WithField("stack", stack).Error("failed to HandlePendingOrdersOfSecondStageFromBTCToEVM")

			if r := recover(); r != nil {
				log.Logger().WithField("error", r).Error("failed to recover HandlePendingOrdersOfSecondStageFromBTCToEVM")
			}
		}()

		task.HandlePendingOrdersOfSecondStageFromBTCToEVM()
	}()
}

func runTONTask() {
	go func() {
		var err error
		defer func() {
			stack := string(debug.Stack())
			log.Logger().WithField("error", err).WithField("stack", stack).Error("failed to HandlePendingOrdersOfFirstStageFromTONToEVM")

			if r := recover(); r != nil {
				log.Logger().WithField("error", r).Error("failed to recover HandlePendingOrdersOfFirstStageFromTONToEVM")
			}
		}()

		err = task.HandlePendingOrdersOfFirstStageFromTONToEVM()
	}()

	go func() {
		defer func() {
			stack := string(debug.Stack())
			log.Logger().WithField("stack", stack).Error("failed to HandleConfirmedOrdersOfFirstStageFromTONToEVM")

			if r := recover(); r != nil {
				log.Logger().WithField("error", r).Error("failed to recover HandleConfirmedOrdersOfFirstStageFromTONToEVM")
			}
		}()

		task.HandleConfirmedOrdersOfFirstStageFromTONToEVM()
	}()

	go func() {
		defer func() {
			stack := string(debug.Stack())
			log.Logger().WithField("stack", stack).Error("failed to HandlePendingOrdersOfSecondStageFromTONToEVM")

			if r := recover(); r != nil {
				log.Logger().WithField("error", r).Error("failed to recover HandlePendingOrdersOfSecondStageFromTONToEVM")
			}
		}()

		task.HandlePendingOrdersOfSecondStageFromTONToEVM()
	}()
}
