package tx

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/mapprotocol/fe-backend/bindings/ferouter"
	"github.com/mapprotocol/fe-backend/utils"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/mapprotocol/fe-backend/resource/log"
)

const (
	RetryTimes  = 20
	Interval    = 2 * time.Second
	NonceTooLow = "nonce too low"
)

func (t *Transactor) SendTransaction(to common.Address, value *big.Int, input []byte) (common.Hash, error) {
	var txHash common.Hash

	for i := 0; i < RetryTimes; i++ {
		opts, err := t.NewTransactOpts()
		if err != nil {
			return common.Hash{}, err
		}
		//t.gas()

		log.Logger().WithField("nonce", opts.Nonce).Info("will send transaction")
		txHash, err = t.sendTransaction(t.privateKey, to, value, input)
		if err != nil {
			if isNonceTooLow(err.Error()) {
				log.Logger().WithField("nonce", opts.Nonce).Warn("nonce too low, will try again in 2 second")
				time.Sleep(Interval)
				continue
			}
			return common.Hash{}, err
		}
		break
	}
	return txHash, nil
}

//func (t *Transactor) DeliverAndSwap(orderID [32]byte, initiator common.Address, token common.Address, amount *big.Int, swapData []byte, bridgeData []byte, feeData []byte, value *big.Int) (common.Hash, error) {
//	var txHash common.Hash
//
//	for i := 0; i < RetryTimes; i++ {
//		opts, err := t.NewTransactOpts()
//		if err != nil {
//			return common.Hash{}, err
//		}
//
//		gasPrice, gasLimit, err := t.gas(ferouter.FerouterMetaData.ABI, "deliverAndSwap", t.chainPoolContract, orderID, initiator, token, amount, swapData, bridgeData, feeData)
//		if err != nil {
//			return common.Hash{}, err
//		}
//		//opts.GasPrice = gasPrice
//		//opts.GasLimit = gasLimit
//		//opts.Value = value
//
//		log.Logger().WithField("nonce", opts.Nonce).Info("will send transaction")
//		tx, err := t.chainPoolTransactor.DeliverAndSwap(opts, orderID, initiator, token, amount, swapData, bridgeData, feeData)
//		if err != nil {
//			if isNonceTooLow(err.Error()) {
//				log.Logger().WithField("nonce", opts.Nonce).Warn("nonce too low, will try again in 2 second")
//				time.Sleep(Interval)
//				continue
//			}
//			//if i == RetryTimes-1 {
//			//	log.Logger().WithField("orderID", orderID).WithField("error", err).Warn("failed to send deliver and swap transaction")
//			//	alarm.Slack(context.Background(), fmt.Sprintf("failed to send deliver and swap transaction, orderID: %v, error: %v", orderID, err))
//			//	return common.Hash{}, err
//			//}
//			return common.Hash{}, err
//		}
//		if tx == nil {
//			log.Logger().WithField("orderID", orderID).WithField("error", err).Warn("completed to send deliver and swap transaction but tx is nil")
//			break
//		}
//		return tx.Hash(), err
//	}
//	return txHash, nil
//}

func (t *Transactor) DeliverAndSwap(orderID [32]byte, initiator common.Address, token common.Address, amount *big.Int, swapData []byte, bridgeData []byte, feeData []byte, value *big.Int) (common.Hash, error) {
	var txHash common.Hash

	for i := 0; i < RetryTimes; i++ {
		opts, err := t.NewTransactOpts()
		if err != nil {
			return common.Hash{}, err
		}

		input, err := pack(ferouter.FerouterMetaData.ABI, "deliverAndSwap", orderID, initiator, token, amount, swapData, bridgeData, feeData)
		if err != nil {
			log.Logger().Error("failed to pack params")
			return common.Hash{}, err
		}
		//opts.GasPrice = gasPrice
		//opts.GasLimit = gasLimit
		//opts.Value = value

		log.Logger().WithField("nonce", opts.Nonce).Info("will send transaction")
		txHash, err = t.sendTransaction(t.privateKey, t.chainPoolContract, value, input)
		if err != nil {
			if isNonceTooLow(err.Error()) {
				log.Logger().WithField("nonce", opts.Nonce).Warn("nonce too low, will try again in 2 second")
				time.Sleep(Interval)
				continue
			}
			return common.Hash{}, err
		}
		break
	}
	return txHash, nil
}

func (t *Transactor) sendTransaction(privateKey *ecdsa.PrivateKey, to common.Address, value *big.Int, input []byte) (common.Hash, error) {
	nonce, err := t.client.PendingNonceAt(context.Background(), t.address)
	if err != nil {
		fields := map[string]interface{}{
			"rpc":     t.endpoint,
			"address": t.address.Hex(),
			"error":   err,
		}
		log.Logger().WithFields(fields).Error("failed to get chain id")
		return common.Hash{}, err
	}

	gasPrice, err := t.client.SuggestGasPrice(context.Background())
	if err != nil {
		fields := map[string]interface{}{
			"rpc":   t.endpoint,
			"error": err,
		}
		log.Logger().WithFields(fields).Error("failed to get suggest gas price")
		return common.Hash{}, err
	}

	msg := ethereum.CallMsg{From: t.address, To: &to, GasPrice: gasPrice, Value: value, Data: input}
	gasLimit, err := t.client.EstimateGas(context.Background(), msg)
	if err != nil {
		fields := map[string]interface{}{
			"rpc":   t.endpoint,
			"msg":   utils.JSON(msg),
			"error": err,
		}
		log.Logger().WithFields(fields).Error("failed to get estimate gas")
		return common.Hash{}, err
	}

	fmt.Println("============================== gasLimit: ", gasLimit)
	if t.gasLimitMultiplier > 1 && gasLimit > 0 {
		gasLimit = uint64(float64(gasLimit) * t.gasLimitMultiplier)
	}
	if gasLimit < 1 {
		gasLimit = 2100000
	}
	fmt.Println("============================== gasLimit: ", gasLimit)

	txData := &types.LegacyTx{
		Nonce:    nonce,
		To:       &to,
		Value:    value,
		Gas:      gasLimit,
		GasPrice: gasPrice,
		Data:     input,
	}
	tx := types.NewTx(txData)
	chainID, err := t.client.ChainID(context.Background())
	if err != nil {
		log.Logger().WithField("rpc", t.endpoint).WithField("error", err).Error("failed to get chain id")
	}

	signer := types.LatestSignerForChainID(chainID)
	signedTx, err := types.SignTx(tx, signer, privateKey)
	if err != nil {
		fields := map[string]interface{}{
			"txData":  utils.JSON(txData),
			"chainID": chainID,
			"error":   err,
		}
		log.Logger().WithFields(fields).Error("failed to sign tx")
		return common.Hash{}, err
	}

	err = t.client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		fields := map[string]interface{}{
			"txData":  utils.JSON(txData),
			"chainID": chainID,
			"error":   err,
		}
		log.Logger().WithFields(fields).Error("failed to send tx")
		return common.Hash{}, err
	}
	return signedTx.Hash(), nil
}

func (t *Transactor) estimateGas(to common.Address, input []byte) (*big.Int, uint64, error) {
	gasPrice, err := t.client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, 0, err
	}

	msg := ethereum.CallMsg{From: t.address, To: &to, GasPrice: gasPrice, Value: nil, Data: input}
	fmt.Println("============================== ")
	gasLimit, err := t.client.EstimateGas(context.Background(), msg)
	if err != nil {
		fmt.Println("============================== estimate gas error: ", err)
		return nil, 0, err
	}
	fmt.Printf("============================== gasPrice: %v, gasLimit: %v\n", gasPrice, gasLimit)
	if t.gasLimitMultiplier > 1 {
		gasLimit = uint64(float64(gasLimit) * t.gasLimitMultiplier)
	}
	fmt.Printf("============================== gasPrice: %v, gasLimit: %v\n", gasPrice, gasLimit)

	return gasPrice, gasLimit, nil
}

func pack(abiStr, method string, args ...interface{}) ([]byte, error) {
	parsed, err := abi.JSON(strings.NewReader(abiStr))
	if err != nil {
		return nil, err
	}
	input, err := parsed.Pack(method, args...) // todo do once
	if err != nil {
		return nil, err
	}
	return input, nil
}

func isNonceTooLow(s string) bool {
	return strings.Contains(s, NonceTooLow)
}
