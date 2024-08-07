package task

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	blog "log"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/mapprotocol/fe-backend/dao"
	"github.com/mapprotocol/fe-backend/params"
	"github.com/mapprotocol/fe-backend/resource/tonclient"
	"github.com/mapprotocol/fe-backend/third-party/butter"
	"github.com/mapprotocol/fe-backend/third-party/filter"
	"github.com/mapprotocol/fe-backend/third-party/tonrouter"
	"github.com/mapprotocol/fe-backend/utils"
	"github.com/mapprotocol/fe-backend/utils/alarm"
	"github.com/mapprotocol/fe-backend/utils/tx"
	"github.com/spf13/viper"
	"github.com/xssnick/tonutils-go/ton/wallet"
	"github.com/xssnick/tonutils-go/tvm/cell"

	"github.com/mapprotocol/fe-backend/resource/log"
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tlb"
)

var FeeRate = big.NewFloat(70) // 70/10000

const (
	EventIDTONToEVM = "34a7e0e8"
	EventIDEVMToTON = "1a6c0a51"
)

var isMultiChainPool = false

func Init() {
	isMultiChainPool = viper.GetBool("isMultiChainPool")
	blog.Println("isMultiChainPool: ", isMultiChainPool)
}

// HandlePendingOrdersOfFirstStageFromTONToEVM
// create order from ton to evm( action=1, stage=1, status=2)
func HandlePendingOrdersOfFirstStageFromTONToEVM() {
	filterLog := dao.NewFilterLog(params.TONChainID, EventIDTONToEVM)
	for {
		gotLog, err := filterLog.First()
		if err != nil {
			fields := map[string]interface{}{
				"chainID": params.TONChainID,
				"topic":   EventIDTONToEVM,
				"error":   err.Error(),
			}
			log.Logger().WithFields(fields).Error("failed to get filter log info")
			alarm.Slack(context.Background(), "failed to get filter log info")
			time.Sleep(5 * time.Second)
			continue
		}

		logs, err := filter.GetLogs(gotLog.LatestLogID, params.TONChainID, EventIDTONToEVM, uint8(20))
		if err != nil {
			fields := map[string]interface{}{
				"id":      gotLog.LatestLogID,
				"chainID": params.TONChainID,
				"topic":   EventIDTONToEVM,
				"limit":   uint8(20),
				"error":   err.Error(),
			}
			log.Logger().WithFields(fields).Error("failed to get logs")
			alarm.Slack(context.Background(), "failed to get logs")
			continue
		}

		for _, lg := range logs {
			if lg.Id <= gotLog.LatestLogID {
				continue
			}

			logData, err := hex.DecodeString(lg.LogData)
			if err != nil {
				fields := map[string]interface{}{
					"id":      lg.Id,
					"chainID": params.TONChainID,
					"topic":   EventIDTONToEVM,
					"logData": lg.LogData,
					"error":   err.Error(),
				}
				log.Logger().WithFields(fields).Error("failed to decode log data")
				alarm.Slack(context.Background(), "failed to decode log data")
				continue
			}

			body := &cell.Cell{}
			if err := json.Unmarshal(logData, &body); err != nil {
				log.Logger().WithField("logData", lg.LogData).WithField("error", err.Error()).Error("failed to unmarshal data")
				alarm.Slack(context.Background(), "failed to unmarshal data")
				continue
			}
			slice := body.BeginParse()
			orderID, err := slice.LoadUInt(64)
			if err != nil {
				// todo add unique flag to log
				log.Logger().WithField("error", err.Error()).Error("failed to load order id")
				alarm.Slack(context.Background(), "failed to load order id from")
				continue
			}
			from, err := slice.LoadRef()
			if err != nil {
				log.Logger().WithField("error", err.Error()).Error("failed to load from ref")
				alarm.Slack(context.Background(), "failed to load from ref")
				continue
			}
			to, err := slice.LoadRef()
			if err != nil {
				log.Logger().WithField("error", err.Error()).Error("failed to load to ref")
				alarm.Slack(context.Background(), "failed to load to ref")
				continue
			}
			srcChain, err := from.LoadUInt(64)
			if err != nil {
				log.Logger().WithField("error", err.Error()).Error("failed to load from chain id")
				alarm.Slack(context.Background(), "failed to load from chain id")
				continue
			}
			sender, err := from.LoadAddr()
			if err != nil {
				log.Logger().WithField("error", err.Error()).Error("failed to load sender")
				alarm.Slack(context.Background(), "failed to load sender")
				continue
			}
			srcToken, err := from.LoadAddr()
			if err != nil {
				log.Logger().WithField("error", err.Error()).Error("failed to load src token")
				alarm.Slack(context.Background(), "failed to load rc token")
				continue
			}
			inAmount, err := from.LoadUInt(64)
			if err != nil {
				log.Logger().WithField("error", err.Error()).Error("failed to load amount in")
				alarm.Slack(context.Background(), "failed to load amount in")
				continue
			}
			slippage, err := from.LoadUInt(16)
			if err != nil {
				log.Logger().WithField("error", err.Error()).Error("failed to load slippage")
				alarm.Slack(context.Background(), "failed to load slippage")
				continue
			}
			dstChain, err := to.LoadUInt(64)
			if err != nil {
				log.Logger().WithField("error", err.Error()).Error("failed to load to chain id")
				alarm.Slack(context.Background(), "failed to load to chain id")
				continue
			}
			receiver, err := to.LoadBigUInt(160)
			if err != nil {
				log.Logger().WithField("error", err.Error()).Error("failed to load receiver")
				alarm.Slack(context.Background(), "failed to load receiver")
				continue
			}
			dstToken, err := to.LoadBigUInt(160)
			if err != nil {
				log.Logger().WithField("error", err.Error()).Error("failed to load token out address")
				alarm.Slack(context.Background(), "failed to load token out address")
				continue
			}
			relayAmount, err := slice.LoadUInt(32)
			if err != nil {
				log.Logger().WithField("error", err.Error()).Error("failed to load jetton amount")
				alarm.Slack(context.Background(), "failed to load jetton amount")
				continue
			}

			srcTokenStr := srcToken.String()
			if srcTokenStr == params.NoneAddress {
				srcTokenStr = params.NativeOfTON
			}
			_, afterAmount := deductFees(new(big.Float).SetUint64(relayAmount), FeeRate)
			// convert token to float like 0.089
			afterAmountFloat := new(big.Float).Quo(afterAmount, big.NewFloat(params.USDTDecimalOfTON))
			inAmountFloat := new(big.Float).Quo(new(big.Float).SetUint64(inAmount), big.NewFloat(params.InAmountDecimalOfTON))
			order := &dao.Order{
				OrderIDFromContract: orderID,
				SrcChain:            strconv.FormatUint(srcChain, 10),
				SrcToken:            srcTokenStr,
				Sender:              sender.String(),
				InAmount:            inAmountFloat.Text('f', -1),
				RelayToken:          params.USDTOfTON,
				RelayAmount:         afterAmountFloat.String(),
				DstChain:            strconv.FormatUint(dstChain, 10),
				DstToken:            common.BytesToAddress(dstToken.Bytes()).String(),
				Receiver:            common.BytesToAddress(receiver.Bytes()).String(),
				Action:              dao.OrderActionToEVM,
				Stage:               dao.OrderStag1,
				Status:              dao.OrderStatusConfirmed,
				Slippage:            slippage,
			}

			if err := order.Create(); err != nil {
				log.Logger().WithField("order", utils.JSON(order)).WithField("error", err).Error("failed to create order")
				alarm.Slack(context.Background(), "failed to update order status")
				continue
			}

			if err := filterLog.UpdateLatestLogID(lg.Id); err != nil {
				fields := map[string]interface{}{
					"chainID":     params.ChainIDOfChainPool,
					"topic":       params.OnReceivedTopic,
					"latestLogID": lg.Id,
					"error":       err.Error(),
				}
				log.Logger().WithFields(fields).Error("failed to update filter log")
				alarm.Slack(context.Background(), "failed to update filter log")
				continue
			}
		}

		time.Sleep(10 * time.Second)
	}
}

// HandleConfirmedOrdersOfFirstStageFromTONToEVM
// action=1, stage=1, status=2 ==> stage=2, status=1
func HandleConfirmedOrdersOfFirstStageFromTONToEVM() {
	order := dao.Order{
		SrcChain: params.TONChainID,
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
				log.Logger().WithField("time", time.Now()).Info("not found confirmed status order from ton to evm")
				time.Sleep(10 * time.Second)
				break
			}

			for i, o := range orders {
				if i == length-1 {
					id = o.ID + 1
				}

				//amountFloat, ok := new(big.Float).SetString(o.RelayAmount)
				//if !ok {
				//	fields := map[string]interface{}{
				//		"orderId": o.ID,
				//		"amount":  o.InAmount,
				//		"error":   err,
				//	}
				//	log.Logger().WithFields(fields).Error("failed to parse string to big float")
				//	alarm.Slack(context.Background(), "failed to parse string to big float")
				//	continue
				//}

				usdt := params.USDTOfChainPool
				decimal := params.USDTDecimalOfChainPool
				chainID := params.ChainIDOfChainPool
				chainInfo := &dao.ChainPool{}
				if isMultiChainPool && o.SrcChain == params.ChainIDOfEthereum {
					usdt = params.USDTOfEthereum
					decimal = params.USDTDecimalOfEthereum
					chainID = params.ChainIDOfEthereum
					chainInfo, err = dao.NewChainPoolWithChainID(params.ChainIDOfEthereum).First()
					if err != nil {
						log.Logger().WithField("chainID", params.ChainIDOfChainPool).WithField("error", err.Error()).Error("failed to get chain info")
						alarm.Slack(context.Background(), "failed to get chain info")
						time.Sleep(5 * time.Second)
						continue

					}
				} else {
					chainInfo, err = dao.NewChainPoolWithChainID(params.ChainIDOfChainPool).First()
					if err != nil {
						log.Logger().WithField("chainID", params.ChainIDOfChainPool).WithField("error", err.Error()).Error("failed to get chain info")
						alarm.Slack(context.Background(), "failed to get chain info")
						time.Sleep(5 * time.Second)
						continue
					}
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
				transactor, err := tx.NewTransactor(chainInfo.ChainRPC, chainInfo.FeRouterContract, multiplier)
				if err != nil {
					fields := map[string]interface{}{
						"chainID":            params.ChainIDOfChainPool,
						"rpc":                chainInfo.ChainRPC,
						"feRouterContract":   chainInfo.FeRouterContract,
						"gasLimitMultiplier": chainInfo.GasLimitMultiplier,
						"error":              err,
					}
					log.Logger().WithFields(fields).Error("failed to create transactor")
					alarm.Slack(context.Background(), "failed to create transactor")
					time.Sleep(5 * time.Second)
					continue
				}

				orderID := utils.Uint64ToByte32(o.OrderIDFromContract)
				amount, ok := new(big.Float).SetString(o.RelayAmount)
				if !ok {
					fields := map[string]interface{}{
						"orderId": o.ID,
						"amount":  o.RelayAmount,
						"error":   err,
					}
					log.Logger().WithFields(fields).Error("failed to parse string to big int")
					alarm.Slack(context.Background(), "failed to parse string to big int")
					continue
				}
				amount = new(big.Float).Mul(amount, big.NewFloat(decimal))
				amountInt, _ := amount.Int(nil)

				txHash := common.Hash{}
				if o.DstChain == chainID && strings.ToLower(o.DstToken) == strings.ToLower(usdt) {
					txHash, err = deliver(transactor, common.HexToAddress(usdt), orderID, amountInt, common.HexToAddress(o.Receiver))
					if err != nil {
						time.Sleep(5 * time.Second)
						continue
					}
				} else {
					request := &butter.RouterAndSwapRequest{
						//FromChainID:     o.SrcChain,
						FromChainID: params.ChainIDOfChainPool, // todo
						ToChainID:   o.DstChain,
						//Amount:          o.InAmount,  //
						Amount:          o.RelayAmount, // todo
						TokenInAddress:  params.USDTOfChainPool,
						TokenOutAddress: o.DstToken,
						Type:            SwapType,
						Slippage:        o.Slippage / 3 * 2, // todo
						//From:            o.Sender,
						From:     Sender, // todo
						Receiver: o.Receiver,
					}
					txHash, err = deliverAndSwap(request, transactor, common.HexToAddress(usdt), orderID, amountInt)
					if err != nil {
						time.Sleep(5 * time.Second)
						continue
					}
				}

				update := &dao.Order{
					Stage:     dao.OrderStag2,
					Status:    dao.OrderStatusPending,
					OutTxHash: txHash.String(),
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

// HandlePendingOrdersOfSecondStageFromTONToEVM
// action=1, stage=2, status=1 ==> status=2/3
func HandlePendingOrdersOfSecondStageFromTONToEVM() {
	order := dao.Order{
		SrcChain: params.TONChainID,
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

				chainInfo := &dao.ChainPool{}
				if isMultiChainPool && o.SrcChain == params.ChainIDOfEthereum {
					chainInfo, err = dao.NewChainPoolWithChainID(params.ChainIDOfEthereum).First()
					if err != nil {
						log.Logger().WithField("chainID", params.ChainIDOfChainPool).WithField("error", err.Error()).Error("failed to get chain info")
						alarm.Slack(context.Background(), "failed to get chain info")
						time.Sleep(5 * time.Second)
						continue

					}
				} else {
					chainInfo, err = dao.NewChainPoolWithChainID(params.ChainIDOfChainPool).First()
					if err != nil {
						log.Logger().WithField("chainID", params.ChainIDOfChainPool).WithField("error", err.Error()).Error("failed to get chain info")
						alarm.Slack(context.Background(), "failed to get chain info")
						time.Sleep(5 * time.Second)
						continue
					}
				}

				transactor, err := tx.NewTransactor(chainInfo.ChainRPC, chainInfo.FeRouterContract, 0)
				if err != nil {
					fields := map[string]interface{}{
						"chainID":            params.ChainIDOfChainPool,
						"rpc":                chainInfo.ChainRPC,
						"feRouterContract":   chainInfo.FeRouterContract,
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

// HandleConfirmedOrdersOfSecondStageFromTONToEVM
// action=1, stage=2, status=2 ==> status=4
func HandleConfirmedOrdersOfSecondStageFromTONToEVM() {
	filterLog := dao.NewFilterLog(params.ChainIDOfChainPool, EventIDTONToEVM)
	for {
		gotLog, err := filterLog.First()
		if err != nil {
			fields := map[string]interface{}{
				"chainID": params.ChainIDOfChainPool,
				"topic":   EventIDTONToEVM,
				"error":   err.Error(),
			}
			log.Logger().WithFields(fields).Error("failed to get filter log info")
			alarm.Slack(context.Background(), "failed to get filter log info")
			time.Sleep(5 * time.Second)
			continue
		}

		logs, err := filter.GetLogs(gotLog.LatestLogID, params.ChainIDOfChainPool, EventIDTONToEVM, uint8(20))
		if err != nil {
			fields := map[string]interface{}{
				"id":      gotLog.LatestLogID,
				"chainID": params.ChainIDOfChainPool,
				"topic":   EventIDTONToEVM,
				"limit":   uint8(20),
				"error":   err.Error(),
			}
			log.Logger().WithFields(fields).Error("failed to get logs")
			alarm.Slack(context.Background(), "failed to get logs")
			continue
		}

		for _, lg := range logs {
			if lg.Id <= gotLog.LatestLogID {
				continue
			}

			logData, err := hex.DecodeString(lg.LogData)
			if err != nil {
				fields := map[string]interface{}{
					"id":      lg.Id,
					"chainID": params.ChainIDOfChainPool,
					"topic":   params.DeliverAndSwapTopic,
					"logData": lg.LogData,
					"error":   err.Error(),
				}
				log.Logger().WithFields(fields).Error("failed to decode log data")
				alarm.Slack(context.Background(), "failed to decode log data")
				continue
			}
			deliverAndSwap, err := UnpackDeliverAndSwap(logData)
			if err != nil {
				fields := map[string]interface{}{
					"id":      lg.Id,
					"chainID": params.ChainIDOfChainPool,
					"topic":   params.DeliverAndSwapTopic,
					"logData": lg.LogData,
					"error":   err.Error(),
				}
				log.Logger().WithFields(fields).Error("failed to unpack deliver and swap log data")
				alarm.Slack(context.Background(), "failed to unpack deliver and swap log data")
				continue
			}
			_ = deliverAndSwap

			if err := filterLog.UpdateLatestLogID(lg.Id); err != nil {
				fields := map[string]interface{}{
					"chainID":     params.ChainIDOfChainPool,
					"topic":       params.OnReceivedTopic,
					"latestLogID": lg.Id,
					"error":       err.Error(),
				}
				log.Logger().WithFields(fields).Error("failed to update filter log")
				alarm.Slack(context.Background(), "failed to update filter log")
				continue
			}
		}

		time.Sleep(10 * time.Second)
	}
}

// HandlePendingOrdersOfFirstStageFromEVM filter the OnReceived event log.
// If this log is found, update order status to confirmed
// action = 2, stage = 1, status = 1 ==> status = 2
//func HandlePendingOrdersOfFirstStageFromEVM() {}

// HandleConfirmedOrdersOfFirstStageFromEVMToTON
// action = 2, stage = 1, status = 2 ==> stage = 2, status = 1
func HandleConfirmedOrdersOfFirstStageFromEVMToTON() {
	order := dao.Order{
		DstChain: params.TONChainID,
		Action:   dao.OrderActionFromEVM,
		Stage:    dao.OrderStag1,
		Status:   dao.OrderStatusConfirmed,
	}
	for {
		for id := uint64(1); ; {
			orders, err := order.GetOldest10ByID(id)
			if err != nil {
				fields := map[string]interface{}{
					"id":    id,
					"order": utils.JSON(order),
					"error": err.Error(),
				}

				log.Logger().WithFields(fields).Error("failed to get confirmed status order of first stage from evm to ton")
				alarm.Slack(context.Background(), "failed to get confirmed status order from evm to ton")
				time.Sleep(5 * time.Second)
				continue
			}

			length := len(orders)
			if length == 0 {
				log.Logger().Info("not found confirmed status order order of first stage from evm to ton", "time", time.Now())
				time.Sleep(10 * time.Second)
				break
			}

			for i, o := range orders {
				if i == length-1 {
					id = o.ID + 1
				}

				request := &tonrouter.BridgeSwapRequest{
					Amount:          o.RelayAmount,
					Slippage:        o.Slippage / 3 * 1,
					TokenOutAddress: o.DstToken,
					Receiver:        o.Receiver,
					OrderID:         o.OrderIDFromContract,
				}

				txParams, err := tonrouter.BridgeSwap(request)
				if err != nil {
					log.Logger().WithField("error", err.Error()).Error("failed to request ton swap")
					alarm.Slack(context.Background(), "failed to request ton swap")
					time.Sleep(1 * time.Second)
					continue
				}
				dstAddr, err := address.ParseAddr(txParams.To)
				if err != nil {
					log.Logger().WithField("address", txParams.To).WithField("error", err.Error()).Error("failed to parse ton address")
					alarm.Slack(context.Background(), "failed to parse ton address")
					continue
				}
				amount, err := tlb.FromNanoTONStr(txParams.Value)
				if err != nil {
					log.Logger().WithField("amount", txParams.Value).WithField("error", err.Error()).Error("failed to parse ton amount")
					alarm.Slack(context.Background(), "failed to parse ton amount")
					continue
				}

				body := &cell.Cell{}
				if err := json.Unmarshal([]byte(fmt.Sprintf(`"%s"`, txParams.Data)), body); err != nil {
					log.Logger().WithField("data", txParams.Data).WithField("error", err.Error()).Error("failed to unmarshal ton data")
					alarm.Slack(context.Background(), "failed to unmarshal ton data")
					continue
				}

				// todo no blocking mode
				// send transaction to chain pool on ton
				t, _, err := tonclient.Wallet().SendWaitTransaction(context.Background(), &wallet.Message{
					Mode: wallet.PayGasSeparately, // pay fees separately (from balance, not from amount)
					InternalMessage: &tlb.InternalMessage{
						Bounce:  true, // return amount in case of processing error
						DstAddr: dstAddr,
						Amount:  amount,
						Body:    body,
					},
				})
				if err != nil {
					log.Logger().WithField("error", err).Error("failed to send transaction to chain pool on ton")
					alarm.Slack(context.Background(), "failed to send transaction to chain pool on ton")
					continue
				}

				log.Logger().Info("transaction sent, confirmed at block, hash:", hex.EncodeToString(t.Hash))

				//balance, err := tonclient.Wallet().GetBalance(context.Background(), block)
				//if err != nil {
				//	log.Logger().WithField("wallet", tonclient.Wallet().WalletAddress()).WithField("error", err).Error("failed to get ton account balance")
				//	alarm.Slack(context.Background(), "failed to get ton account balance")
				//}
				//if balance.Nano().Uint64() < 3000000 {
				//	log.Logger().Info("ton account not enough balance:", balance.String())
				//	alarm.Slack(context.Background(), "ton account not enough balance")
				//}

				update := &dao.Order{
					ID:        o.ID,
					Stage:     dao.OrderStag2,
					Status:    dao.OrderStatusPending,
					OutTxHash: hex.EncodeToString(t.Hash),
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

// HandlePendingOrdersOfSecondSStageFromEVMToTON
// action = 2, stage = 2, status = 1 ==> status = 2
func HandlePendingOrdersOfSecondSStageFromEVMToTON() {
	filterLog := dao.NewFilterLog(params.TONChainID, EventIDEVMToTON)
	for {
		gotLog, err := filterLog.First()
		if err != nil {
			fields := map[string]interface{}{
				"chainID": params.TONChainID,
				"topic":   EventIDEVMToTON,
				"error":   err.Error(),
			}
			log.Logger().WithFields(fields).Error("failed to get filter log info")
			alarm.Slack(context.Background(), "failed to get filter log info")
			time.Sleep(5 * time.Second)
			continue
		}

		logs, err := filter.GetLogs(gotLog.LatestLogID, params.TONChainID, EventIDEVMToTON, uint8(20))
		if err != nil {
			fields := map[string]interface{}{
				"id":      gotLog.LatestLogID,
				"chainID": params.TONChainID,
				"topic":   EventIDEVMToTON,
				"limit":   uint8(20),
				"error":   err.Error(),
			}
			log.Logger().WithFields(fields).Error("failed to get logs")
			alarm.Slack(context.Background(), "failed to get logs")
			continue
		}

		for _, lg := range logs {
			if lg.Id <= gotLog.LatestLogID {
				continue
			}

			logData, err := hex.DecodeString(lg.LogData)
			if err != nil {
				fields := map[string]interface{}{
					"id":      lg.Id,
					"chainID": params.TONChainID,
					"topic":   EventIDTONToEVM,
					"logData": lg.LogData,
					"error":   err.Error(),
				}
				log.Logger().WithFields(fields).Error("failed to decode log data")
				alarm.Slack(context.Background(), "failed to decode log data")
				continue
			}

			body := &cell.Cell{}
			if err := json.Unmarshal(logData, &body); err != nil {
				log.Logger().WithField("logData", lg.LogData).WithField("error", err.Error()).Error("failed to unmarshal data")
				alarm.Slack(context.Background(), "failed to unmarshal data")
				continue
			}

			slice := body.BeginParse()
			orderIDFromContract, err := slice.LoadUInt(64)
			if err != nil {
				log.Logger().WithField("error", err.Error()).Error("failed to load order id")
				alarm.Slack(context.Background(), "failed to load order id from")
				continue
			}

			order, err := dao.NewOrderWithOrderIDFromContract(orderIDFromContract).First()
			if err != nil {
				fields := map[string]interface{}{
					"orderIDFromContract": orderIDFromContract,
					"error":               err.Error(),
				}
				log.Logger().WithFields(fields).Error("failed to get order")
				alarm.Slack(context.Background(), "failed to get order")
				continue
			}
			if order.DstChain != params.TONChainID {
				fields := map[string]interface{}{
					"orderID": order.ID,
					"action":  order.Action,
				}
				log.Logger().WithFields(fields).Error("order dst chain not match")
				alarm.Slack(context.Background(), "order dst chain not match")
				continue
			}
			if order.Action != dao.OrderActionFromEVM {
				fields := map[string]interface{}{
					"orderID": order.ID,
					"action":  order.Action,
				}
				log.Logger().WithFields(fields).Error("order action not match")
				alarm.Slack(context.Background(), "order action not match")
				continue
			}
			if order.Stage != dao.OrderStag2 {
				fields := map[string]interface{}{
					"orderID": order.ID,
					"stage":   order.Stage,
				}
				log.Logger().WithFields(fields).Error("order stage not match")
				alarm.Slack(context.Background(), "order stage not match")
				continue
			}
			if order.Status != dao.OrderStatusPending {
				fields := map[string]interface{}{
					"orderID": order.ID,
					"status":  order.Status,
				}
				log.Logger().WithFields(fields).Error("order status not pending")
				alarm.Slack(context.Background(), "order status not pending")
				continue
			}
			update := &dao.Order{
				ID:     order.ID,
				Status: dao.OrderStatusConfirmed,
			}
			if err := dao.NewOrder().Updates(update); err != nil {
				fields := map[string]interface{}{
					"orderID": order.ID,
					"update":  utils.JSON(update),
					"error":   err.Error(),
				}
				log.Logger().WithFields(fields).Error("failed to update order status")
				alarm.Slack(context.Background(), "failed to update order status")
				continue
			}

			if err := filterLog.UpdateLatestLogID(lg.Id); err != nil {
				fields := map[string]interface{}{
					"chainID":     params.ChainIDOfChainPool,
					"topic":       params.OnReceivedTopic,
					"latestLogID": lg.Id,
					"error":       err.Error(),
				}
				log.Logger().WithFields(fields).Error("failed to update filter log")
				alarm.Slack(context.Background(), "failed to update filter log")
				continue
			}
		}

		time.Sleep(10 * time.Second)
	}
}

func deliver(transactor *tx.Transactor, usdt common.Address, orderID [32]byte, amount *big.Int, receiver common.Address) (txHash common.Hash, err error) {
	txHash, err = transactor.Deliver(orderID, usdt, amount, receiver)
	if err != nil {
		fields := map[string]interface{}{
			"orderID":  hex.EncodeToString(orderID[:]),
			"token":    usdt,
			"amount":   amount,
			"receiver": receiver,
			"error":    err.Error(),
		}
		log.Logger().WithFields(fields).Error("failed to send deliver transaction")
		alarm.Slack(context.Background(), "failed to send deliver transaction")
		return txHash, err
	}
	log.Logger().WithField("hash", txHash).Info("completed send deliver transaction")
	return txHash, nil
}

func deliverAndSwap(request *butter.RouterAndSwapRequest, transactor *tx.Transactor, usdt common.Address, orderID [32]byte, amount *big.Int) (txHash common.Hash, err error) {
	data, err := butter.RouteAndSwap(request)
	if err != nil {
		log.Logger().WithField("error", err.Error()).Error("failed to create router and swap request")
		alarm.Slack(context.Background(), "failed to create router and swap request")
		return txHash, err
	}

	decodeData, err := DecodeData(data.Data)
	if err != nil {
		log.Logger().WithField("data", data.Data).WithField("error", err.Error()).Error("failed to decode call data")
		alarm.Slack(context.Background(), "failed to decode call data")
		return txHash, err
	}

	value, ok := new(big.Int).SetString(utils.TrimHexPrefix(data.Value), 16)
	if !ok {
		log.Logger().WithField("value", utils.TrimHexPrefix(data.Value)).Error("failed to parse string to big int")
		alarm.Slack(context.Background(), "failed to parse string to big int")
		return txHash, err
	}

	txHash, err = transactor.DeliverAndSwap(orderID, Initiator, usdt, amount, decodeData.SwapData, decodeData.BridgeData, decodeData.FeeData, value)
	if err != nil {
		fields := map[string]interface{}{
			"orderID":    hex.EncodeToString(orderID[:]),
			"initiator":  Initiator,
			"token":      usdt,
			"amount":     amount,
			"SwapData":   hex.EncodeToString(decodeData.SwapData),
			"BridgeData": hex.EncodeToString(decodeData.BridgeData),
			"FeeData":    hex.EncodeToString(decodeData.FeeData),
			"value":      value,
			"error":      err.Error(),
		}
		log.Logger().WithFields(fields).Error("failed to send deliver and swap transaction")
		alarm.Slack(context.Background(), "failed to send deliver and swap transaction")
		return txHash, err
	}
	log.Logger().WithField("hash", txHash).Info("completed send deliver and swap transaction")
	return txHash, nil
}
