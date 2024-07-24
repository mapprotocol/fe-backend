package task

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/mapprotocol/fe-backend/dao"
	"github.com/mapprotocol/fe-backend/params"
	"github.com/mapprotocol/fe-backend/resource/tonclient"
	"github.com/mapprotocol/fe-backend/third-party/butter"
	"github.com/mapprotocol/fe-backend/third-party/tonrouter"
	"github.com/mapprotocol/fe-backend/utils"
	"github.com/mapprotocol/fe-backend/utils/alarm"
	"github.com/mapprotocol/fe-backend/utils/tx"
	"github.com/xssnick/tonutils-go/ton/wallet"
	"github.com/xssnick/tonutils-go/tvm/cell"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/mapprotocol/fe-backend/resource/log"
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tlb"
)

var FeeRate = big.NewInt(7) // 7/1000

const (
	EventIDTONToEVM = "34a7e0e8"
	EventIDEVMToTON = "1a6c0a51"
)

func HandlePendingOrdersOfFirstStageFromTONToEVM() error {
	filterLog := dao.NewFilterLog(params.ChainIDOfTon, params.SwapLogTopic)
	gotLog, err := filterLog.First()
	if err != nil {
		fields := map[string]interface{}{
			"chainID": params.ChainIDOfTon,
			"topic":   params.SwapLogTopic,
			"error":   err.Error(),
		}
		log.Logger().WithFields(fields).Error("failed to get filter log info")
		alarm.Slack(context.Background(), "failed to get filter log info")
		return err
	}

	master, err := tonclient.Client().CurrentMasterchainInfo(context.Background()) // we fetch block just to trigger chain proof check
	if err != nil {
		log.Logger().WithField("error", err.Error()).Error("failed to get master chain info")
		return err
	}
	log.Logger().Info("master proof checks are completed successfully, now communication is 100% safe!, master: ", master)

	// address on which we are accepting payments
	treasuryAddress := address.MustParseAddr(params.ChainPoolAddressOnTON) // todo contract name
	lastProcessedLT := gotLog.LatestLogID
	// channel with new transactions
	transactions := make(chan *tlb.Transaction)

	// it is a blocking call, so we start it asynchronously
	go tonclient.Client().SubscribeOnTransactions(context.Background(), treasuryAddress, lastProcessedLT, transactions) // 47840872000005

	log.Logger().Info("waiting for transfers...")

	// listen for new transactions from channel
	for t := range transactions {
		out := t.IO.Out
		if out == nil {
			continue
		}
		messages, err := out.ToSlice()
		if err != nil {
			log.Logger().WithField("error", err.Error()).Error("failed to get out messages")
			continue
		}
		for _, msg := range messages {
			if msg.MsgType != tlb.MsgTypeExternalOut {
				continue
			}

			externalOut := msg.AsExternalOut()
			payload := externalOut.Payload()
			if payload == nil {
				fmt.Println("payload is ni")
				continue
			}

			slice := payload.BeginParse()
			if slice == nil {
				continue
			}

			if strings.Contains(strings.ToLower(externalOut.DstAddr.String()), EventIDTONToEVM) {
				orderID, err := slice.LoadUInt(64)
				if err != nil {
					log.Logger().WithField("error", err.Error()).Error("failed to load order id")
					alarm.Slack(context.Background(), "failed to load order id from")
					continue
				}
				//from, err := slice.LoadRef()
				//if err != nil {
				//	log.Logger().WithField("error", err.Error()).Error("failed to load from ref")
				//	alarm.Slack(context.Background(), "failed to load from ref")
				//	continue
				//}
				//to, err := slice.LoadRef()
				//if err != nil {
				//	log.Logger().WithField("error", err.Error()).Error("failed to load to ref")
				//	alarm.Slack(context.Background(), "failed to load to ref")
				//	continue
				//}
				//sender, err := from.LoadAddr()
				//if err != nil {
				//	log.Logger().WithField("error", err.Error()).Error("failed to load sender")
				//	alarm.Slack(context.Background(), "failed to load sender")
				//	return err
				//}
				//srcToken, err := from.LoadAddr()
				//if err != nil {
				//	log.Logger().WithField("error", err.Error()).Error("failed to load src token")
				//	alarm.Slack(context.Background(), "failed to load rc token")
				//	return err
				//}
				//amountIn, err := slice.LoadUInt(32)
				//if err != nil {
				//	log.Logger().WithField("error", err.Error()).Error("failed to load amount in")
				//	alarm.Slack(context.Background(), "failed to load amount in")
				//	continue
				//}
				//toChainID, err := slice.LoadUInt(32)
				//if err != nil {
				//	log.Logger().WithField("error", err.Error()).Error("failed to load to chain id")
				//	alarm.Slack(context.Background(), "failed to load to chain id")
				//	continue
				//}
				//r, err := to.LoadBigUInt(160)
				//if err != nil {
				//	log.Logger().WithField("error", err.Error()).Error("failed to load receiver")
				//	alarm.Slack(context.Background(), "failed to load receiver")
				//	continue
				//}
				//receiver := "0x" + hex.EncodeToString(r.Bytes())
				//tokenOUt, err := to.LoadBigUInt(160)
				//if err != nil {
				//	log.Logger().WithField("error", err.Error()).Error("failed to load token out address")
				//	alarm.Slack(context.Background(), "failed to load token out address")
				//	continue
				//}
				//tokenOutAddress := "0x" + hex.EncodeToString(tokenOUt.Bytes())
				amount, err := slice.LoadUInt(32)
				if err != nil {
					log.Logger().WithField("error", err.Error()).Error("failed to load jetton amount")
					alarm.Slack(context.Background(), "failed to load jetton amount")
					continue
				}

				order, err := dao.NewOrderWithID(orderID).First()
				if err != nil {
					fields := map[string]interface{}{
						"orderID": orderID,
						"error":   err.Error(),
					}
					log.Logger().WithFields(fields).Error("failed to get order")
					alarm.Slack(context.Background(), "failed to get order")
					continue
				}
				if order.Action != dao.OrderActionToEVM {
					fields := map[string]interface{}{
						"orderID": orderID,
						"action":  order.Action,
					}
					log.Logger().WithFields(fields).Error("order action not match")
					alarm.Slack(context.Background(), "order action not match")
					continue
				}
				if order.Stage != dao.OrderStag1 {
					fields := map[string]interface{}{
						"orderID": orderID,
						"stage":   order.Stage,
					}
					log.Logger().WithFields(fields).Error("order stage not match")
					alarm.Slack(context.Background(), "order stage not match")
					continue
				}
				if order.Status != dao.OrderStatusPending {
					fields := map[string]interface{}{
						"orderID": orderID,
						"status":  order.Status,
					}
					log.Logger().WithFields(fields).Error("order status not pending")
					alarm.Slack(context.Background(), "order status not pending")
					continue
				}

				_, afterAmount := fees(new(big.Int).SetUint64(amount), FeeRate)

				update := &dao.Order{
					ID:       order.ID,
					InAmount: afterAmount.String(),
					Status:   dao.OrderStatusConfirmed,
				}
				if err := dao.NewOrder().Updates(update); err != nil {
					fields := map[string]interface{}{
						"orderID": orderID,
						"update":  utils.JSON(update),
						"error":   err.Error(),
					}
					log.Logger().WithFields(fields).Error("failed to update order status")
					alarm.Slack(context.Background(), "failed to update order status")
					return err
				}
			} else if strings.Contains(strings.ToLower(externalOut.DstAddr.String()), EventIDEVMToTON) {
				orderID, err := slice.LoadUInt(64)
				if err != nil {
					log.Logger().WithField("error", err.Error()).Error("failed to load order id")
					alarm.Slack(context.Background(), "failed to load order id from")
					continue
				}

				order, err := dao.NewOrderWithID(orderID).First()
				if err != nil {
					fields := map[string]interface{}{
						"orderID": orderID,
						"error":   err.Error(),
					}
					log.Logger().WithFields(fields).Error("failed to get order")
					alarm.Slack(context.Background(), "failed to get order")
					continue
				}
				if order.Action != dao.OrderActionFromEVM {
					fields := map[string]interface{}{
						"orderID": orderID,
						"action":  order.Action,
					}
					log.Logger().WithFields(fields).Error("order action not match")
					alarm.Slack(context.Background(), "order action not match")
					continue
				}
				if order.Stage != dao.OrderStag2 {
					fields := map[string]interface{}{
						"orderID": orderID,
						"stage":   order.Stage,
					}
					log.Logger().WithFields(fields).Error("order stage not match")
					alarm.Slack(context.Background(), "order stage not match")
					continue
				}
				if order.Status != dao.OrderStatusPending {
					fields := map[string]interface{}{
						"orderID": orderID,
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
						"orderID": orderID,
						"update":  utils.JSON(update),
						"error":   err.Error(),
					}
					log.Logger().WithFields(fields).Error("failed to update order status")
					alarm.Slack(context.Background(), "failed to update order status")
					return err
				}
			}
		}

		// update last processed lt and save it in db
		if err := filterLog.UpdateLatestLogID(t.LT); err != nil {
			fields := map[string]interface{}{
				"chainID":     params.ChainIDOfChainPool,
				"topic":       params.OnReceivedTopic,
				"latestLogID": t.LT,
				"error":       err.Error(),
			}
			log.Logger().WithFields(fields).Error("failed to update filter log")
			alarm.Slack(context.Background(), "failed to update filter log")
			return err
		}
		lastProcessedLT = t.LT
	}

	return errors.New("transaction listening unexpectedly finished")
}

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
				log.Logger().Info("not found confirmed status order from ton to evm", "time", time.Now())
				time.Sleep(10 * time.Second)
				break
			}

			for i, o := range orders {
				if i == length-1 {
					id = o.ID + 1
				}

				amountFloat, ok := new(big.Float).SetString(o.InAmount)
				if !ok {
					fields := map[string]interface{}{
						"orderId": o.ID,
						"amount":  o.InAmount,
						"error":   err,
					}
					log.Logger().WithFields(fields).Error("failed to parse string to big float")
					alarm.Slack(context.Background(), "failed to parse string to big float")
					continue
				}
				amountFloat = new(big.Float).Quo(amountFloat, big.NewFloat(1e6)) // todo
				request := &butter.RouterAndSwapRequest{
					//FromChainID:     o.SrcChain,
					FromChainID: params.ChainIDOfChainPool, // todo
					ToChainID:   o.DstChain,
					//Amount:          o.InAmount,  //
					Amount:          amountFloat.String(), // todo
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
						"token":      params.USDTOfChainPool,
						"amount":     value,
						"SwapData":   hex.EncodeToString(decodeData.SwapData),
						"BridgeData": hex.EncodeToString(decodeData.BridgeData),
						"FeeData":    hex.EncodeToString(decodeData.FeeData),
						"value":      value,
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

				amountFloat, ok := new(big.Float).SetString(o.InAmount)
				if !ok {
					fields := map[string]interface{}{
						"orderId": o.ID,
						"amount":  o.InAmount,
						"error":   err,
					}
					log.Logger().WithFields(fields).Error("failed to parse string to big float")
					alarm.Slack(context.Background(), "failed to parse string to big float")
					continue
				}
				amountFloat = new(big.Float).Quo(amountFloat, big.NewFloat(1e6)) // todo

				request := &tonrouter.BridgeSwapRequest{
					Amount:          amountFloat.String(),
					Slippage:        o.Slippage / 3 * 1,
					TokenOutAddress: o.DstToken,
					Receiver:        o.Receiver,
					OrderID:         o.ID,
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
				t, block, err := tonclient.Wallet().SendWaitTransaction(context.Background(), &wallet.Message{
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
					return
				}

				log.Logger().Info("transaction sent, confirmed at block, hash:", base64.StdEncoding.EncodeToString(t.Hash))

				balance, err := tonclient.Wallet().GetBalance(context.Background(), block)
				if err != nil {
					log.Logger().WithField("wallet", tonclient.Wallet().WalletAddress()).WithField("error", err).Error("failed to get ton account balance")
					alarm.Slack(context.Background(), "failed to get ton account balance")
				}
				if balance.Nano().Uint64() < 3000000 {
					log.Logger().Info("ton account not enough balance:", balance.String())
					alarm.Slack(context.Background(), "ton account not enough balance")
				}

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
