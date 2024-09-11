package main

import (
	"fmt"
	blog "log"
	"os"
	"os/signal"
	"runtime/debug"
	"strings"
	"syscall"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/viper"
	"golang.org/x/term"

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

	tips := "Enter the chain pool manager password on the EVM chain: "
	privateKey, err := GetWalletKey(tips, viper.GetStringMapString("chainPool")["keystore"])
	if err != nil {
		panic(err)
	}
	passwd, err := getPassphrase("Enter the chain pool manager password on the TON chain: ")
	if err != nil {
		panic(err)
	}

	// init log
	log.Init(viper.GetString("env"), viper.GetString("logDir"))
	// init db
	dbConfig := viper.GetStringMapString("database")
	db.Init(dbConfig["user"], dbConfig["password"], dbConfig["host"], dbConfig["port"], dbConfig["name"])

	task.InitMempoolClient(viper.GetString("network"))
	tx.InitTransactor(privateKey)
	tonclient.Init(viper.GetStringMapString("ton")["words"], passwd)

	filter.Init()
	butter.Init()
	tonrouter.Init()
	task.Init()

	runTask()
	runBTCTask()
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
		defer func() {
			stack := string(debug.Stack())
			log.Logger().WithField("stack", stack).Error("failed to HandlePendingOrdersOfFirstStageFromEVM")

			if r := recover(); r != nil {
				log.Logger().WithField("error", r).Error("failed to recover HandlePendingOrdersOfFirstStageFromEVM")
			}
		}()

		task.HandlePendingOrdersOfFirstStageFromEVM()
	}()
}

func runBTCTask() {
	go func() {
		defer func() {
			stack := string(debug.Stack())
			log.Logger().WithField("stack", stack).Error("failed to GetFeeRate")

			if r := recover(); r != nil {
				log.Logger().WithField("error", r).Error("failed to recover GetFeeRate")
			}
		}()

		task.GetFeeRate()
	}()

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
		defer func() {
			stack := string(debug.Stack())
			log.Logger().WithField("stack", stack).Error("failed to HandlePendingOrdersOfFirstStageFromTONToEVM")

			if r := recover(); r != nil {
				log.Logger().WithField("error", r).Error("failed to recover HandlePendingOrdersOfFirstStageFromTONToEVM")
			}
		}()

		task.HandlePendingOrdersOfFirstStageFromTONToEVM()
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

	go func() {
		defer func() {
			stack := string(debug.Stack())
			log.Logger().WithField("stack", stack).Error("failed to HandlePendingOrdersOfSecondSStageFromEVMToTON")

			if r := recover(); r != nil {
				log.Logger().WithField("error", r).Error("failed to recover HandlePendingOrdersOfSecondSStageFromEVMToTON")
			}
		}()

		task.HandlePendingOrdersOfSecondSStageFromEVMToTON()
	}()
}

func getPassphrase(tips string) (string, error) {
	fmt.Print(tips)
	password, err := term.ReadPassword(syscall.Stdin)
	if err != nil {
		return "", err
	}
	fmt.Printf("\n")
	return strings.TrimSpace(string(password)), nil
}

func GetWalletKey(tips, filepath string) (string, error) {
	// Read key from file.
	keyjson, err := os.ReadFile(filepath)
	if err != nil {
		return "", err
	}
	// Decrypt key with passphrase.
	passphrase, err := getPassphrase(tips)
	if err != nil {
		return "", err
	}
	key, err := keystore.DecryptKey(keyjson, passphrase)
	if err != nil {
		return "", err
	}
	return hexutil.Encode(crypto.FromECDSA(key.PrivateKey)), nil
}
