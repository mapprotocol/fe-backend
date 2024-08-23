package task

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/mapprotocol/fe-backend/third-party/filter"
	blog "log"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/mapprotocol/fe-backend/dao"
	"github.com/mapprotocol/fe-backend/params"
	"github.com/mapprotocol/fe-backend/resource/log"
	"github.com/mapprotocol/fe-backend/third-party/butter"
	"github.com/mapprotocol/fe-backend/third-party/mempool"
	"github.com/mapprotocol/fe-backend/utils"
	"github.com/mapprotocol/fe-backend/utils/alarm"
	"github.com/mapprotocol/fe-backend/utils/tx"
)

const (
	SwapType = "exactIn"
)

const (
	NetworkMainnet = "mainnet"
	NetworkTestnet = "testnet"
)

var Initiator = common.HexToAddress("0x22776") // todo

var (
	netParams    = &chaincfg.Params{}
	btcApiClient = &mempool.MempoolClient{}
)

//var (
//	senderAddress    btcutil.Address
//	senderPrivateKey *btcec.PrivateKey
//)

func InitMempoolClient(network string) {
	switch network {
	case NetworkMainnet, "":
		netParams = &chaincfg.MainNetParams
		blog.Print("initialized network: ", NetworkMainnet)
	case NetworkTestnet:
		netParams = &chaincfg.TestNet3Params
		blog.Print("initialized network: ", NetworkTestnet)
	default:
		panic("unknown network")
	}

	btcApiClient = mempool.NewClient(netParams)
	blog.Print("initialized mempool client，network: ", network)
}

// HandlePendingOrdersOfFirstStageFromBTCToEVM filter pending orders of first stage and check relay address balance.
// if balance is enough update order status to confirmed
// action=1, stage=1, status=1(TxPrepareSend) ==> stage=1, status=4(TxConfirmed)
func HandlePendingOrdersOfFirstStageFromBTCToEVM() {
	order := dao.BitcoinOrder{
		SrcChain: params.BTCChainID,
		Action:   dao.OrderActionToEVM,
		Stage:    dao.OrderStag1,
		Status:   dao.OrderStatusTxPrepareSend,
	}
	for {
		for id := uint64(1); ; {
			orders, err := order.GetOldest10ByID(id)
			if err != nil {
				log.Logger().WithField("error", err.Error()).Error("failed to get pending status order")
				alarm.Slack(context.Background(), "failed to get pending status order")
				time.Sleep(5 * time.Second)
				continue
			}

			length := len(orders)
			if length == 0 {
				log.Logger().Info("not found pending status order", "time", time.Now())
				time.Sleep(10 * time.Second)
				break
			}

			for i, o := range orders {
				if i == length-1 {
					id = o.ID + 1
				}
				relayer, err := btcutil.DecodeAddress(o.Relayer, netParams)
				if err != nil {
					log.Logger().WithField("relayer", o.Relayer).WithField("error", err.Error()).Error("failed to decode relayer")
					alarm.Slack(context.Background(), "failed to decode relayer")
					continue
				}

				utxo, err := btcApiClient.ListUnspent(relayer)
				if err != nil {
					log.Logger().WithField("relayer", relayer.String()).WithField("error", err.Error()).Error("failed to list unspent")
					alarm.Slack(context.Background(), "failed to list unspent")
					continue
				}
				if len(utxo) != 1 { // todo get first utxo ?
					log.Logger().WithField("relayer", o.Relayer).WithField("total", len(utxo)).Debug("invalid utxo")
					continue
				}

				// todo value >= in amount
				inAmount := new(big.Float).Quo(new(big.Float).SetInt64(utxo[0].Output.Value), big.NewFloat(params.BTCDecimal))
				_, afterAmount := deductFees(new(big.Float).SetInt64(utxo[0].Output.Value), FeeRate)
				afterAmountFloat := new(big.Float).Quo(afterAmount, big.NewFloat(params.BTCDecimal))
				update := &dao.BitcoinOrder{
					InAmount:    inAmount.Text('f', -1),
					InAmountSat: strconv.FormatInt(utxo[0].Output.Value, 10),
					RelayToken:  params.BTCTokenAddress,
					RelayAmount: afterAmountFloat.Text('f', 8),
					Status:      dao.OrderStatusTxConfirmed,
				}
				if err := dao.NewBitcoinOrderWithID(o.ID).Updates(update); err != nil {
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

// HandleConfirmedOrdersOfFirstStageFromBTCToEVM filter confirmed orders of first stage and send transaction to chain pool.
// set order stage to stage2 and set status to pending
// action=1, stage=1, status=4(TxConfirmed) ==> stage=2, status=1(TxPrepareSend) ==> stage=2, status=2(TxSent)
func HandleConfirmedOrdersOfFirstStageFromBTCToEVM() {
	order := dao.BitcoinOrder{
		SrcChain: params.BTCChainID,
		Action:   dao.OrderActionToEVM,
		Stage:    dao.OrderStag1,
		Status:   dao.OrderStatusTxConfirmed,
	}
	for {
		for id := uint64(1); ; {
			orders, err := order.GetOldest10ByID(id)
			if err != nil {
				// todo query parameters write to log
				log.Logger().WithField("error", err.Error()).Error("failed to get confirmed status order")
				alarm.Slack(context.Background(), "failed to get confirmed status order")
				time.Sleep(5 * time.Second)
				continue
			}

			length := len(orders)
			if length == 0 {
				log.Logger().WithField("time", time.Now()).Info("not found confirmed status order")
				time.Sleep(10 * time.Second)
				break
			}

			for i, o := range orders {
				if i == length-1 {
					id = o.ID + 1
				}

				wbtc := params.WBTCOfChainPool
				decimal := params.WBTCDecimalOfChainPool
				chainIDOfChainPool := params.ChainIDOfChainPool
				chainInfo := &dao.ChainPool{}
				if isMultiChainPool && o.SrcChain == params.ChainIDOfEthereum {
					wbtc = params.WBTCOfEthereum
					decimal = params.WBTCDecimalOfEthereum
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

				orderID := utils.Uint64ToByte32(o.ID)
				amount, ok := new(big.Rat).SetString(o.RelayAmount)
				if !ok {
					fields := map[string]interface{}{
						"orderId": o.ID,
						"amount":  o.RelayAmount,
						"error":   err,
					}
					log.Logger().WithFields(fields).Error("failed to parse string to big rat")
					alarm.Slack(context.Background(), "failed to parse string to big rat")
					continue
				}
				amount = new(big.Rat).Mul(amount, new(big.Rat).SetUint64(decimal))
				amountInt := amount.Num()
				inAmountSat, ok := new(big.Int).SetString(o.InAmountSat, 10)
				if !ok {
					fields := map[string]interface{}{
						"orderId": o.ID,
						"amount":  o.InAmountSat,
						"error":   err,
					}
					log.Logger().WithFields(fields).Error("failed to parse string to big int")
					alarm.Slack(context.Background(), "failed to parse string to big int")
					continue
				}

				fee := calcProtocolFees(inAmountSat, o.FeeRatio, decimal)
				txHash := common.Hash{}
				if o.DstChain == chainIDOfChainPool && strings.ToLower(o.DstToken) == strings.ToLower(wbtc) {
					update := &dao.BitcoinOrder{
						Stage:  dao.OrderStag2,
						Status: dao.OrderStatusTxPrepareSend,
					}
					if err := dao.NewBitcoinOrderWithID(o.ID).Updates(update); err != nil {
						log.Logger().WithField("update", utils.JSON(update)).WithField("error", err.Error()).Error("failed to update order status")
						alarm.Slack(context.Background(), "failed to update order status")
						time.Sleep(5 * time.Second)
						continue
					}

					txHash, err = deliver(transactor, common.HexToAddress(wbtc), orderID, amountInt, common.HexToAddress(o.Receiver), fee, common.HexToAddress(o.FeeCollector))
					if err != nil {
						log.Logger().WithField("error", err.Error()).Error("failed to send deliver transaction")
						alarm.Slack(context.Background(), "failed to send deliver transaction")
						time.Sleep(5 * time.Second)
						continue
					}
				} else {
					relayAmountStr := o.RelayAmount
					if fee.Cmp(big.NewInt(0)) == 1 {
						relayAmount, err := strconv.ParseFloat(o.RelayAmount, 10)
						if err != nil {
							log.Logger().WithField("amount", o.RelayAmount).WithField("error", err.Error()).Error("failed to parse relay amount")
							alarm.Slack(context.Background(), "failed to parse relay amount")
							continue
						}
						relayAmountSat, err := btcutil.NewAmount(relayAmount)
						if err != nil {
							log.Logger().WithField("amount", relayAmount).WithField("error", err.Error()).Error("failed to convert relay amount to satoshi")
							alarm.Slack(context.Background(), "failed to convert relay amount to satoshi")
							continue
						}

						fee := new(big.Int).Mul(inAmountSat, new(big.Int).SetUint64(o.FeeRatio))
						fee = new(big.Int).Div(fee, big.NewInt(10000))
						relayAmountStr = strconv.FormatFloat(float64(int64(relayAmountSat)-fee.Int64())/params.BTCDecimal, 'f', -1, 64)
					}

					request := &butter.RouterAndSwapRequest{
						FromChainID:     params.ChainIDOfChainPool,
						ToChainID:       o.DstChain,
						Amount:          relayAmountStr,
						TokenInAddress:  params.WBTCOfChainPool,
						TokenOutAddress: o.DstToken,
						Type:            SwapType,
						Slippage:        o.Slippage / 3 * 2,
						From:            sender,
						Receiver:        o.Receiver,
					}
					data, err := butter.RouteAndSwap(request)
					if err != nil {
						log.Logger().WithField("error", err.Error()).Error("failed to create router and swap request")
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

					update := &dao.BitcoinOrder{
						Stage:  dao.OrderStag2,
						Status: dao.OrderStatusTxPrepareSend,
					}
					if err := dao.NewBitcoinOrderWithID(o.ID).Updates(update); err != nil {
						log.Logger().WithField("update", utils.JSON(update)).WithField("error", err.Error()).Error("failed to update order status")
						alarm.Slack(context.Background(), "failed to update order status")
						time.Sleep(5 * time.Second)
						continue
					}

					txHash, err = deliverAndSwap(transactor, common.HexToAddress(wbtc), orderID, amountInt, decodeData, fee, common.HexToAddress(o.FeeCollector), value)
					if err != nil {
						log.Logger().WithField("error", err.Error()).Error("failed to send deliver and swap transaction")
						alarm.Slack(context.Background(), "failed to send deliver and swap transaction")
						time.Sleep(5 * time.Second)
						continue
					}
				}

				update := &dao.BitcoinOrder{
					Stage:     dao.OrderStag2,
					Status:    dao.OrderStatusTxSent,
					OutTxHash: txHash.String(),
				}
				if err := dao.NewBitcoinOrderWithID(o.ID).Updates(update); err != nil {
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

// HandlePendingOrdersOfSecondStageFromBTCToEVM filter pending orders of second stage and check transaction is confirmed.
// if transaction is confirmed, update order status to confirmed.
// action=1, stage=2, status=2(TxSent) ==> status=3(TxFailed)/4(TxConfirmed)
func HandlePendingOrdersOfSecondStageFromBTCToEVM() {
	order := dao.BitcoinOrder{
		SrcChain: params.BTCChainID,
		Action:   dao.OrderActionToEVM,
		Stage:    dao.OrderStag2,
		Status:   dao.OrderStatusTxSent,
	}
	for {
		for id := uint64(1); ; {
			orders, err := order.GetOldest10ByID(id)
			if err != nil {
				log.Logger().WithField("error", err.Error()).Error("failed to get confirmed status order")
				alarm.Slack(context.Background(), "failed to get confirmed status order")
				time.Sleep(5 * time.Second)
				continue
			}

			length := len(orders)
			if length == 0 {
				log.Logger().Info("not found confirmed status order", "time", time.Now())
				time.Sleep(10 * time.Second)
				break
			}

			for i, o := range orders {
				if i == length-1 {
					id = o.ID + 1
				}

				chainInfo, err := dao.NewChainPoolWithChainID(params.ChainIDOfChainPool).First()
				if err != nil {
					log.Logger().WithField("chainID", o.DstChain).WithField("error", err.Error()).Error("failed to get chain info")
					alarm.Slack(context.Background(), "failed to get chain info")
					time.Sleep(5 * time.Second)
					continue
				}

				// todo call NewCaller
				transactor, err := tx.NewTransactor(chainInfo.ChainRPC, chainInfo.FeRouterContract, 0)
				if err != nil {
					log.Logger().WithField("endpoint", chainInfo.ChainRPC).WithField("error", err.Error()).Error("failed to create transactor")
					alarm.Slack(context.Background(), "failed to create transactor")
					time.Sleep(5 * time.Second)
					continue
				}
				pending, err := transactor.TransactionIsPending(common.HexToHash(o.OutTxHash))
				if err != nil {
					log.Logger().WithField("endpoint", chainInfo.ChainRPC).WithField("error", err.Error()).Error("failed to create transactor")
					alarm.Slack(context.Background(), "failed to create transactor")
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

				update := &dao.BitcoinOrder{
					Status: dao.OrderStatusTxConfirmed,
				}
				if status == types.ReceiptStatusFailed {
					update.Status = dao.OrderStatusTxFailed
				}
				if err := dao.NewBitcoinOrderWithID(o.ID).Updates(update); err != nil {
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

// HandlePendingOrdersOfFirstStageFromEVM filter the OnReceived event log.
// If this log is found, update order status to confirmed
// create order from evm ( action=2, stage=1, status=4(TxConfirmed))
// todo multi chain pool
func HandlePendingOrdersOfFirstStageFromEVM() {
	chainID := params.ChainIDOfChainPool
	topic := params.OnReceivedTopic
	filterLog := dao.NewFilterLog(chainID, topic)
	for {
		gotLog, err := filterLog.First()
		if err != nil {
			fields := map[string]interface{}{
				"chainID": params.ChainIDOfChainPool,
				"topic":   params.OnReceivedTopic,
				"error":   err.Error(),
			}
			log.Logger().WithFields(fields).Error("failed to get filter log info")
			alarm.Slack(context.Background(), "failed to get filter log info")
			time.Sleep(5 * time.Second)
			continue
		}

		logs, err := filter.GetLogs(gotLog.LatestLogID, params.ChainIDOfChainPool, params.OnReceivedTopic, uint8(20))
		if err != nil {
			fields := map[string]interface{}{
				"id":      gotLog.LatestLogID,
				"chainID": params.ChainIDOfChainPool,
				"topic":   params.OnReceivedTopic,
				"limit":   uint8(20),
				"error":   err.Error(),
			}
			log.Logger().WithFields(fields).Error("failed to get logs")
			alarm.Slack(context.Background(), "failed to get logs")
			continue
		}
		if len(logs) == 0 {
			log.Logger().WithField("id", gotLog.LatestLogID).WithField("time", time.Now()).Info("not found on received logs")
			time.Sleep(5 * time.Second)
			continue
		}

		for _, lg := range logs {
			if lg.Id <= gotLog.LatestLogID {
				continue
			}
			// 1. parse log data
			logData, err := hex.DecodeString(lg.LogData)
			if err != nil {
				fields := map[string]interface{}{
					"id":      lg.Id,
					"chainID": params.ChainIDOfChainPool,
					"topic":   params.OnReceivedTopic,
					"logData": lg.LogData,
					"error":   err.Error(),
				}
				log.Logger().WithFields(fields).Error("failed to decode log data")
				alarm.Slack(context.Background(), "failed to decode log data")
				UpdateLogID(chainID, topic, lg.Id)
				continue
			}

			onReceived, err := UnpackOnReceived(logData)
			if err != nil {
				fields := map[string]interface{}{
					"id":      lg.Id,
					"chainID": params.ChainIDOfChainPool,
					"topic":   params.OnReceivedTopic,
					"logData": lg.LogData,
					"error":   err.Error(),
				}
				log.Logger().WithFields(fields).Error("failed to unpack log data")
				alarm.Slack(context.Background(), "failed to unpack log data")
				UpdateLogID(chainID, topic, lg.Id)
				continue
			}

			_, afterAmount := deductFees(new(big.Float).SetInt(onReceived.ChainPoolTokenAmount), FeeRate)
			if onReceived.DstChain.String() == params.TONChainID {
				afterAmountFloat := new(big.Float).Quo(afterAmount, big.NewFloat(params.USDTDecimalOfChainPool))
				order := &dao.Order{
					OrderIDFromContract: onReceived.BridgeId,
					SrcChain:            onReceived.SrcChain.String(),
					SrcToken:            string(onReceived.SrcToken),
					Sender:              string(onReceived.Sender),
					InAmount:            onReceived.InAmount,
					RelayToken:          params.USDTOfChainPool,
					RelayAmount:         afterAmountFloat.String(),
					DstChain:            onReceived.DstChain.String(),
					DstToken:            string(onReceived.DstToken),
					Receiver:            string(onReceived.Receiver),
					Action:              dao.OrderActionFromEVM,
					Stage:               dao.OrderStag1,
					Status:              dao.OrderStatusTxConfirmed,
					Slippage:            onReceived.Slippage,
				}
				if err := order.Create(); err != nil {
					log.Logger().WithField("order", utils.JSON(order)).WithField("error", err).Error("failed to create order")
					alarm.Slack(context.Background(), "failed to update order status")
					UpdateLogID(chainID, topic, lg.Id)
					continue
				}
			} else if onReceived.DstChain.String() == params.BTCChainID {
				inAmount, err := strconv.ParseFloat(onReceived.InAmount, 10)
				if err != nil {
					log.Logger().WithField("amount", onReceived.InAmount).WithField("error", err.Error()).Error("failed to parse in amount")
					alarm.Slack(context.Background(), "failed to parse in amount")
					UpdateLogID(chainID, topic, lg.Id)
					continue
				}
				inAmountSat, err := btcutil.NewAmount(inAmount)
				if err != nil {
					log.Logger().WithField("amount", onReceived.InAmount).WithField("error", err.Error()).Error("failed to convert in amount to satoshi")
					alarm.Slack(context.Background(), "failed to convert in amount to satoshi")
					UpdateLogID(chainID, topic, lg.Id)
					continue
				}

				afterAmountFloat := new(big.Float).Quo(afterAmount, new(big.Float).SetUint64(params.WBTCDecimalOfChainPool))
				order := &dao.BitcoinOrder{
					SrcChain:    onReceived.SrcChain.String(),
					SrcToken:    string(onReceived.SrcToken),
					Sender:      string(onReceived.Sender),
					InAmount:    onReceived.InAmount,
					InAmountSat: strconv.FormatInt(int64(inAmountSat), 10),
					RelayToken:  params.WBTCOfChainPool,
					RelayAmount: afterAmountFloat.String(),
					DstChain:    onReceived.DstChain.String(),
					DstToken:    string(onReceived.DstToken),
					Receiver:    string(onReceived.Receiver),
					Action:      dao.OrderActionFromEVM,
					Stage:       dao.OrderStag1,
					Status:      dao.OrderStatusTxConfirmed,
					Slippage:    onReceived.Slippage,
				}

				if err := order.Create(); err != nil {
					log.Logger().WithField("order", utils.JSON(order)).WithField("error", err).Error("failed to create order")
					alarm.Slack(context.Background(), "failed to update order status")
					UpdateLogID(chainID, topic, lg.Id)
					continue
				}
			} else {
				log.Logger().WithField("dstChain", onReceived.DstChain.String()).Error("unsupported dst chain")
				alarm.Slack(context.Background(), fmt.Sprintf("unsupported dst chain: %s", onReceived.DstChain.String()))
				UpdateLogID(chainID, topic, lg.Id)
				continue
			}

			UpdateLogID(params.ChainIDOfChainPool, params.OnReceivedTopic, lg.Id)
		}

		time.Sleep(10 * time.Second)
	}
}

//func HandleConfirmedOrdersOfFirstStageFromEVM() {
//	order := dao.NewOrder()
//	for {
//		for id := uint64(1); ; {
//			orders, err := order.GetOldest10ByStatus(id, dao.OrderActionFromEVM, dao.OrderStag1, dao.OrderStatusConfirmed)
//			if err != nil {
//				// todo query parameters write to log
//				log.Logger().WithField("error", err.Error()).Error("failed to get confirmed status order")
//				alarm.Slack(context.Background(), "failed to get confirmed status order")
//				time.Sleep(5 * time.Second)
//				continue
//			}
//
//			length := len(orders)
//			if length == 0 {
//				log.Logger().Info("not found confirmed status order", "time", time.Now())
//				time.Sleep(10 * time.Second)
//				break
//			}
//
//			for i, o := range orders {
//				if i == length-1 {
//					id = o.ID + 1
//				}
//				// send transaction to bitcoin
//				fees, err := btcApiClient.RecommendedFees()
//				if err != nil {
//					log.Logger().WithField("error", err).Error("failed to get fee rate from mempool")
//					alarm.Slack(context.Background(), "failed to get fee rate from mempool")
//					time.Sleep(10 * time.Second)
//					continue
//				}
//
//				amount, err := strconv.ParseInt(o.InAmount, 10, 64)
//				if err != nil {
//					log.Logger().WithField("amount", o.InAmount).WithField("error", err.Error()).Error("failed to parse amount")
//					alarm.Slack(context.Background(), "failed to parse amount")
//					continue
//				}
//				receiver, err := btcutil.DecodeAddress(o.Receiver, netParams)
//				if err != nil {
//					log.Logger().WithField("address", o.Receiver).WithField("error", err.Error()).Error("failed to decode receiver address")
//					alarm.Slack(context.Background(), "failed to decode receiver address")
//					continue
//				}
//				hash, err := SendTransaction(btcApiClient, fees.FastestFee, amount, senderAddress, receiver, senderPrivateKey)
//				if err != nil {
//					log.Logger().WithField("error", err.Error()).Error("failed to send transaction to bitcoin")
//					alarm.Slack(context.Background(), "failed to send transaction to bitcoin")
//					time.Sleep(10 * time.Second)
//					continue
//				}
//
//				update := &dao.BitcoinOrder{
//					ID:        o.ID,
//					Stage:     dao.OrderStag2,
//					Status:    dao.OrderStatusPending,
//					OutTxHash: hash.String(),
//				}
//				if err := dao.NewOrderWithID(o.ID).Updates(update); err != nil {
//					log.Logger().WithField("update", utils.JSON(update)).WithField("error", err.Error()).Error("failed to update order status")
//					alarm.Slack(context.Background(), "failed to update order status")
//					time.Sleep(5 * time.Second)
//					continue
//				}
//
//			}
//			time.Sleep(10 * time.Second)
//		}
//	}
//}
//
//func HandlePendingOrdersOfSecondStageFromEVM() {
//	order := dao.NewOrder()
//	for {
//		for id := uint64(1); ; {
//			orders, err := order.GetOldest10ByStatus(id, dao.OrderActionFromEVM, dao.OrderStag2, dao.OrderStatusPending)
//			if err != nil {
//				log.Logger().WithField("error", err.Error()).Error("failed to get confirmed status order")
//				alarm.Slack(context.Background(), "failed to get confirmed status order")
//				time.Sleep(5 * time.Second)
//				continue
//			}
//
//			length := len(orders)
//			if length == 0 {
//				log.Logger().Info("not found confirmed status order", "time", time.Now())
//				time.Sleep(10 * time.Second)
//				break
//			}
//
//			for i, o := range orders {
//				if i == length-1 {
//					id = o.ID + 1
//				}
//
//				hash, err := chainhash.NewHashFromStr(o.OutTxHash)
//				if err != nil {
//					log.Logger().WithField("txHash", o.InTxHash).WithField("error", err.Error()).Error("failed to creates a hash from a hash string")
//					alarm.Slack(context.Background(), "failed to creates a hash from a hash string")
//					continue
//				}
//				ret, err := btcApiClient.TransactionStatus(hash)
//				if err != nil {
//					log.Logger().WithField("error", err.Error()).Error("failed to get transaction status from bitcoin")
//					alarm.Slack(context.Background(), "failed to get transaction status from bitcoin")
//					continue
//				}
//				// todo tx failed
//				if !ret.Confirmed {
//					continue
//				}
//				// todo get block number to judge whether the transaction is confirmed
//
//				update := &dao.BitcoinOrder{
//					Status: dao.OrderStatusConfirmed,
//				}
//				if err := dao.NewOrderWithID(o.ID).Updates(update); err != nil {
//					log.Logger().WithField("update", utils.JSON(update)).WithField("error", err.Error()).Error("failed to update order status")
//					alarm.Slack(context.Background(), "failed to update order status")
//					time.Sleep(5 * time.Second)
//					continue
//				}
//			}
//			time.Sleep(10 * time.Second)
//		}
//	}
//}

func TransactionIsConfirmed(txHash string) (bool, error) {
	hash, err := chainhash.NewHashFromStr(txHash)
	if err != nil {
		log.Logger().WithField("txHash", txHash).WithField("error", err.Error()).Error("failed to creates a hash from a hash string")
		return false, err
	}

	ret, err := btcApiClient.TransactionStatus(hash)
	if err != nil {
		log.Logger().WithField("error", err.Error()).Error("failed to get transaction status")
		return false, err
	}

	return ret.Confirmed, nil
}

//func fees(amount, feeRate *big.Int) (feeAmount, afterAmount *big.Int) {
//	feeRate = new(big.Int).Div(feeRate, big.NewInt(1000))
//	feeAmount = new(big.Int).Mul(amount, feeRate)
//	afterAmount = new(big.Int).Sub(amount, feeAmount)
//	return feeAmount, afterAmount
//}

func UpdateLogID(chainID, topic string, logID uint64) {
	if err := dao.NewFilterLog(chainID, topic).UpdateLatestLogID(logID); err != nil {
		fields := map[string]interface{}{
			"chainID":     chainID,
			"topic":       topic,
			"latestLogID": logID,
			"error":       err.Error(),
		}
		log.Logger().WithFields(fields).Error("failed to update filter log")
		alarm.Slack(context.Background(), "failed to update filter log")
	}
}

func deductFees(amount, feeRate *big.Float) (feeAmount, afterAmount *big.Float) {
	log.Logger().WithField("amount", amount).WithField("feeRate", feeRate).Info("before deduction of fees")
	feeRate = new(big.Float).Quo(feeRate, big.NewFloat(10000))
	feeAmount = new(big.Float).Mul(amount, feeRate)
	afterAmount = new(big.Float).Sub(amount, feeAmount)
	log.Logger().WithField("amount", amount).WithField("feeAmount", feeAmount).Info("after deduction of fees")
	return feeAmount, afterAmount
}

func calcProtocolFees(inAmountSat *big.Int, feeRatio uint64, wbtcDecimal uint64) *big.Int {
	fee := new(big.Int).Mul(inAmountSat, new(big.Int).SetUint64(feeRatio))
	fee = new(big.Int).Div(fee, big.NewInt(10000))
	if wbtcDecimal != params.BTCDecimal {
		fee = new(big.Int).Mul(fee, new(big.Int).SetUint64(wbtcDecimal))
		fee = new(big.Int).Div(fee, new(big.Int).SetUint64(params.BTCDecimal))
	}
	return fee
}
