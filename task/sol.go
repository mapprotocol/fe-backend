package task

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/mapprotocol/fe-backend/dao"
	"github.com/mapprotocol/fe-backend/params"
	"github.com/mapprotocol/fe-backend/resource/log"
	"github.com/mapprotocol/fe-backend/third-party/filter"
	"github.com/mapprotocol/fe-backend/utils"
	"github.com/mapprotocol/fe-backend/utils/alarm"
	"github.com/spf13/viper"
	"github.com/xssnick/tonutils-go/tvm/cell"
	"io"
	"math/big"
	"net/http"
	"strconv"
	"time"
)

const (
	MaticChainId = 137
)

func FilterEventToSol() {
	chainID := params.ChainIDOfMaticPool
	topic := params.OnReceivedTopic
	filterLog := dao.NewFilterLog(chainID, topic)
	for {
		gotLog, err := filterLog.First()
		if err != nil {
			fields := map[string]interface{}{
				"chainID": params.ChainIDOfMaticPool,
				"topic":   params.OnReceivedTopic,
				"error":   err.Error(),
			}
			log.Logger().WithFields(fields).Error("failed to get filter log info")
			alarm.Slack(context.Background(), "failed to get filter log info")
			time.Sleep(5 * time.Second)
			continue
		}

		logs, err := filter.GetLogs(gotLog.LatestLogID, params.ChainIDOfMaticPool, params.OnReceivedTopic, uint8(20))
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

			if onReceived.DstChain.String() == params.SolChainID {
				orderId := make([]byte, 0, 32)
				for _, v := range onReceived.OrderId {
					orderId = append(orderId, v)
				}
				order := &dao.SolOrder{
					SrcChain:       onReceived.SrcChain.String(),
					SrcToken:       string(onReceived.SrcToken),
					Sender:         "0x" + common.Bytes2Hex(onReceived.Sender),
					InAmount:       onReceived.InAmount,
					RelayToken:     params.USDTOfChainPool,
					RelayAmount:    onReceived.ChainPoolTokenAmount.String(),
					DstChain:       onReceived.DstChain.String(),
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

			UpdateLogID(params.ChainIDOfMaticPool, params.OnReceivedTopic, lg.Id)
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
	endpoint := getEndpoint()
	client := rpc.New(endpoint)

	solCfg := viper.GetStringMapString("sol")
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
			break
		}

		for i, o := range orders {
			if i == length-1 {
				id = o.ID + 1
			}

			ele := o
			data, err := requestSolButter(solCfg["host"], routerPri.PublicKey().String(), ele)
			if err != nil {
				log.Logger().WithField("error", err.Error()).Error("failed to request sol swap")
				alarm.Slack(context.Background(), "failed to request sol swap")
				time.Sleep(5 * time.Second)
				continue
			}

			bbs, err := hex.DecodeString(data)
			if err != nil {
				log.Logger().WithField("error", err.Error()).Error("failed to hex data")
				alarm.Slack(context.Background(), "failed to hex data")
				time.Sleep(5 * time.Second)
				continue
			}
			trx, err := solana.TransactionFromBytes(bbs)
			if err != nil {
				log.Logger().WithField("error", err.Error()).Error("failed to get trx")
				alarm.Slack(context.Background(), "failed to get trx")
				time.Sleep(5 * time.Second)
				continue
			}
			resp, err := client.GetLatestBlockhash(context.Background(), rpc.CommitmentFinalized)
			if err != nil {
				log.Logger().WithField("error", err.Error()).Error("failed to getLatestBlockHash")
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
				log.Logger().WithField("error", err.Error()).Error("failed to sign trx")
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
				log.Logger().WithField("update", utils.JSON(update)).WithField("error", err.Error()).Error("failed to update sol order status")
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

func requestSolButter(host, router string, param *dao.SolOrder) (string, error) {
	url := fmt.Sprintf("%s/solanaCrossIn?fromChainId=%d&chainPoolChain=%d&"+
		"chainPoolTokenAddress=%s&chainPoolTokenAmount=%s&"+
		"tokenOutAddress=%s&fromChainTokenInAddress=%s&"+
		"fromChainTokenAmount=%s&slippage=%d&"+
		"router=%s&minAmountOut=%d&from=%s&orderIdHex=%d&receiver=%s",
		host, MaticChainId, MaticChainId,
		param.ChainPoolToken, param.RelayAmount,
		param.DstToken, param.SrcToken,
		param.InAmount, 100,
		router, param.Slippage, param.Sender, param.BridgeId, param.Receiver,
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
	fmt.Println("requestSolButter data ", string(body))
	data := SolButterData{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Logger().WithField("err", err).Error("failed to unmarshal body")
		return "", err
	}
	if data.Errno != 0 {
		return "", fmt.Errorf("code %d, mess:%s", data.Errno, data.Message)
	}
	return data.Data[0].TxParam[0].Data, nil
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
	} `json:"data"`
}
