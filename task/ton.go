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

var (
	Big0         = big.NewInt(0)
	EmptyAddress = common.Address{}
)

var (
	sender           string
	isMultiChainPool bool
)

func Init() {
	sender = viper.GetStringMapString("chainPool")["sender"]
	if utils.IsEmpty(sender) {
		panic("chainPool:sender is empty")
	}
	isMultiChainPool = viper.GetBool("isMultiChainPool")

	blog.Println("chainPool:sender: ", sender)
	blog.Println("isMultiChainPool: ", isMultiChainPool)
}

// HandlePendingOrdersOfFirstStageFromTONToEVM
// create order from ton to evm( action=1, stage=1, status=4(TxConfirmed))
func HandlePendingOrdersOfFirstStageFromTONToEVM() {
	chainID := params.TONChainID
	topic := EventIDTONToEVM
	filterLog := dao.NewFilterLog(chainID, topic)
	for {
		gotLog, err := filterLog.First()
		if err != nil {
			fields := map[string]interface{}{
				"chainID": chainID,
				"topic":   topic,
				"error":   err.Error(),
			}
			log.Logger().WithFields(fields).Error("failed to get filter log info")
			alarm.Slack(context.Background(), "failed to get filter log info")
			time.Sleep(5 * time.Second)
			continue
		}

		logs, err := filter.GetLogs(gotLog.LatestLogID, chainID, topic, uint8(20))
		if err != nil {
			fields := map[string]interface{}{
				"id":      gotLog.LatestLogID,
				"chainID": chainID,
				"topic":   topic,
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
					"logID":   lg.Id,
					"chainID": chainID,
					"topic":   topic,
					"logData": lg.LogData,
					"error":   err.Error(),
				}
				log.Logger().WithFields(fields).Error("failed to decode log data")
				alarm.Slack(context.Background(), "failed to decode log data")
				UpdateLogID(chainID, topic, lg.Id)
				continue
			}

			body := &cell.Cell{}
			if err := json.Unmarshal(logData, &body); err != nil {
				fields := map[string]interface{}{
					"logID":   lg.Id,
					"chainID": chainID,
					"topic":   topic,
					"logData": hex.EncodeToString(logData),
					"error":   err.Error(),
				}
				log.Logger().WithFields(fields).Error("failed to unmarshal log data")
				alarm.Slack(context.Background(), "failed to unmarshal log data")
				UpdateLogID(chainID, topic, lg.Id)
				continue
			}
			slice := body.BeginParse()
			orderID, err := slice.LoadUInt(64)
			if err != nil {
				fields := map[string]interface{}{
					"logID":   lg.Id,
					"chainID": chainID,
					"topic":   topic,
					"error":   err.Error(),
				}
				log.Logger().WithFields(fields).Error("failed to load order id")
				alarm.Slack(context.Background(), "failed to load order id from")
				UpdateLogID(chainID, topic, lg.Id)
				continue
			}
			from, err := slice.LoadRef()
			if err != nil {
				fields := map[string]interface{}{
					"logID":   lg.Id,
					"chainID": chainID,
					"topic":   topic,
					"error":   err.Error(),
				}
				log.Logger().WithFields(fields).Error("failed to load from ref")
				alarm.Slack(context.Background(), "failed to load from ref")
				UpdateLogID(chainID, topic, lg.Id)
				continue
			}
			to, err := slice.LoadRef()
			if err != nil {
				fields := map[string]interface{}{
					"logID":   lg.Id,
					"chainID": chainID,
					"topic":   topic,
					"error":   err.Error(),
				}
				log.Logger().WithFields(fields).Error("failed to load to ref")
				alarm.Slack(context.Background(), "failed to load to ref")
				UpdateLogID(chainID, topic, lg.Id)
				continue
			}
			srcChain, err := from.LoadUInt(64)
			if err != nil {
				fields := map[string]interface{}{
					"logID":   lg.Id,
					"chainID": chainID,
					"topic":   topic,
					"error":   err.Error(),
				}
				log.Logger().WithFields(fields).Error("failed to load from chain id")
				alarm.Slack(context.Background(), "failed to load from chain id")
				UpdateLogID(chainID, topic, lg.Id)
				continue
			}
			sender, err := from.LoadAddr()
			if err != nil {
				fields := map[string]interface{}{
					"logID":   lg.Id,
					"chainID": chainID,
					"topic":   topic,
					"error":   err.Error(),
				}
				log.Logger().WithFields(fields).Error("failed to load sender")
				alarm.Slack(context.Background(), "failed to load sender")
				UpdateLogID(chainID, topic, lg.Id)
				continue
			}
			srcToken, err := from.LoadAddr()
			if err != nil {
				fields := map[string]interface{}{
					"logID":   lg.Id,
					"chainID": chainID,
					"topic":   topic,
					"error":   err.Error(),
				}
				log.Logger().WithFields(fields).Error("failed to load src token")
				alarm.Slack(context.Background(), "failed to load rc token")
				UpdateLogID(chainID, topic, lg.Id)
				continue
			}
			inAmount, err := from.LoadUInt(64)
			if err != nil {
				fields := map[string]interface{}{
					"logID":   lg.Id,
					"chainID": chainID,
					"topic":   topic,
					"error":   err.Error(),
				}
				log.Logger().WithFields(fields).Error("failed to load amount in")
				alarm.Slack(context.Background(), "failed to load amount in")
				UpdateLogID(chainID, topic, lg.Id)
				continue
			}
			slippage, err := from.LoadUInt(16)
			if err != nil {
				fields := map[string]interface{}{
					"logID":   lg.Id,
					"chainID": chainID,
					"topic":   topic,
					"error":   err.Error(),
				}
				log.Logger().WithFields(fields).Error("failed to load slippage")
				alarm.Slack(context.Background(), "failed to load slippage")
				UpdateLogID(chainID, topic, lg.Id)
				continue
			}
			dstChain, err := to.LoadUInt(64)
			if err != nil {
				fields := map[string]interface{}{
					"logID":   lg.Id,
					"chainID": chainID,
					"topic":   topic,
					"error":   err.Error(),
				}
				log.Logger().WithFields(fields).Error("failed to load to chain id")
				alarm.Slack(context.Background(), "failed to load to chain id")
				UpdateLogID(chainID, topic, lg.Id)
				continue
			}
			receiverRef, err := to.LoadRef()
			if err != nil {
				fields := map[string]interface{}{
					"logID":   lg.Id,
					"chainID": chainID,
					"topic":   topic,
					"error":   err.Error(),
				}
				log.Logger().WithFields(fields).Error("failed to load receiver ref")
				alarm.Slack(context.Background(), "failed to load receiver ref")
				UpdateLogID(chainID, topic, lg.Id)
				continue
			}
			heightBitsReceiver, err := receiverRef.LoadBigUInt(256)
			if err != nil {
				fields := map[string]interface{}{
					"logID":   lg.Id,
					"chainID": chainID,
					"topic":   topic,
					"error":   err.Error(),
				}
				log.Logger().WithFields(fields).Error("failed to load height bits receiver")
				alarm.Slack(context.Background(), "failed to load height receiver")
				UpdateLogID(chainID, topic, lg.Id)
				continue
			}
			lowBitsReceiver, err := receiverRef.LoadBigUInt(256)
			if err != nil {
				fields := map[string]interface{}{
					"logID":   lg.Id,
					"chainID": chainID,
					"topic":   topic,
					"error":   err.Error(),
				}
				log.Logger().WithFields(fields).Error("failed to load low bits receiver")
				alarm.Slack(context.Background(), "failed to load low receiver")
				UpdateLogID(chainID, topic, lg.Id)
				continue
			}
			receiver := new(big.Int).Or(new(big.Int).Lsh(heightBitsReceiver, 256), lowBitsReceiver)

			dstTokenRef, err := to.LoadRef()
			if err != nil {
				fields := map[string]interface{}{
					"logID":   lg.Id,
					"chainID": chainID,
					"topic":   topic,
					"error":   err.Error(),
				}
				log.Logger().WithFields(fields).Error("failed to load token out address ref")
				alarm.Slack(context.Background(), "failed to load token out address ref")
				UpdateLogID(chainID, topic, lg.Id)
				continue
			}
			heightBitsDstToken, err := dstTokenRef.LoadBigUInt(256)
			if err != nil {
				fields := map[string]interface{}{
					"logID":   lg.Id,
					"chainID": chainID,
					"topic":   topic,
					"error":   err.Error(),
				}
				log.Logger().WithFields(fields).Error("failed to load height bits token out address")
				alarm.Slack(context.Background(), "failed to load height bits token out address")
				UpdateLogID(chainID, topic, lg.Id)
				continue
			}
			lowBitsDstToken, err := dstTokenRef.LoadBigUInt(256)
			if err != nil {
				fields := map[string]interface{}{
					"logID":   lg.Id,
					"chainID": chainID,
					"topic":   topic,
					"error":   err.Error(),
				}
				log.Logger().WithFields(fields).Error("failed to load low bits token out address")
				alarm.Slack(context.Background(), "failed to load low bits token out address")
				UpdateLogID(chainID, topic, lg.Id)
				continue
			}
			dstToken := new(big.Int).Or(new(big.Int).Lsh(heightBitsDstToken, 256), lowBitsDstToken)

			relayAmount, err := slice.LoadUInt(32)
			if err != nil {
				fields := map[string]interface{}{
					"logID":   lg.Id,
					"chainID": chainID,
					"topic":   topic,
					"error":   err.Error(),
				}
				log.Logger().WithFields(fields).Error("failed to load jetton amount")
				alarm.Slack(context.Background(), "failed to load jetton amount")
				UpdateLogID(chainID, topic, lg.Id)
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
				Status:              dao.OrderStatusTxConfirmed,
				Slippage:            slippage,
			}

			if err := order.Create(); err != nil {
				fields := map[string]interface{}{
					"logID":   lg.Id,
					"chainID": chainID,
					"topic":   topic,
					"order":   utils.JSON(order),
					"error":   err.Error(),
				}
				log.Logger().WithFields(fields).Error("failed to create order")
				alarm.Slack(context.Background(), "failed to update order status")
				UpdateLogID(chainID, topic, lg.Id)
				continue
			}

			UpdateLogID(chainID, topic, lg.Id)
		}

		time.Sleep(10 * time.Second)
	}
}

// HandleConfirmedOrdersOfFirstStageFromTONToEVM
// action=1, stage=1, status=4(TxConfirmed) ==> stage=2, status=1(TxPrepareSend) ==> stage=2, status=2(TxSent)
func HandleConfirmedOrdersOfFirstStageFromTONToEVM() {
	order := dao.Order{
		SrcChain: params.TONChainID,
		Action:   dao.OrderActionToEVM,
		Stage:    dao.OrderStag1,
		Status:   dao.OrderStatusTxConfirmed,
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

				usdt := params.USDTOfChainPool
				decimal := params.USDTDecimalOfChainPool
				chainIDOfChainPool := params.ChainIDOfChainPool
				chainInfo := &dao.ChainPool{}
				if isMultiChainPool && o.SrcChain == params.ChainIDOfEthereum {
					usdt = params.USDTOfEthereum
					decimal = params.USDTDecimalOfEthereum
					chainIDOfChainPool = params.ChainIDOfEthereum
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
				if o.DstChain == chainIDOfChainPool && strings.ToLower(o.DstToken) == strings.ToLower(usdt) {
					update := &dao.Order{
						Stage:  dao.OrderStag2,
						Status: dao.OrderStatusTxPrepareSend,
					}
					if err := dao.NewOrderWithID(o.ID).Updates(update); err != nil {
						fields := map[string]interface{}{
							"orderId": o.ID,
							"update":  utils.JSON(update),
							"error":   err,
						}
						log.Logger().WithFields(fields).WithField("error", err.Error()).Error("failed to update order status")
						alarm.Slack(context.Background(), "failed to update order status")
						time.Sleep(5 * time.Second)
						continue
					}

					txHash, err = deliver(transactor, common.HexToAddress(usdt), orderID, amountInt, common.HexToAddress(o.Receiver), Big0, EmptyAddress)
					if err != nil {
						log.Logger().WithField("error", err.Error()).Error("failed to send deliver transaction")
						alarm.Slack(context.Background(), "failed to send deliver transaction")
						time.Sleep(5 * time.Second)
						continue
					}
				} else {
					request := &butter.RouterAndSwapRequest{
						FromChainID:     params.ChainIDOfChainPool,
						ToChainID:       o.DstChain,
						Amount:          o.RelayAmount,
						TokenInAddress:  params.USDTOfChainPool,
						TokenOutAddress: o.DstToken,
						Type:            SwapType,
						Slippage:        o.Slippage / 3 * 2,
						From:            sender, // todo
						Receiver:        o.Receiver,
					}
					data, err := butter.RouteAndSwap(request)
					if err != nil {
						log.Logger().WithField("request", utils.JSON(request)).WithField("error", err.Error()).Error("failed to create router and swap request")
						alarm.Slack(context.Background(), "failed to create router and swap request")
						continue
					}

					decodeData, err := DecodeData(data.Data)
					if err != nil {
						log.Logger().WithField("data", data.Data).WithField("error", err.Error()).Error("failed to decode call data")
						alarm.Slack(context.Background(), "failed to decode call data")
						continue
					}

					value, ok := new(big.Int).SetString(utils.TrimHexPrefix(data.Value), 16)
					if !ok {
						log.Logger().WithField("value", utils.TrimHexPrefix(data.Value)).Error("failed to parse string to big int")
						alarm.Slack(context.Background(), "failed to parse string to big int")
						continue
					}

					update := &dao.Order{
						Stage:  dao.OrderStag2,
						Status: dao.OrderStatusTxPrepareSend,
					}
					if err := dao.NewOrderWithID(o.ID).Updates(update); err != nil {
						fields := map[string]interface{}{
							"orderId": o.ID,
							"update":  utils.JSON(update),
							"error":   err,
						}
						log.Logger().WithFields(fields).Error("failed to update order status")
						alarm.Slack(context.Background(), "failed to update order status")
						time.Sleep(5 * time.Second)
						continue
					}

					txHash, err = deliverAndSwap(transactor, common.HexToAddress(usdt), orderID, amountInt, decodeData, Big0, EmptyAddress, value)
					if err != nil {
						log.Logger().WithField("error", err.Error()).Error("failed to send deliver and swap transaction")
						alarm.Slack(context.Background(), "failed to send deliver and swap transaction")
						time.Sleep(5 * time.Second)
						continue
					}
				}

				update := &dao.Order{
					Stage:     dao.OrderStag2,
					Status:    dao.OrderStatusTxSent,
					OutTxHash: txHash.String(),
				}
				if err := dao.NewOrderWithID(o.ID).Updates(update); err != nil {
					fields := map[string]interface{}{
						"orderId": o.ID,
						"update":  utils.JSON(update),
						"error":   err,
					}
					log.Logger().WithFields(fields).Error("failed to update order status")
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
// action=1, stage=2, status=2(TxSent) ==> status=3(TxFailed)/4(TxConfirmed)
func HandlePendingOrdersOfSecondStageFromTONToEVM() {
	order := dao.Order{
		SrcChain: params.TONChainID,
		Action:   dao.OrderActionToEVM,
		Stage:    dao.OrderStag2,
		Status:   dao.OrderStatusTxSent,
	}
	for {
		for id := uint64(1); ; {
			orders, err := order.GetOldest10ByID(id)
			if err != nil {
				fields := map[string]interface{}{
					"id":    id,
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
					fields := map[string]interface{}{
						"chainID": chainInfo.ChainRPC,
						"txHash":  o.OutTxHash,
						"error":   err,
					}
					log.Logger().WithFields(fields).Error("failed to get transaction status")
					alarm.Slack(context.Background(), "failed to get transaction status")
					time.Sleep(5 * time.Second)
					continue
				}

				update := &dao.Order{
					Status: dao.OrderStatusTxConfirmed,
				}
				if status == types.ReceiptStatusFailed {
					update.Status = dao.OrderStatusTxFailed
				}
				if err := dao.NewOrderWithID(o.ID).Updates(update); err != nil {
					fields := map[string]interface{}{
						"orderId": o.ID,
						"status":  update.Status,
						"error":   err,
					}
					log.Logger().WithFields(fields).Error("failed to update order status")
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
// action=1, stage=2, status=4(TxConfirmed) ==> status=5(Completed)
func HandleConfirmedOrdersOfSecondStageFromTONToEVM() {
	chainID := params.ChainIDOfChainPool
	topic := EventIDTONToEVM
	filterLog := dao.NewFilterLog(chainID, topic)
	for {
		gotLog, err := filterLog.First()
		if err != nil {
			fields := map[string]interface{}{
				"chainID": chainID,
				"topic":   topic,
				"error":   err.Error(),
			}
			log.Logger().WithFields(fields).Error("failed to get filter log info")
			alarm.Slack(context.Background(), "failed to get filter log info")
			time.Sleep(5 * time.Second)
			continue
		}

		logs, err := filter.GetLogs(gotLog.LatestLogID, chainID, topic, uint8(20))
		if err != nil {
			fields := map[string]interface{}{
				"id":      gotLog.LatestLogID,
				"chainID": chainID,
				"topic":   topic,
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
					"chainID": chainID,
					"topic":   topic,
					"logData": lg.LogData,
					"error":   err.Error(),
				}
				log.Logger().WithFields(fields).Error("failed to decode log data")
				alarm.Slack(context.Background(), "failed to decode log data")
				UpdateLogID(chainID, topic, lg.Id)
				continue
			}
			deliverAndSwap, err := UnpackDeliverAndSwap(logData)
			if err != nil {
				fields := map[string]interface{}{
					"id":      lg.Id,
					"chainID": chainID,
					"topic":   topic,
					"logData": lg.LogData,
					"error":   err.Error(),
				}
				log.Logger().WithFields(fields).Error("failed to unpack deliver and swap log data")
				alarm.Slack(context.Background(), "failed to unpack deliver and swap log data")
				UpdateLogID(chainID, topic, lg.Id)
				continue
			}
			_ = deliverAndSwap

			UpdateLogID(chainID, topic, lg.Id)
			continue
		}

		time.Sleep(10 * time.Second)
	}
}

// HandlePendingOrdersOfFirstStageFromEVM filter the OnReceived event log.
// If this log is found, update order status to confirmed
// action = 2, stage = 1, status = 2
// create order from evm ( action=2, stage=1, status=4(TxConfirmed))
// todo multi chain pool
//func HandlePendingOrdersOfFirstStageFromEVM() {}

// HandleConfirmedOrdersOfFirstStageFromEVMToTON
// action = 2, stage = 1, status = 2 ==> stage = 2, status = 1
// action=2, stage=1, status=4(TxConfirmed) ==> stage=2, status=1(TxPrepareSend) ==> stage=2, status=2(TxSent)
func HandleConfirmedOrdersOfFirstStageFromEVMToTON() {
	order := dao.Order{
		DstChain: params.TONChainID,
		Action:   dao.OrderActionFromEVM,
		Stage:    dao.OrderStag1,
		Status:   dao.OrderStatusTxConfirmed,
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
					log.Logger().WithField("request", utils.JSON(request)).WithField("error", err.Error()).Error("failed to request ton swap")
					alarm.Slack(context.Background(), "failed to request ton swap")
					time.Sleep(5 * time.Second)
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

				update := &dao.Order{
					Stage:  dao.OrderStag2,
					Status: dao.OrderStatusTxPrepareSend,
				}
				if err := dao.NewOrderWithID(o.ID).Updates(update); err != nil {
					fields := map[string]interface{}{
						"id":     o.ID,
						"update": utils.JSON(update),
						"error":  err.Error(),
					}
					log.Logger().WithFields(fields).Error("failed to update order status")
					alarm.Slack(context.Background(), "failed to update order status")
					time.Sleep(5 * time.Second)
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
					fields := map[string]interface{}{
						"id":      o.ID,
						"dstAddr": dstAddr,
						"amount":  amount,
						"txData":  txParams.Data,
						"error":   err.Error(),
					}
					log.Logger().WithFields(fields).Error("failed to send transaction to chain pool on ton")
					alarm.Slack(context.Background(), "failed to send transaction to chain pool on ton")
					continue
				}

				log.Logger().Info("transaction sent, confirmed at block, hash:", hex.EncodeToString(t.Hash))

				update = &dao.Order{
					Stage:     dao.OrderStag2,
					Status:    dao.OrderStatusTxSent,
					OutTxHash: hex.EncodeToString(t.Hash),
				}
				if err := dao.NewOrderWithID(o.ID).Updates(update); err != nil {
					fields := map[string]interface{}{
						"id":     o.ID,
						"update": utils.JSON(update),
						"error":  err.Error(),
					}
					log.Logger().WithFields(fields).Error("failed to update order status")
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
// action=2, stage=2, status=2(TxSent) ==> status=4(TxConfirmed)
func HandlePendingOrdersOfSecondSStageFromEVMToTON() {
	chainID := params.TONChainID
	topic := EventIDEVMToTON
	filterLog := dao.NewFilterLog(chainID, topic)
	for {
		gotLog, err := filterLog.First()
		if err != nil {
			fields := map[string]interface{}{
				"chainID": chainID,
				"topic":   topic,
				"error":   err.Error(),
			}
			log.Logger().WithFields(fields).Error("failed to get filter log info")
			alarm.Slack(context.Background(), "failed to get filter log info")
			time.Sleep(5 * time.Second)
			continue
		}

		logs, err := filter.GetLogs(gotLog.LatestLogID, chainID, topic, uint8(20))
		if err != nil {
			fields := map[string]interface{}{
				"id":      gotLog.LatestLogID,
				"chainID": chainID,
				"topic":   topic,
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
					"chainID": chainID,
					"topic":   topic,
					"logData": lg.LogData,
					"error":   err.Error(),
				}
				log.Logger().WithFields(fields).Error("failed to decode log data")
				alarm.Slack(context.Background(), "failed to decode log data")
				UpdateLogID(chainID, topic, lg.Id)
				continue
			}

			body := &cell.Cell{}
			if err := json.Unmarshal(logData, &body); err != nil {
				fields := map[string]interface{}{
					"id":      lg.Id,
					"chainID": chainID,
					"topic":   topic,
					"logData": hex.EncodeToString(logData),
					"error":   err.Error(),
				}
				log.Logger().WithFields(fields).Error("failed to unmarshal data")
				alarm.Slack(context.Background(), "failed to unmarshal data")
				UpdateLogID(chainID, topic, lg.Id)
				continue
			}

			slice := body.BeginParse()
			orderIDFromContract, err := slice.LoadUInt(64)
			if err != nil {
				fields := map[string]interface{}{
					"id":      lg.Id,
					"chainID": chainID,
					"topic":   topic,
					"error":   err.Error(),
				}
				log.Logger().WithFields(fields).Error("failed to load order id")
				alarm.Slack(context.Background(), "failed to load order id from")
				UpdateLogID(chainID, topic, lg.Id)
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
				UpdateLogID(chainID, topic, lg.Id)
				continue
			}
			if order.DstChain != params.TONChainID {
				fields := map[string]interface{}{
					"orderID":  order.ID,
					"dstChain": order.DstChain,
				}
				log.Logger().WithFields(fields).Error("order dst chain not match")
				alarm.Slack(context.Background(), "order dst chain not match")
				UpdateLogID(chainID, topic, lg.Id)
				continue
			}
			if order.Action != dao.OrderActionFromEVM {
				fields := map[string]interface{}{
					"orderID": order.ID,
					"action":  order.Action,
				}
				log.Logger().WithFields(fields).Error("order action not match")
				alarm.Slack(context.Background(), "order action not match")
				UpdateLogID(chainID, topic, lg.Id)
				continue
			}
			if order.Stage != dao.OrderStag2 {
				fields := map[string]interface{}{
					"orderID": order.ID,
					"stage":   order.Stage,
				}
				log.Logger().WithFields(fields).Error("order stage not match")
				alarm.Slack(context.Background(), "order stage not match")
				UpdateLogID(chainID, topic, lg.Id)
				continue
			}
			if order.Status != dao.OrderStatusTxSent {
				fields := map[string]interface{}{
					"orderID": order.ID,
					"status":  order.Status,
				}
				log.Logger().WithFields(fields).Error("order status not tx sent")
				alarm.Slack(context.Background(), "order status not tx sent")
				UpdateLogID(chainID, topic, lg.Id)
				continue
			}
			update := &dao.Order{
				ID:     order.ID,
				Status: dao.OrderStatusTxConfirmed,
			}
			if err := dao.NewOrder().Updates(update); err != nil {
				fields := map[string]interface{}{
					"orderID": order.ID,
					"status":  dao.OrderStatusTxConfirmed,
					"error":   err.Error(),
				}
				log.Logger().WithFields(fields).Error("failed to update order status")
				alarm.Slack(context.Background(), "failed to update order status")
				UpdateLogID(chainID, topic, lg.Id)
				continue
			}

			UpdateLogID(chainID, topic, lg.Id)
		}

		time.Sleep(10 * time.Second)
	}
}

func deliver(transactor *tx.Transactor, token common.Address, orderID [32]byte, amount *big.Int, receiver common.Address, fee *big.Int, feeReceiver common.Address) (txHash common.Hash, err error) {
	txHash, err = transactor.Deliver(orderID, token, amount, receiver, fee, feeReceiver)
	if err != nil {
		fields := map[string]interface{}{
			"orderID":     hex.EncodeToString(orderID[:]),
			"token":       token,
			"amount":      amount,
			"receiver":    receiver,
			"fee":         fee,
			"feeReceiver": feeReceiver,
			"error":       err.Error(),
		}
		log.Logger().WithFields(fields).Error("failed to send deliver transaction")
		return txHash, err
	}
	log.Logger().WithField("hash", txHash).Info("completed send deliver transaction")
	return txHash, nil
}

func deliverAndSwap(transactor *tx.Transactor, token common.Address, orderID [32]byte, amount *big.Int, params *SwapAndBridgeFunctionParams, fee *big.Int, feeReceiver common.Address, value *big.Int) (txHash common.Hash, err error) {
	fields := map[string]interface{}{
		"orderID":     hex.EncodeToString(orderID[:]),
		"initiator":   Initiator,
		"token":       token,
		"amount":      amount,
		"swapData":    hex.EncodeToString(params.SwapData),
		"bridgeData":  hex.EncodeToString(params.BridgeData),
		"feeData":     hex.EncodeToString(params.FeeData),
		"fee":         fee,
		"feeReceiver": feeReceiver,
		"value":       value,
	}
	txHash, err = transactor.DeliverAndSwap(orderID, Initiator, token, amount, params.SwapData, params.BridgeData, params.FeeData, fee, feeReceiver, value)
	if err != nil {
		fields["error"] = err.Error()
		log.Logger().WithFields(fields).Error("failed to send deliver and swap transaction")
		return txHash, err
	}
	fields["hash"] = txHash.Hex()
	log.Logger().WithFields(fields).Info("completed send deliver and swap transaction")
	return txHash, nil
}
