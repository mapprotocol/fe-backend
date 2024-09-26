package task

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/mapprotocol/fe-backend/third-party/butter"
	"io"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/mapprotocol/fe-backend/dao"
	"github.com/mapprotocol/fe-backend/params"
	"github.com/mapprotocol/fe-backend/resource/log"
	"github.com/mapprotocol/fe-backend/third-party/filter"
	"github.com/mapprotocol/fe-backend/utils"
	"github.com/mapprotocol/fe-backend/utils/alarm"
	"github.com/mapprotocol/fe-backend/utils/tx"
	"github.com/mr-tron/base58"
	"github.com/spf13/viper"
)

const (
	MaticChainId = 137
)

func FilterEventToSol() {
	chainID := params.ChainIDOfSolChainPool
	topic := params.OnReceivedTopic
	filterLog := dao.NewFilterLog(chainID, topic)
	for {
		gotLog, err := filterLog.First()
		if err != nil {
			fields := map[string]interface{}{
				"chainID": params.ChainIDOfSolChainPool,
				"topic":   params.OnReceivedTopic,
				"error":   err.Error(),
			}
			log.Logger().WithFields(fields).Error("failed to get filter log info")
			alarm.Slack(context.Background(), "failed to get filter log info")
			time.Sleep(5 * time.Second)
			continue
		}

		logs, err := filter.GetLogs(gotLog.LatestLogID, params.ChainIDOfSolChainPool, params.OnReceivedTopic, uint8(20))
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
			log.Logger().WithField("id", gotLog.LatestLogID).Info("running")
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
					"chainID": params.ChainIDOfSolChainPool,
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
					"chainID": params.ChainIDOfSolChainPool,
					"topic":   params.OnReceivedTopic,
					"logData": lg.LogData,
					"error":   err.Error(),
				}
				log.Logger().WithFields(fields).Error("failed to unpack log data")
				alarm.Slack(context.Background(), "failed to unpack log data")
				UpdateLogID(chainID, topic, lg.Id)
				continue
			}

			if onReceived.DstChain.String() == params.SolChainID {
				orderId := make([]byte, 0, 32)
				for _, v := range onReceived.OrderId {
					orderId = append(orderId, v)
				}
				order := &dao.SolOrder{
					SrcHash:        lg.TxHash,
					SrcChain:       onReceived.SrcChain.String(),
					SrcToken:       common.BytesToAddress(onReceived.SrcToken).String(),
					Sender:         "0x" + common.Bytes2Hex(onReceived.Sender),
					InAmount:       onReceived.InAmount,
					RelayToken:     params.USDTOfChainPool,
					RelayAmount:    onReceived.ChainPoolTokenAmount.String(),
					DstChain:       onReceived.DstChain.String(), // onReceived.DstChain.String(),
					DstToken:       string(onReceived.DstToken),
					Receiver:       string(onReceived.Receiver),
					Action:         dao.OrderActionFromEVM,
					Stage:          dao.OrderStag1,
					Status:         dao.OrderStatusTxConfirmed,
					Slippage:       onReceived.Slippage,
					OrderId:        "0x" + common.Bytes2Hex(orderId),
					ChainPoolToken: onReceived.ChainPoolToken.Hex(),
					BridgeId:       onReceived.BridgeId,
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

			UpdateLogID(params.ChainIDOfSolChainPool, params.OnReceivedTopic, lg.Id)
			time.Sleep(time.Second * 5)
		}
		time.Sleep(time.Second * 5)
	}
}

func HandlerEvm2Sol() {
	order := dao.SolOrder{
		DstChain: params.SolChainID,
		Status:   dao.OrderStatusTxConfirmed,
	}
	//endpoint := getEndpoint()
	endpointCfg := viper.GetStringMapString("endpoints")
	solCfg := viper.GetStringMapString("sol")
	client := rpc.New(endpointCfg["solana"])
	routerPri, err := solana.PrivateKeyFromBase58(solCfg["pri"])
	if err != nil {
		panic(err)
	}
	for id := uint64(1); ; {
		orders, err := order.GetOldest10ByID(id)
		if err != nil {
			fields := map[string]interface{}{
				"id":    id,
				"func":  "HandlerEvm2Sol",
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
			time.Sleep(10 * time.Second)
			log.Logger().Info("HandlerEvm2Sol is running, wait new transaction")
			continue
		}

		for i, o := range orders {
			if i == length-1 {
				id = o.ID + 1
			}

			log.Logger().Info("HandlerEvm2Sol srcHash =", o.SrcHash)
			ele := o
			data, err := requestSolButter(endpointCfg["butter"], routerPri.PublicKey().String(), ele)
			if err != nil {
				log.Logger().WithField("id", o.ID).WithField("error", err.Error()).Error("failed to request sol butter")
				alarm.Slack(context.Background(), "failed to request sol butter")
				time.Sleep(5 * time.Second)
				continue
			}

			bbs, err := hex.DecodeString(data)
			if err != nil {
				log.Logger().WithField("id", o.ID).WithField("error", err.Error()).Error("failed to hex data")
				alarm.Slack(context.Background(), "failed to hex data")
				time.Sleep(5 * time.Second)
				continue
			}
			trx, err := solana.TransactionFromBytes(bbs)
			if err != nil {
				log.Logger().WithField("id", o.ID).WithField("error", err.Error()).Error("failed to get trx")
				alarm.Slack(context.Background(), "failed to get trx")
				time.Sleep(5 * time.Second)
				continue
			}
			resp, err := client.GetLatestBlockhash(context.Background(), rpc.CommitmentFinalized)
			if err != nil {
				log.Logger().WithField("id", o.ID).WithField("error", err.Error()).Error("failed to getLatestBlockHash")
				alarm.Slack(context.Background(), "failed to getLatestBlockHash")
				time.Sleep(5 * time.Second)
				continue
			}
			trx.Message.RecentBlockhash = resp.Value.Blockhash
			// sign
			_, err = trx.Sign(func(key solana.PublicKey) *solana.PrivateKey {
				if key == routerPri.PublicKey() {
					return &routerPri
				}
				return nil
			})
			if err != nil {
				log.Logger().WithField("id", o.ID).WithField("error", err.Error()).Error("failed to sign trx")
				alarm.Slack(context.Background(), "failed to sign trx")
				time.Sleep(5 * time.Second)
				continue
			}
			sig, err := client.SendTransaction(context.TODO(), trx)
			if err != nil {
				log.Logger().WithField("error", err.Error()).Error("failed to send trx")
				alarm.Slack(context.Background(), "failed to send trx")
				time.Sleep(5 * time.Second)
				continue
			}

			update := &dao.SolOrder{
				Stage:     dao.OrderStag2,
				Status:    dao.OrderStatusTxSent,
				OutTxHash: sig.String(),
			}
			if err := dao.NewSolOrderWithID(o.ID).Updates(update); err != nil {
				log.Logger().WithField("id", o.ID).WithField("update", utils.JSON(update)).WithField("error", err.Error()).Error("failed to update sol order status")
				alarm.Slack(context.Background(), "failed to update sol order status")
				time.Sleep(5 * time.Second)
				continue
			}
			time.Sleep(time.Second)
		}
		time.Sleep(10 * time.Second)
	}
}

func FilterSol2Evm() {
	chainID := params.SolChainID
	topic := params.EventIDSolToEVM
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

			logData := make(map[string]interface{})
			err = json.Unmarshal([]byte(lg.LogData), &logData)
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

			srcChainStr := logData["fromChainId"].(string)
			srcChain, ok := big.NewInt(0).SetString(srcChainStr, 16)
			if !ok {
				fields := map[string]interface{}{
					"logID":       lg.Id,
					"chainID":     chainID,
					"srcChainStr": srcChainStr,
				}
				log.Logger().WithFields(fields).Error("parse fromChainId failed")
				alarm.Slack(context.Background(), "parse fromChainId failed, str("+srcChainStr+")")
				UpdateLogID(chainID, topic, lg.Id)
				continue
			}

			toChainStr := logData["toChain"].(string)
			toChain, ok := big.NewInt(0).SetString(toChainStr, 16)
			if !ok {
				fields := map[string]interface{}{
					"logID":      lg.Id,
					"chainID":    chainID,
					"toChainStr": toChainStr,
				}
				log.Logger().WithFields(fields).Error("parse toChain failed")
				alarm.Slack(context.Background(), "parse toChain failed, str("+toChainStr+")")
				UpdateLogID(chainID, topic, lg.Id)
				continue
			}

			minAmountOutStr := logData["minAmountOut"].(string)
			minAmountOut, ok := big.NewInt(0).SetString(minAmountOutStr, 16)
			if !ok {
				fields := map[string]interface{}{
					"logID":      lg.Id,
					"chainID":    chainID,
					"toChainStr": toChainStr,
				}
				log.Logger().WithFields(fields).Error("parse minAmountOut failed")
				alarm.Slack(context.Background(), "parse minAmountOut failed, str("+toChainStr+")")
				UpdateLogID(chainID, topic, lg.Id)
				continue
			}

			amountOutStr := logData["amountOut"].(string)
			amountOut, ok := big.NewInt(0).SetString(amountOutStr, 16)
			if !ok {
				fields := map[string]interface{}{
					"logID":        lg.Id,
					"chainID":      chainID,
					"amountOutStr": amountOutStr,
				}
				log.Logger().WithFields(fields).Error("parse amountOut failed")
				alarm.Slack(context.Background(), "parse amountOut failed, str("+amountOutStr+")")
				UpdateLogID(chainID, topic, lg.Id)
				continue
			}

			orderIdStr := logData["orderId"].(string)
			orderId := big.NewInt(0).SetBytes(common.Hex2Bytes(strings.TrimPrefix(orderIdStr, "ff")))

			fromStr := logData["from"].([]interface{})
			from := base58.Encode(convert2Bytes(fromStr))

			fromTokenStr := logData["fromToken"].([]interface{})
			fromToken := base58.Encode(convert2Bytes(fromTokenStr))

			toTokenStr := logData["toToken"].([]interface{})
			toToken := common.BytesToAddress(convert2Bytes(toTokenStr[:20]))

			receiverTokenStr := logData["receiver"].([]interface{})
			receiver := common.BytesToAddress(convert2Bytes(receiverTokenStr[:20]))

			afterBalanceStr := logData["afterBalance"].(string)
			afterBalance, ok := big.NewInt(0).SetString(afterBalanceStr, 16)
			if !ok {
				fields := map[string]interface{}{
					"logID":           lg.Id,
					"chainID":         chainID,
					"afterBalanceStr": afterBalanceStr,
				}
				log.Logger().WithFields(fields).Error("parse tokenAmount failed")
				alarm.Slack(context.Background(), "parse tokenAmount failed, str("+afterBalanceStr+")")
				UpdateLogID(chainID, topic, lg.Id)
				continue
			}

			swapTokenOutBeforeBalanceStr := logData["swapTokenOutBeforeBalance"].(string)
			swapTokenOutBeforeBalance, ok := big.NewInt(0).SetString(swapTokenOutBeforeBalanceStr, 16)
			if !ok {
				fields := map[string]interface{}{
					"logID":                     lg.Id,
					"chainID":                   chainID,
					"swapTokenOutBeforeBalance": swapTokenOutBeforeBalanceStr,
				}
				log.Logger().WithFields(fields).Error("parse tokenAmount failed")
				alarm.Slack(context.Background(), "parse tokenAmount failed, str("+swapTokenOutBeforeBalanceStr+")")
				UpdateLogID(chainID, topic, lg.Id)
				continue
			}
			relayAmount := big.NewInt(0).Sub(afterBalance, swapTokenOutBeforeBalance)

			order := &dao.Order{
				OrderIDFromContract: uint64(orderId.Int64()),
				SrcChain:            srcChain.String(),
				DstChain:            toChain.String(),
				SrcToken:            fromToken,
				Sender:              from,
				InAmount:            amountOut.String(),
				RelayToken:          params.USDCOfSOL,
				RelayAmount:         relayAmount.String(),
				DstToken:            toToken.Hex(),
				Receiver:            receiver.Hex(),
				Action:              dao.OrderActionToEVM,
				Stage:               dao.OrderStag1,
				Status:              dao.OrderStatusTxConfirmed,
				MinAmountOut:        minAmountOut.String(),
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
				// UpdateLogID(chainID, topic, lg.Id)
				continue
			}

			UpdateLogID(chainID, topic, lg.Id)
		}

		time.Sleep(10 * time.Second)
	}
}

// HandleSol2EvmButter
// action=1, stage=1, status=4(TxConfirmed) ==> stage=2, status=1(TxPrepareSend) ==> stage=2, status=2(TxSent)
func HandleSol2EvmButter() {
	order := dao.Order{
		SrcChain: params.SolChainID,
		Action:   dao.OrderActionToEVM,
		Stage:    dao.OrderStag1,
		Status:   dao.OrderStatusTxConfirmed,
	}
	//endpointCfg := viper.GetStringMapString("endpoints")
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
				log.Logger().Info("HandleSol2EvmButter is running")
				time.Sleep(10 * time.Second)
				break
			}

			for i, o := range orders {
				if i == length-1 {
					id = o.ID + 1
				}

				decimal := params.USDTDecimalOfEthereum
				chainIDOfChainPool := params.ChainIDOfSolChainPool
				chainInfo := &dao.ChainPool{}
				chainInfo, err = dao.NewChainPoolWithChainID(params.ChainIDOfSolChainPool).First()
				if err != nil {
					log.Logger().WithField("chainID", params.ChainIDOfSolChainPool).WithField("error", err.Error()).Error("failed to get chain info")
					alarm.Slack(context.Background(), "failed to get chain info")
					time.Sleep(5 * time.Second)
					continue
				}

				before, _ := big.NewFloat(0).SetString(o.InAmount)
				amount := before.Quo(before, big.NewFloat(decimal)).String()
				// step1: 请求 route 接口，获取路由
				req := butter.RouterRequest{
					FromChainID:     chainIDOfChainPool,
					ToChainID:       o.DstChain,
					Amount:          amount, // 处理精度
					TokenInAddress:  chainInfo.USDTContract,
					TokenOutAddress: o.DstToken,
					Type:            SwapType,
					Slippage:        150,
				}
				data, err := butter.Route(&req)
				if err != nil {
					log.Logger().WithField("error", err.Error()).Error("failed to butter route info")
					alarm.Slack(context.Background(), "failed to butter route info")
					time.Sleep(5 * time.Second)
					continue
				}

				fmt.Println("order ------------------- ", "0xff"+common.Bytes2Hex(big.NewInt(0).SetUint64(o.OrderIDFromContract).Bytes()))
				// step2: 请求 evmCrossInSwap接口，获取交易信息
				swapReq := butter.EvmCrossInSwapRequest{
					Hash:         data.Data[0].Hash,
					SrcChainId:   o.SrcChain,
					From:         o.Sender,
					Router:       viper.GetStringMapString("chainPool")["sender"],
					Receiver:     o.Receiver,
					MinAmountOut: o.MinAmountOut,
					OrderIdHex:   "0xff" + common.Bytes2Hex(big.NewInt(0).SetUint64(o.OrderIDFromContract).Bytes()), // orderId处理
					Fee:          "0",
					FeeReceiver:  "0x0000000000000000000000000000000000000000",
				}
				swapResp, err := butter.EvmCrossInSwap(&swapReq)
				if err != nil {
					log.Logger().WithField("error", err.Error()).Error("failed to butter evmCrossInSwap")
					alarm.Slack(context.Background(), "failed to butter evmCrossInSwap")
					time.Sleep(5 * time.Second)
					continue
				}
				// step3: 发送交易
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
				transactor, err := tx.NewTransactor(chainInfo.ChainRPC, swapResp.Data[0].To, multiplier)
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

				value, _ := big.NewInt(0).SetString(strings.TrimPrefix(swapResp.Data[0].Value, "0x"), 16)
				txHash, err := transactor.SendCustom(common.HexToAddress(swapResp.Data[0].To), value,
					common.Hex2Bytes(strings.TrimPrefix(swapResp.Data[0].Data, "0x")))
				if err != nil {
					fields := map[string]interface{}{
						"error": err,
					}
					log.Logger().WithFields(fields).Error("failed to send transactor")
					alarm.Slack(context.Background(), "failed to send transactor")
					time.Sleep(5 * time.Second)
					continue
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

func HandleSol2EvmFinish() {
	order := dao.Order{
		SrcChain: params.SolChainID,
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
						log.Logger().WithField("chainID", params.ChainIDOfSolChainPool).WithField("error", err.Error()).Error("failed to get chain info")
						alarm.Slack(context.Background(), "failed to get chain info")
						time.Sleep(5 * time.Second)
						continue

					}
				} else {
					chainInfo, err = dao.NewChainPoolWithChainID(params.ChainIDOfSolChainPool).First()
					if err != nil {
						log.Logger().WithField("chainID", params.ChainIDOfSolChainPool).WithField("error", err.Error()).Error("failed to get chain info")
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

func convert2Bytes(data []interface{}) []byte {
	ret := make([]byte, 0, len(data))
	for _, d := range data {
		switch d.(type) {
		case float64:
			ret = append(ret, byte(d.(float64)))
		}
	}

	return ret
}

func requestSolButter(host, router string, param *dao.SolOrder) (string, error) {
	orderIdHex := big.NewInt(0).SetUint64(param.BridgeId)
	url := fmt.Sprintf("%s/solanaCrossIn?fromChainId=%s&chainPoolChain=%d&"+
		"chainPoolTokenAddress=%s&chainPoolTokenAmount=%s&"+
		"tokenOutAddress=%s&fromChainTokenInAddress=%s&"+
		"fromChainTokenAmount=%s&slippage=%d&"+
		"router=%s&minAmountOut=%d&from=%s&orderIdHex=%s&receiver=%s",
		host, param.SrcChain, MaticChainId,
		param.ChainPoolToken, param.RelayAmount,
		param.DstToken, param.SrcToken,
		param.InAmount, 100,
		router, param.Slippage, param.Sender, "0x"+common.Bytes2Hex(orderIdHex.Bytes()), param.Receiver,
	)
	fmt.Println("requestSolButter url ", url)
	resp, err := http.Get(url)
	if err != nil {
		log.Logger().WithField("err", err).Error("failed to get response")
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Logger().WithField("err", err).Error("failed to readAll body")
		return "", err
	}
	//fmt.Println("requestSolButter data ", string(body))
	data := SolButterData{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Logger().WithField("err", err).Error("failed to unmarshal body")
		return "", err
	}
	if data.Errno != 0 {
		return "", fmt.Errorf("code %d, mess:%s", data.Errno, data.Message)
	}

	if data.Data[0].Error.Response.Errno != 0 {
		if data.Data[0].Error.Response.Message == "Invalid min amount out" {
			swapData, err := requestRouteAndSwap(param)
			if err != nil {
				return "", fmt.Errorf("invalid min amount out failed to requestRouteAndSwap, err:%vv", err)
			}
			min, ok := big.NewFloat(0).SetString(swapData.Data[0].Route.MinAmountOut.Amount)
			if !ok {
				return "", fmt.Errorf("invalid min amount out, swapData MinAmountOut not float64")
			}
			slippage, err := convertSol(swapData.Data[0].Route.MinAmountOut.Symbol, min)
			if err != nil {
				return "", fmt.Errorf("invalid min amount out, swapData MinAmountOut convert failed, err:%v", err)
			}
			fmt.Println("slippage ---------------------- ", slippage)
			//time.Sleep(time.Minute * 10)
			param.Slippage = slippage
			return requestSolButter(host, router, param)
		}
		return "", fmt.Errorf("code %d, mess:%s", data.Data[0].Error.Response.Errno, data.Data[0].Error.Response.Message)
	}

	return data.Data[0].TxParam[0].Data, nil
}

func convertSol(symbol string, data *big.Float) (uint64, error) {
	switch strings.ToLower(symbol) {
	case "sol":
		data = data.Mul(data, big.NewFloat(1e9))
	}
	ret, _ := data.Int64()
	return uint64(ret), nil
}

// amount=10&slippage=300&receiver=FjgdFp1zt8A2fttq71fprwjbnJvRYEg2VLoQaLxTdFgR&from=0xa8EE0cf2Af6fE245090801d36E69281BC6610F29&referrer=0xa8EE0cf2Af6fE245090801d36E69281BC6610F29&rateOrNativeFee=50&feeType=1
func requestRouteAndSwap(param *dao.SolOrder) (*butter.RouterAndSwapResponse, error) {
	before, _ := big.NewFloat(0).SetString(param.RelayAmount)
	amount := before.Quo(before, big.NewFloat(params.USDTDecimalOfEthereum)).String()
	request := &butter.RouterAndSwapRequest{
		FromChainID:     params.ChainIDOfSolChainPool,
		ToChainID:       param.DstChain,
		Amount:          amount, // decimal
		TokenInAddress:  params.USDTOfChainPool,
		TokenOutAddress: param.DstToken,
		Type:            SwapType,
		Slippage:        300,
		From:            param.Sender,
		Receiver:        param.Receiver,
		Referrer:        sender,
		RateOrNativeFee: "50",
	}
	return butter.RouteAndSwapSol(request)
}

type SolButterData struct {
	Errno   int    `json:"errno"`
	Message string `json:"message"`
	Data    []struct {
		Route struct {
			Diff      string `json:"diff"`
			BridgeFee struct {
				Amount string `json:"amount"`
			} `json:"bridgeFee"`
			TradeType int `json:"tradeType"`
			GasFee    struct {
				Amount string `json:"amount"`
				Symbol string `json:"symbol"`
			} `json:"gasFee"`
			SwapFee struct {
				NativeFee string `json:"nativeFee"`
				TokenFee  string `json:"tokenFee"`
			} `json:"swapFee"`
			FeeConfig struct {
				FeeType         int    `json:"feeType"`
				Referrer        string `json:"referrer"`
				RateOrNativeFee int    `json:"rateOrNativeFee"`
			} `json:"feeConfig"`
			GasEstimated       string `json:"gasEstimated"`
			GasEstimatedTarget string `json:"gasEstimatedTarget"`
			TimeEstimated      int    `json:"timeEstimated"`
			Hash               string `json:"hash"`
			Timestamp          int64  `json:"timestamp"`
			HasLiquidity       bool   `json:"hasLiquidity"`
			SrcChain           struct {
				ChainID string `json:"chainId"`
				TokenIn struct {
					Address  string `json:"address"`
					Name     string `json:"name"`
					Decimals int    `json:"decimals"`
					Symbol   string `json:"symbol"`
					Icon     string `json:"icon"`
				} `json:"tokenIn"`
				TokenOut struct {
					Address  string `json:"address"`
					Name     string `json:"name"`
					Decimals int    `json:"decimals"`
					Symbol   string `json:"symbol"`
					Icon     string `json:"icon"`
				} `json:"tokenOut"`
				TotalAmountIn  string `json:"totalAmountIn"`
				TotalAmountOut string `json:"totalAmountOut"`
				Route          []struct {
					AmountIn  string        `json:"amountIn"`
					AmountOut string        `json:"amountOut"`
					DexName   string        `json:"dexName"`
					Path      []interface{} `json:"path"`
				} `json:"route"`
				Bridge string `json:"bridge"`
			} `json:"srcChain"`
			MinAmountOut struct {
				Amount string `json:"amount"`
				Symbol string `json:"symbol"`
			} `json:"minAmountOut"`
		} `json:"route"`
		TxParam []struct {
			To      string `json:"to"`
			ChainID string `json:"chainId"`
			Data    string `json:"data"`
			Value   string `json:"value"`
			Method  string `json:"method"`
		} `json:"txParam"`
		Error struct {
			Response struct {
				Errno   int    `json:"errno"`
				Message string `json:"message"`
			} `json:"response"`
			Status  int    `json:"status"`
			Message string `json:"message"`
			Name    string `json:"name"`
		} `json:"error"`
	} `json:"data"`
}
