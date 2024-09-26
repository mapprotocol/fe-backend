package task

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/mapprotocol/fe-backend/third-party/filter"
	"github.com/shopspring/decimal"
	blog "log"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
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

const (
	FeeRateMultiple = 2
	FeeRateLimit    = 200
)

const (
	BridgeFeeRate            int64 = 30
	BridgeFeeRateDenominator int64 = 10000
)

const BaseTxFeeMultiplier = 1.5

var Initiator = common.HexToAddress("0x22776") // todo

var (
	netParams    = &chaincfg.Params{}
	btcApiClient = &mempool.MempoolClient{}
)

var globalFeeRate int64 = 20

var ToTONBaseTxFee = new(big.Int).SetUint64(uint64(params.USDTDecimalOfTON * 1.5)) // 1.5 USDT
var TONToEVMBaseTxFee = new(big.Int).SetUint64(params.USDTDecimalOfChainPool)      // 1 USDT
var BitcoinToEVMBaseTxFee = new(big.Int).SetUint64(700 * BaseTxFeeMultiplier)      // 0.0000105 WBTC
var BitcoinTxBytes = new(big.Int).SetUint64(200)

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
// if balance is enough update bitcoin order status to confirmed
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
				if len(utxo) != 1 {
					log.Logger().WithField("relayer", o.Relayer).WithField("total", len(utxo)).Debug("invalid utxo")
					continue
				}

				value := utxo[0].Output.Value
				if uint64(value) < o.InAmountSat {
					log.Logger().WithField("relayer", o.Relayer).WithField("value", value).WithField("inAmountSat", o.InAmountSat).Debug("get utxo value is less than in amount sat")
					alarm.Slack(context.Background(), "get utxo value is less than in amount sat")

					update := &dao.BitcoinOrder{
						Status: dao.OrderStatusTxFailed,
					}
					if err := dao.NewBitcoinOrderWithID(o.ID).Updates(update); err != nil {
						log.Logger().WithField("update", utils.JSON(update)).WithField("error", err.Error()).Error("failed to update bitcoin order status")
						alarm.Slack(context.Background(), "failed to update bitcoin order status")
						time.Sleep(5 * time.Second)
					}
					continue
				}

				//inAmount := new(big.Float).Quo(new(big.Float).SetInt64(value), big.NewFloat(params.BTCDecimal))
				inAmount := decimal.NewFromInt(value).Div(decimal.NewFromFloat(params.BTCDecimal))

				bridgeFees, afterAmount := deductBitcoinToEVMBridgeFees(new(big.Int).SetInt64(value), big.NewInt(BridgeFeeRate))
				//afterAmountFloat := new(big.Float).Quo(new(big.Float).SetInt(afterAmount), big.NewFloat(params.BTCDecimal))
				//afterAmountFloat := decimal.NewFromBigInt(afterAmount, 0).Div(decimal.NewFromFloat(params.BTCDecimal))
				afterAmountWithFixedDecimal := convertDecimal(afterAmount, params.BTCDecimalNumber, params.FixedDecimalNumber)

				update := &dao.BitcoinOrder{
					//InAmount:    inAmount.Text('f', params.BTCDecimalNumber),
					InAmount:       inAmount.StringFixedBank(params.BTCDecimalNumber),
					InAmountSat:    uint64(utxo[0].Output.Value),
					InTxHash:       utxo[0].Outpoint.Hash.String(),
					BridgeFee:      bridgeFees.Uint64(),
					RelayToken:     params.WBTCOfChainPool,
					RelayAmountInt: afterAmountWithFixedDecimal.Uint64(),
					Status:         dao.OrderStatusTxConfirmed,
				}
				if err := dao.NewBitcoinOrderWithID(o.ID).Updates(update); err != nil {
					log.Logger().WithField("update", utils.JSON(update)).WithField("error", err.Error()).Error("failed to update bitcoin order status")
					alarm.Slack(context.Background(), "failed to update bitcoin order status")
					time.Sleep(5 * time.Second)
					continue
				}

			}
			time.Sleep(30 * time.Second)
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
				wbtcDecimal := params.WBTCDecimalOfChainPool
				wbtcDecimalNumber := params.WBTCDecimalNumberOfChainPool
				chainIDOfChainPool := params.ChainIDOfChainPool
				chainInfo := &dao.ChainPool{}
				if isMultiChainPool && o.SrcChain == params.ChainIDOfEthereum {
					wbtc = params.WBTCOfEthereum
					wbtcDecimal = params.WBTCDecimalOfEthereum
					wbtcDecimalNumber = params.WBTCDecimalNumberOfEthereum
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

				//amount, ok := new(big.Rat).SetString(o.RelayAmountInt)
				//if !ok {
				//	fields := map[string]interface{}{
				//		"orderId": o.ID,
				//		"amount":  o.RelayAmountInt,
				//		"error":   err,
				//	}
				//	log.Logger().WithFields(fields).Error("failed to parse string to big rat")
				//	alarm.Slack(context.Background(), "failed to parse string to big rat")
				//	continue
				//}
				//amount = new(big.Rat).Mul(amount, new(big.Rat).SetUint64(decimal))
				//amountInt := amount.Num()

				//relayAmount, err := decimal.NewFromString(o.RelayAmountInt)
				//if err != nil {
				//	fields := map[string]interface{}{
				//		"orderId": o.ID,
				//		"amount":  o.RelayAmountInt,
				//		"error":   err,
				//	}
				//	log.Logger().WithFields(fields).Error("failed to parse string to decimal")
				//	continue
				//}
				//
				//relayAmount = relayAmount.Mul(decimal.NewFromUint64(wbtcDecimal))
				//amountInt := relayAmount.BigInt()

				//inAmountSat, ok := new(big.Int).SetString(o.InAmountSat, 10)
				//if !ok {
				//	fields := map[string]interface{}{
				//		"orderId": o.ID,
				//		"amount":  o.InAmountSat,
				//		"error":   err,
				//	}
				//	log.Logger().WithFields(fields).Error("failed to parse string to big int")
				//	alarm.Slack(context.Background(), "failed to parse string to big int")
				//	continue
				//}

				txHash := common.Hash{}
				value := big.NewInt(0)
				fee := calcProtocolFees(new(big.Int).SetUint64(o.InAmountSat), o.FeeRatio, wbtcDecimal)
				fee = convertDecimal(fee, uint64(wbtcDecimalNumber), params.FixedDecimalNumber)
				srcChain, ok := new(big.Int).SetString(o.SrcChain, 10)
				if !ok {
					log.Logger().WithField("srcChain", o.SrcChain).Error("failed to parse src chain to big int")
					alarm.Slack(context.Background(), "failed to parse src chain to big int")
					continue
				}

				dstChain, ok := new(big.Int).SetString(o.DstChain, 10)
				if !ok {
					log.Logger().WithField("dstChain", o.DstChain).Error("failed to parse dst chain to big int")
					alarm.Slack(context.Background(), "failed to parse dst chain to big int")
					continue
				}

				deliverParam := &tx.DeliverParam{
					OrderId:     utils.Uint64ToByte32(o.ID),
					Receiver:    common.HexToAddress(o.Receiver),
					Token:       common.HexToAddress(wbtc),
					Amount:      convertDecimal(new(big.Int).SetUint64(o.RelayAmountInt), params.FixedDecimalNumber, uint64(wbtcDecimalNumber)),
					FromChain:   srcChain,
					ToChain:     dstChain,
					Fee:         fee,
					FeeReceiver: common.HexToAddress(o.FeeCollector),
					From:        []byte(o.Sender),
					//ButterData:  []byte{},
				}
				//if o.DstChain == chainIDOfChainPool && strings.ToLower(o.DstToken) == strings.ToLower(wbtc) {
				//	update := &dao.BitcoinOrder{
				//		Stage:  dao.OrderStag2,
				//		Status: dao.OrderStatusTxPrepareSend,
				//	}
				//	if err := dao.NewBitcoinOrderWithID(o.ID).Updates(update); err != nil {
				//		log.Logger().WithField("update", utils.JSON(update)).WithField("error", err.Error()).Error("failed to update bitcoin order status")
				//		alarm.Slack(context.Background(), "failed to update bitcoin order status")
				//		time.Sleep(5 * time.Second)
				//		continue
				//	}
				//
				//	//amountBigInt := convertDecimal(new(big.Int).SetUint64(o.RelayAmountInt), params.FixedDecimalNumber, uint64(wbtcDecimalNumber))
				//	txHash, err = deliver(transactor, common.HexToAddress(wbtc), orderID, amountBigInt, common.HexToAddress(o.Receiver), fee, common.HexToAddress(o.FeeCollector))
				//	if err != nil {
				//		log.Logger().WithField("error", err.Error()).Error("failed to send deliver transaction")
				//		alarm.Slack(context.Background(), "failed to send deliver transaction")
				//		time.Sleep(5 * time.Second)
				//		continue
				//	}
				//}

				if o.DstChain != chainIDOfChainPool || strings.ToLower(o.DstToken) != strings.ToLower(wbtc) {
					//relayAmountStr := strconv.FormatFloat(float64(o.RelayAmountInt-fee.Uint64())/params.FixedDecimal, 'f', params.WBTCDecimalNumberOfChainPool, 64)
					//relayAmount := decimal.NewFromUint64(o.RelayAmountInt).Sub(decimal.NewFromUint64(fee.Uint64())) // todo big int calc
					//relayAmountStr := unwrapFixedDecimal(relayAmount).StringFixedBank(params.WBTCDecimalNumberOfChainPool)
					//relayAmount, err := decimal.NewFromString(relayAmountStr)

					relayAmount := new(big.Int).Sub(new(big.Int).SetUint64(o.RelayAmountInt), fee)
					relayAmountStr := unwrapFixedDecimal(decimal.NewFromBigInt(relayAmount, 0)).String()

					//relayAmount, err := decimal.NewFromString(relayAmountStr)
					//if err != nil {
					//	fields := map[string]interface{}{
					//		"orderId": o.ID,
					//		"amount":  o.RelayAmountInt,
					//		"error":   err,
					//	}
					//	log.Logger().WithFields(fields).Error("failed to parse string to float")
					//	continue
					//}
					//
					//relayAmountBigInt := relayAmount.Add(decimal.NewFromBigInt(fee, 0)).Mul(decimal.NewFromUint64(wbtcDecimal)).BigInt()

					//if fee.Cmp(big.NewInt(0)) == 1 {
					//	//relayAmount, err := strconv.ParseFloat(o.RelayAmountInt, 10)
					//	//if err != nil {
					//	//	log.Logger().WithField("amount", o.RelayAmountInt).WithField("error", err.Error()).Error("failed to parse relay amount")
					//	//	alarm.Slack(context.Background(), "failed to parse relay amount")
					//	//	continue
					//	//}
					//	//relayAmountSat, err := btcutil.NewAmount(relayAmount)
					//	//if err != nil {
					//	//	log.Logger().WithField("amount", relayAmount).WithField("error", err.Error()).Error("failed to convert relay amount to satoshi")
					//	//	alarm.Slack(context.Background(), "failed to convert relay amount to satoshi")
					//	//	continue
					//	//}
					//
					//	fee := new(big.Int).Mul(inAmountSat, new(big.Int).SetUint64(o.FeeRatio))
					//	fee = new(big.Int).Div(fee, big.NewInt(10000))
					//	relayAmountStr = strconv.FormatFloat(float64(int64(relayAmountSat)-fee.Int64())/params.BTCDecimal, 'f', -1, 64)
					//}
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
					butterData, err := EncodeButterData(Initiator, common.HexToAddress(o.DstToken), decodeData.SwapData, decodeData.BridgeData, decodeData.FeeData)
					if err != nil {
						log.Logger().WithField("error", err.Error()).Error("failed to encode butter data")
						alarm.Slack(context.Background(), "failed to encode butter data")
						continue
					}

					v, ok := new(big.Int).SetString(utils.TrimHexPrefix(data.Value), 16)
					if !ok {
						log.Logger().WithField("value", utils.TrimHexPrefix(data.Value)).Error("failed to parse string to big int")
						alarm.Slack(context.Background(), "failed to parse string to big int")
						continue
					}
					value = v
					deliverParam.ButterData = butterData

					//update := &dao.BitcoinOrder{
					//	Stage:  dao.OrderStag2,
					//	Status: dao.OrderStatusTxPrepareSend,
					//}
					//if err := dao.NewBitcoinOrderWithID(o.ID).Updates(update); err != nil {
					//	log.Logger().WithField("update", utils.JSON(update)).WithField("error", err.Error()).Error("failed to update bitcoin order status")
					//	alarm.Slack(context.Background(), "failed to update bitcoin order status")
					//	time.Sleep(5 * time.Second)
					//	continue
					//}
					//
					//txHash, err = deliverAndSwap(transactor, common.HexToAddress(wbtc), orderID, amountBigInt, decodeData, fee, common.HexToAddress(o.FeeCollector), value)
					//if err != nil {
					//	log.Logger().WithField("error", err.Error()).Error("failed to send deliver and swap transaction")
					//	alarm.Slack(context.Background(), "failed to send deliver and swap transaction")
					//	time.Sleep(5 * time.Second)
					//	continue
					//}
				}
				update := &dao.BitcoinOrder{
					Stage:  dao.OrderStag2,
					Status: dao.OrderStatusTxPrepareSend,
				}
				if err := dao.NewBitcoinOrderWithID(o.ID).Updates(update); err != nil {
					log.Logger().WithField("update", utils.JSON(update)).WithField("error", err.Error()).Error("failed to update bitcoin order status")
					alarm.Slack(context.Background(), "failed to update bitcoin order status")
					time.Sleep(5 * time.Second)
					continue
				}

				txHash, err = deliverAndSwap(transactor, deliverParam, value)
				if err != nil {
					log.Logger().WithField("error", err.Error()).Error("failed to send deliver and swap transaction")
					alarm.Slack(context.Background(), "failed to send deliver and swap transaction")
					time.Sleep(5 * time.Second)
					continue
				}

				update = &dao.BitcoinOrder{
					Stage:     dao.OrderStag2,
					Status:    dao.OrderStatusTxSent,
					OutTxHash: txHash.String(),
				}
				if err := dao.NewBitcoinOrderWithID(o.ID).Updates(update); err != nil {
					log.Logger().WithField("update", utils.JSON(update)).WithField("error", err.Error()).Error("failed to update bitcoin order status")
					alarm.Slack(context.Background(), "failed to update bitcoin order status")
					time.Sleep(5 * time.Second)
					continue
				}

			}
			time.Sleep(10 * time.Second)
		}
	}
}

// HandlePendingOrdersOfSecondStageFromBTCToEVM filter pending orders of second stage and check transaction is confirmed.
// if transaction is confirmed, update bitcoin order status to confirmed.
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
					log.Logger().WithField("endpoint", chainInfo.ChainRPC).WithField("error", err.Error()).Error("failed to judge transaction is pending")
					alarm.Slack(context.Background(), "failed to judge transaction is pending")
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
					log.Logger().WithField("update", utils.JSON(update)).WithField("error", err.Error()).Error("failed to update bitcoin order status")
					alarm.Slack(context.Background(), "failed to update bitcoin order status")
					time.Sleep(5 * time.Second)
					continue
				}
			}
			time.Sleep(10 * time.Second)
		}
	}
}

// HandlePendingOrdersOfFirstStageFromEVM filter the OnReceived event log.
// If this log is found, update bitcoin order status to confirmed
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

			if onReceived.DstChain.String() == params.TONChainID {
				bridgeFees, afterAmount := deductToTONBridgeFees(onReceived.ChainPoolTokenAmount, big.NewInt(BridgeFeeRate))
				//afterAmountFloat := new(big.Float).Quo(new(big.Float).SetInt(afterAmount), big.NewFloat(float64(params.USDTDecimalOfChainPool)))

				//afterAmountFloat := decimal.NewFromBigInt(afterAmount, 0).Div(decimal.NewFromUint64(params.USDTDecimalOfChainPool)) // todo * 100
				//afterAmount := new(big.Int).Div(afterAmount, new(big.Int).SetUint64(params.USDTDecimalOfChainPool))

				afterAmountWithFixedDecimal := convertDecimal(afterAmount, params.USDTDecimalNumberOfChainPool, params.FixedDecimalNumber)

				order := &dao.Order{
					OrderIDFromContract: onReceived.BridgeId,
					SrcChain:            onReceived.SrcChain.String(),
					SrcToken:            string(onReceived.SrcToken),
					Sender:              string(onReceived.Sender),
					InAmount:            onReceived.InAmount,
					BridgeFee:           bridgeFees.Uint64(),
					RelayToken:          params.USDTOfChainPool,
					//RelayAmountInt:         wrapFixedDecimal(afterAmountFloat),
					RelayAmountInt: afterAmountWithFixedDecimal.Uint64(),
					DstChain:       onReceived.DstChain.String(),
					DstToken:       string(onReceived.DstToken),
					Receiver:       string(onReceived.Receiver),
					Action:         dao.OrderActionFromEVM,
					Stage:          dao.OrderStag1,
					Status:         dao.OrderStatusTxConfirmed,
					Slippage:       onReceived.Slippage,
				}
				if err := order.Create(); err != nil {
					log.Logger().WithField("order", utils.JSON(order)).WithField("error", err).Error("failed to create order")
					alarm.Slack(context.Background(), "failed to update bitcoin order status")
					UpdateLogID(chainID, topic, lg.Id)
					continue
				}
			} else if onReceived.DstChain.String() == params.BTCChainID {
				bridgeFees, afterAmount := deductToBitcoinBridgeFees(onReceived.ChainPoolTokenAmount, big.NewInt(BridgeFeeRate))
				//afterAmountFloat := new(big.Float).Quo(new(big.Float).SetInt(afterAmount), big.NewFloat(params.WBTCDecimalOfChainPool))

				//afterAmountFloat := decimal.NewFromBigInt(afterAmount, 0).Div(decimal.NewFromUint64(params.WBTCDecimalOfChainPool))
				//afterAmount := new(big.Int).Div(afterAmount, new(big.Int).SetUint64(params.WBTCDecimalOfChainPool))

				afterAmountWithFixedDecimal := convertDecimal(afterAmount, params.WBTCDecimalNumberOfChainPool, params.FixedDecimalNumber)

				order := &dao.BitcoinOrder{
					SrcChain:   onReceived.SrcChain.String(),
					SrcToken:   string(onReceived.SrcToken),
					Sender:     string(onReceived.Sender),
					InAmount:   onReceived.InAmount,
					BridgeFee:  bridgeFees.Uint64(),
					RelayToken: params.WBTCOfChainPool,
					//RelayAmountInt: wrapFixedDecimal(afterAmountFloat), // todo multi chain pool
					RelayAmountInt: afterAmountWithFixedDecimal.Uint64(),
					DstChain:       onReceived.DstChain.String(),
					DstToken:       string(onReceived.DstToken),
					Receiver:       string(onReceived.Receiver),
					Action:         dao.OrderActionFromEVM,
					Stage:          dao.OrderStag1,
					Status:         dao.OrderStatusTxConfirmed,
					Slippage:       onReceived.Slippage,
				}

				if err := order.Create(); err != nil {
					log.Logger().WithField("order", utils.JSON(order)).WithField("error", err).Error("failed to create order")
					alarm.Slack(context.Background(), "failed to update bitcoin order status")
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

func GetFeeRate() {
	logger := log.Logger().WithField("func", "GetFeeRate")
	for {
		fees, err := btcApiClient.RecommendedFees()
		if err != nil {
			logger.WithField("error", err).Error("failed to get fee rate")
			time.Sleep(10 * time.Second)
			continue
		}

		feeRate := fees.FastestFee * FeeRateMultiple
		if feeRate > FeeRateLimit {
			feeRate = FeeRateLimit
		}

		setGlobalFeeRate(feeRate)

		logger.WithField("fastestFee", fees.FastestFee).WithField("feeRate", globalFeeRate).Info("got fee rate")
		time.Sleep(30 * time.Minute)
		continue
	}
}

func setGlobalFeeRate(feeRate int64) { // todo use CAS
	globalFeeRate = feeRate
}

func GetGlobalFeeRate() int64 { // todo use CAS
	return globalFeeRate
}

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

func wrapFixedDecimal(amount decimal.Decimal) uint64 {
	//exp := new(big.Int).Exp(big.NewInt(10), big.NewInt(params.FixedDecimalNumber), nil)
	//amt, _ := new(big.Float).Mul(amount, new(big.Float).SetInt(exp)).Int(nil)
	exp := decimal.NewFromInt(10).Pow(decimal.NewFromUint64(params.FixedDecimalNumber))
	return amount.Mul(exp).BigInt().Uint64()
}

func unwrapFixedDecimal(amount decimal.Decimal) decimal.Decimal {
	exp := decimal.NewFromInt(10).Pow(decimal.NewFromUint64(params.FixedDecimalNumber))
	return amount.Div(exp)
}

//func deductFees(amount, feeRate *big.Float) (feeAmount, afterAmount *big.Float) {
//	log.Logger().WithField("amount", amount).WithField("feeRate", feeRate).Info("before deduction of fees")
//	feeRate = new(big.Float).Quo(feeRate, big.NewFloat(10000))
//	feeAmount = new(big.Float).Mul(amount, feeRate)
//	afterAmount = new(big.Float).Sub(amount, feeAmount)
//	log.Logger().WithField("amount", amount).WithField("feeAmount", feeAmount).Info("after deduction of fees")
//	return feeAmount, afterAmount
//}

func deductFees(amount, feeRate *big.Int) (feeAmount, afterAmount *big.Int) {
	feeAmount = new(big.Int).Mul(amount, feeRate)
	feeAmount = new(big.Int).Div(feeAmount, big.NewInt(10000))
	afterAmount = new(big.Int).Sub(amount, feeAmount)

	fields := map[string]interface{}{
		"amount":      amount,
		"afterAmount": afterAmount,
		"feeAmount":   feeAmount,
	}
	log.Logger().WithFields(fields).Info("completed the deduction bridge fees")
	return feeAmount, afterAmount
}

func deductToTONBridgeFees(amount, bridgeFeeRate *big.Int) (bridgeFees, afterAmount *big.Int) {
	bridgeFees = new(big.Int).Mul(amount, bridgeFeeRate)
	bridgeFees = new(big.Int).Div(bridgeFees, big.NewInt(BridgeFeeRateDenominator))
	bridgeFees = new(big.Int).Add(bridgeFees, ToTONBaseTxFee)

	afterAmount = new(big.Int).Sub(amount, bridgeFees)

	fields := map[string]interface{}{
		"amount":        amount,
		"afterAmount":   afterAmount,
		"bridgeFees":    bridgeFees,
		"bridgeFeeRate": bridgeFeeRate,
		"baseTxFee":     ToTONBaseTxFee,
	}
	log.Logger().WithFields(fields).Info("completed the deduction to ton bridge fees")
	return bridgeFees, afterAmount
}

func deductToBitcoinBridgeFees(amount, bridgeFeeRate *big.Int) (bridgeFees, afterAmount *big.Int) {
	feeRate := GetGlobalFeeRate()
	baseTxFee := new(big.Int).Mul(BitcoinTxBytes, big.NewInt(int64(float64(feeRate)*BaseTxFeeMultiplier)))

	bridgeFees = new(big.Int).Mul(amount, bridgeFeeRate)
	bridgeFees = new(big.Int).Div(bridgeFees, big.NewInt(BridgeFeeRateDenominator))
	bridgeFees = new(big.Int).Add(bridgeFees, baseTxFee)

	afterAmount = new(big.Int).Sub(amount, bridgeFees)

	fields := map[string]interface{}{
		"amount":        amount,
		"afterAmount":   afterAmount,
		"bridgeFees":    bridgeFees,
		"bridgeFeeRate": bridgeFeeRate,
		"feeRate":       feeRate,
		"baseTxFee":     baseTxFee,
	}
	log.Logger().WithFields(fields).Info("completed the deduction to bitcoin bridge fees")
	return bridgeFees, afterAmount
}

func deductTONToEVMBridgeFees(amount, bridgeFeeRate *big.Int) (bridgeFees, afterAmount *big.Int) {
	bridgeFees = new(big.Int).Mul(amount, bridgeFeeRate)
	bridgeFees = new(big.Int).Div(bridgeFees, big.NewInt(BridgeFeeRateDenominator))
	bridgeFees = new(big.Int).Add(bridgeFees, TONToEVMBaseTxFee)

	afterAmount = new(big.Int).Sub(amount, bridgeFees)

	fields := map[string]interface{}{
		"amount":        amount,
		"afterAmount":   afterAmount,
		"bridgeFees":    bridgeFees,
		"bridgeFeeRate": bridgeFeeRate,
		"baseTxFee":     TONToEVMBaseTxFee,
	}
	log.Logger().WithFields(fields).Info("complete the deduction to evm bridge fees")
	return bridgeFees, afterAmount
}

func deductBitcoinToEVMBridgeFees(amount, bridgeFeeRate *big.Int) (bridgeFees, afterAmount *big.Int) {
	bridgeFees = new(big.Int).Mul(amount, bridgeFeeRate)
	bridgeFees = new(big.Int).Div(bridgeFees, big.NewInt(BridgeFeeRateDenominator))
	bridgeFees = new(big.Int).Add(bridgeFees, BitcoinToEVMBaseTxFee)

	afterAmount = new(big.Int).Sub(amount, bridgeFees)

	fields := map[string]interface{}{
		"amount":        amount,
		"afterAmount":   afterAmount,
		"bridgeFees":    bridgeFees,
		"bridgeFeeRate": bridgeFeeRate,
		"baseTxFee":     BitcoinToEVMBaseTxFee,
	}
	log.Logger().WithFields(fields).Info("complete the deduction to evm bridge fees")
	return bridgeFees, afterAmount
}

func calcProtocolFees(inAmountSat *big.Int, feeRatio uint64, wbtcDecimal uint64) *big.Int {
	if feeRatio == 0 {
		return big.NewInt(0)
	}

	fee := new(big.Int).Mul(inAmountSat, new(big.Int).SetUint64(feeRatio))
	fee = new(big.Int).Div(fee, big.NewInt(10000))
	if wbtcDecimal != params.BTCDecimal {
		fee = new(big.Int).Mul(fee, new(big.Int).SetUint64(wbtcDecimal))
		fee = new(big.Int).Div(fee, new(big.Int).SetUint64(params.BTCDecimal))
	}
	return fee
}

func convertDecimal(amount *big.Int, srcDecimal uint64, dstDecimal uint64) *big.Int {
	dstAmount := amount
	if srcDecimal > dstDecimal {
		exp := new(big.Int).Exp(big.NewInt(10), new(big.Int).SetUint64(srcDecimal-dstDecimal), nil)
		dstAmount = new(big.Int).Div(amount, exp)
	} else if srcDecimal < dstDecimal {
		exp := new(big.Int).Exp(big.NewInt(10), new(big.Int).SetUint64(dstDecimal-srcDecimal), nil)
		dstAmount = new(big.Int).Mul(amount, exp)
	}
	return dstAmount
}
