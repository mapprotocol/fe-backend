package task

import (
	"context"
	"encoding/hex"
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/mapprotocol/fe-backend/dao"
	"github.com/mapprotocol/fe-backend/params"
	"github.com/mapprotocol/fe-backend/third-party/butter"
	"github.com/mapprotocol/fe-backend/utils"
	"github.com/mapprotocol/fe-backend/utils/alarm"
	"github.com/mapprotocol/fe-backend/utils/tx"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/viper"
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/liteclient"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/ton"

	"github.com/mapprotocol/fe-backend/resource/log"
)

const (
	EventIDTONToEVM = "" // todo
	EventIDEVMToTON = "" // todo
)

func HandlePendingOrdersOfFirstStageFromTONToEVM() error {
	tonCfg := viper.GetStringMapString("ton")
	client := liteclient.NewConnectionPool()

	cfg, err := liteclient.GetConfigFromUrl(context.Background(), tonCfg["url"])
	if err != nil {
		log.Logger().WithField("error", err.Error()).Error("failed to get config")
		return err
	}

	// connect to mainnet lite servers
	if err = client.AddConnectionsFromConfig(context.Background(), cfg); err != nil {
		log.Logger().WithField("error", err.Error()).Error("failed to add connection")
		return err
	}

	// initialize ton api lite connection wrapper with full proof checks
	api := ton.NewAPIClient(client, ton.ProofCheckPolicySecure).WithRetry()
	api.SetTrustedBlockFromConfig(cfg)

	log.Logger().Info("fetching and checking proofs since config init block, it may take near a minute...")
	master, err := api.CurrentMasterchainInfo(context.Background()) // we fetch block just to trigger chain proof check
	if err != nil {
		log.Logger().WithField("error", err.Error()).Error("failed to get master chain info")
		return err
	}
	log.Logger().Info("master proof checks are completed successfully, now communication is 100% safe!, master: ", master)

	// address on which we are accepting payments
	treasuryAddress := address.MustParseAddr(tonCfg["chainpool"])                  // todo contract name
	lastProcessedLT, err := strconv.ParseUint(tonCfg["lasttxlogicaltime"], 10, 64) // todo replace it
	if err != nil {
		log.Logger().WithField("error", err.Error()).Error("failed to parse last tx logical time")
		return err
	}
	// channel with new transactions
	transactions := make(chan *tlb.Transaction)

	// it is a blocking call, so we start it asynchronously
	go api.SubscribeOnTransactions(context.Background(), treasuryAddress, lastProcessedLT, transactions)

	log.Logger().Info("waiting for transfers...")

	// listen for new transactions from channel
	for t := range transactions {
		messages, err := t.IO.Out.ToSlice()
		if err != nil {
			log.Logger().WithField("error", err.Error()).Error("failed to get out messages")
			continue
		}
		for _, msg := range messages {
			if msg.MsgType != tlb.MsgTypeExternalOut {
				continue
			}

			// parse cell
			out := msg.AsExternalOut()
			//log.Println("src: ", out.SrcAddr)
			//log.Println("dst: ", out.DstAddr)
			//log.Println("dst: ", out.CreatedLT)
			//log.Println("dst: ", out.CreatedAt)
			//
			//slice := out.Payload().BeginParse()
			//orderID := slice.MustLoadBigUInt(256)
			//sender := slice.MustLoadAddr()
			//srcToken := slice.MustLoadAddr()
			//srcAmount := slice.MustLoadBigUInt(256)
			//log.Println("orderID", orderID)
			//log.Println("sender: ", sender)
			//log.Println("srcToken: ", srcToken)
			//log.Println("srcAmount: ", srcAmount)

			if strings.Contains(strings.ToLower(out.DstAddr.String()), EventIDTONToEVM) {

			} else if strings.Contains(strings.ToLower(out.DstAddr.String()), EventIDEVMToTON) {

			}

		}

		// update last processed lt and save it in db
		lastProcessedLT = t.LT // todo store to db
	}

	return errors.New("transaction listening unexpectedly finished")
}

func HandleConfirmedOrdersOfFirstStageFromTONToEVM() {
	//orders, err := order.GetOldest10ByStatus(id, dao.OrderActionToEVM, dao.OrderStag1, dao.OrderStatusConfirmed)
	order := dao.Order{
		SrcToken: params.TONChainID,
		Action:   dao.OrderActionToEVM,
		Stage:    dao.OrderStag1,
		Status:   dao.OrderStatusConfirmed,
	}
	for {
		for id := uint64(1); ; {
			orders, err := order.GetOldest10ByID(id)
			if err != nil {
				fields := map[string]interface{}{
					"order": utils.JSON(order),
					"error": err,
				}
				log.Logger().WithFields(fields).Error("failed to get confirmed status order from ton to evm")
				alarm.Slack(context.Background(), "failed to get confirmed status order from ton to evm")
				time.Sleep(5 * time.Second)
				continue
			}

			length := len(orders)
			if length == 0 {
				log.Logger().Info("not found confirmed status order from ton to evm", "time", time.Now())
				time.Sleep(10 * time.Second)
				break
			}

			for i, o := range orders {
				if i == length-1 {
					id = o.ID + 1
				}

				amount, ok := new(big.Int).SetString(o.InAmount, 10)
				if !ok {
					fields := map[string]interface{}{
						"orderId": o.ID,
						"amount":  o.InAmount,
						"error":   err,
					}
					log.Logger().WithFields(fields).Error("failed to parse string to big int")
					alarm.Slack(context.Background(), "failed to parse string to big int")
					continue
				}
				amount = new(big.Int).Quo(amount, big.NewInt(1e18)) // todo
				request := &butter.RouterAndSwapRequest{
					//FromChainID:     o.SrcChain,
					FromChainID: params.ChainIDOfChainPool, // todo
					ToChainID:   o.DstChain,
					//Amount:          o.InAmount,  //
					Amount:          amount.String(), // todo
					TokenInAddress:  params.USDTOfChainPool,
					TokenOutAddress: o.DstToken,
					Kind:            SwapType,
					Slippage:        strconv.FormatUint(o.Slippage/3, 10), // todo
					Entrance:        Entrance,
					//From:            o.Sender,
					From:     Sender, // todo
					Receiver: o.Receiver,
				}
				data, err := butter.RouterAndSwap(request)
				if err != nil {
					log.Logger().WithField("error", err.Error()).Error("failed to create router and swap request")
					alarm.Slack(context.Background(), "failed to create router and swap request")
					continue
				}

				chainInfo, err := dao.NewSupportedChainWithChainID(params.ChainIDOfChainPool).First()
				if err != nil {
					log.Logger().WithField("chainID", params.ChainIDOfChainPool).WithField("error", err.Error()).Error("failed to get chain info")
					alarm.Slack(context.Background(), "failed to get chain info")
					time.Sleep(5 * time.Second)
					continue
				}

				multiplier, err := strconv.ParseFloat(chainInfo.GasLimitMultiplier, 64)
				if err != nil {
					fields := map[string]interface{}{
						"chainID":            params.ChainIDOfChainPool,
						"gasLimitMultiplier": chainInfo.GasLimitMultiplier,
						"error":              err,
					}
					log.Logger().WithFields(fields).Error("failed to parse string to float")
					alarm.Slack(context.Background(), "failed to parse string to float")
					continue
				}
				transactor, err := tx.NewTransactor(chainInfo.ChainRPC, chainInfo.ChainPoolContract, multiplier)
				if err != nil {
					fields := map[string]interface{}{
						"chainID":            params.ChainIDOfChainPool,
						"rpc":                chainInfo.ChainRPC,
						"chainPoolContract":  chainInfo.ChainPoolContract,
						"gasLimitMultiplier": chainInfo.GasLimitMultiplier,
						"error":              err,
					}
					log.Logger().WithFields(fields).Error("failed to create transactor")
					alarm.Slack(context.Background(), "failed to create transactor")
					time.Sleep(5 * time.Second)
					continue
				}

				decodeData, err := DecodeData(data.Data)
				if err != nil {
					log.Logger().WithField("data", data.Data).WithField("error", err.Error()).Error("failed to decode call data")
					alarm.Slack(context.Background(), "failed to decode call data")
					continue
				}

				orderID := utils.Uint64ToByte32(o.ID)
				v := ""
				if data.Value[:2] == "0x" || data.Value[:2] == "0X" {
					v = data.Value[2:]
				}
				value, ok := new(big.Int).SetString(v, 16)
				if !ok {
					log.Logger().WithField("value", v).Error("failed to parse string to big int")
					alarm.Slack(context.Background(), "failed to parse string to big int")
					continue
				}
				hash, err := transactor.DeliverAndSwap(orderID, Initiator, common.HexToAddress(params.USDTOfChainPool), amount, decodeData.SwapData, decodeData.BridgeData, decodeData.FeeData, value)
				if err != nil {
					fields := map[string]interface{}{
						"orderID":    hex.EncodeToString(orderID[:]),
						"initiator":  Initiator,
						"token":      params.WBTCOfChainPool,
						"amount":     value,
						"SwapData":   hex.EncodeToString(decodeData.SwapData),
						"BridgeData": hex.EncodeToString(decodeData.BridgeData),
						"FeeData":    hex.EncodeToString(decodeData.FeeData),
						"error":      err.Error(),
					}
					log.Logger().WithFields(fields).Error("failed to send deliver and swap transaction")
					alarm.Slack(context.Background(), "failed to send deliver and swap transaction")
					continue
				}

				update := &dao.Order{
					Stage:     dao.OrderStag2,
					Status:    dao.OrderStatusPending,
					OutTxHash: hash.String(),
				}
				if err := dao.NewOrderWithID(o.ID).Updates(update); err != nil {
					log.Logger().WithField("update", utils.JSON(update)).WithField("error", err.Error()).Error("failed to update order status")
					alarm.Slack(context.Background(), "failed to update order status")
					time.Sleep(5 * time.Second)
					continue
				}

			}
			time.Sleep(10 * time.Second)
		}
	}
}

func HandlePendingOrdersOfSecondStageFromTONToEVM() {
	order := dao.Order{
		SrcToken: params.TONChainID,
		Action:   dao.OrderActionToEVM,
		Stage:    dao.OrderStag2,
		Status:   dao.OrderStatusPending,
	}
	for {
		for id := uint64(1); ; {
			orders, err := order.GetOldest10ByID(id)
			if err != nil {
				log.Logger().WithField("error", err.Error()).Error("failed to get confirmed status order from ton to evm")
				alarm.Slack(context.Background(), "failed to get confirmed status order from ton to evm")
				time.Sleep(5 * time.Second)
				continue
			}

			length := len(orders)
			if length == 0 {
				log.Logger().Info("not found confirmed status order from ton to evm", "time", time.Now())
				time.Sleep(10 * time.Second)
				break
			}

			for i, o := range orders {
				if i == length-1 {
					id = o.ID + 1
				}

				chainInfo, err := dao.NewSupportedChainWithChainID(params.ChainIDOfChainPool).First()
				if err != nil {
					log.Logger().WithField("chainID", o.DstChain).WithField("error", err.Error()).Error("failed to get chain info")
					alarm.Slack(context.Background(), "failed to get chain info")
					time.Sleep(5 * time.Second)
					continue
				}

				// todo call NewCaller
				transactor, err := tx.NewTransactor(chainInfo.ChainRPC, chainInfo.ChainPoolContract, 0)
				if err != nil {
					fields := map[string]interface{}{
						"chainID":            params.ChainIDOfChainPool,
						"rpc":                chainInfo.ChainRPC,
						"chainPoolContract":  chainInfo.ChainPoolContract,
						"gasLimitMultiplier": 0,
						"error":              err,
					}
					log.Logger().WithFields(fields).Error("failed to create transactor")
					alarm.Slack(context.Background(), "failed to create transactor")
					time.Sleep(5 * time.Second)
					continue
				}
				pending, err := transactor.TransactionIsPending(common.HexToHash(o.OutTxHash))
				if err != nil {
					fields := map[string]interface{}{
						"chainID": params.ChainIDOfChainPool,
						"rpc":     chainInfo.ChainRPC,
						"txHash":  o.OutTxHash,
						"error":   err,
					}
					log.Logger().WithFields(fields).Error("failed to judge transaction is pending")
					alarm.Slack(context.Background(), "failed to judge transaction is p ending")
					time.Sleep(5 * time.Second)
					continue
				}
				if pending {
					continue
				}

				status, err := transactor.TransactionStatus(common.HexToHash(o.OutTxHash))
				if err != nil {
					log.Logger().WithField("endpoint", chainInfo.ChainRPC).WithField("txHash", o.OutTxHash).WithField("error", err.Error()).Error("get transaction status")
					alarm.Slack(context.Background(), "get transaction status")
					time.Sleep(5 * time.Second)
					continue
				}

				update := &dao.Order{
					Status: dao.OrderStatusConfirmed,
				}
				if status == types.ReceiptStatusFailed {
					update.Status = dao.OrderStatusFailed
				}
				if err := dao.NewOrderWithID(o.ID).Updates(update); err != nil {
					log.Logger().WithField("update", utils.JSON(update)).WithField("error", err.Error()).Error("failed to update order status")
					alarm.Slack(context.Background(), "failed to update order status")
					time.Sleep(5 * time.Second)
					continue
				}
			}
			time.Sleep(10 * time.Second)
		}
	}
}
