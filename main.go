package main

import (
	"fmt"
	"github.com/mapprotocol/fe-backend/resource/tonclient"
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

const (
	VersionMajor = 0        // Major version component of the current release
	VersionMinor = 0        // Minor version component of the current release
	VersionPatch = 1        // Patch version component of the current release
	VersionMeta  = "stable" // Version metadata to append to the version string
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
	// init config
	config.Init()
	// init log
	log.Init(viper.GetString("env"), viper.GetString("logDir"))
	// init db
	dbConfig := viper.GetStringMapString("database")
	db.Init(dbConfig["user"], dbConfig["password"], dbConfig["host"], dbConfig["port"], dbConfig["name"])

	//task.InitMempoolClient(viper.GetString("network"))

	tx.InitTransactor(viper.GetStringMapString("chainpool")["senderprivatekey"])

	tonConfig := viper.GetStringMapString("ton")
	tonclient.Init(tonConfig["words"], tonConfig["password"])

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

	go func() {
		defer func() {
			stack := string(debug.Stack())
			log.Logger().WithField("stack", stack).Error("failed to HandleConfirmedOrdersOfFirstStageFromEVMToTON")

			if r := recover(); r != nil {
				log.Logger().WithField("error", r).Error("failed to recover HandleConfirmedOrdersOfFirstStageFromEVMToTON")
			}
		}()

		task.HandleConfirmedOrdersOfFirstStageFromEVMToTON()
	}()

}
