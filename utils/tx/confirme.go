package tx

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"time"
)

const (
	WaitRetryLimit    = 10
	WaitRetryInterval = time.Second
)

//var ErrTxExecutionFailed = errors.New("transaction execution failed")

//func WaitForTxConfirmation(txHash common.Hash) error {
//	for i := 0; i < WaitRetryLimit; i++ {
//		time.Sleep(WaitRetryInterval)
//		_, isPending, err := transactor.client.TransactionByHash(context.Background(), txHash)
//		if err != nil {
//			log.Error("get transaction by hash failed, will retry", "times", i+1, "hash", txHash, "error", err)
//			continue
//		}
//		if !isPending {
//			break
//		}
//	}
//
//	receipt, err := getTransactionReceipt(txHash)
//	if err != nil {
//		return err
//	}
//	if receipt.Status != types.ReceiptStatusSuccessful {
//		return ErrTxExecutionFailed
//	}
//	return err
//}

//func getTransactionReceipt(client *ethclient.Client, txHash common.Hash) (receipt *types.Receipt, err error) {
//	for i := 0; i < WaitRetryLimit; i++ {
//		receipt, err = client.TransactionReceipt(context.Background(), txHash)
//		if err == nil {
//			return receipt, nil
//		}
//
//		log.Error("get transaction receipt failed, will retry", "times", i+1, "hash", txHash, "error", err)
//		time.Sleep(WaitRetryInterval)
//	}
//
//	return nil, err
//}

func (t *Transactor) TransactionIsPending(txHash common.Hash) (bool, error) {
	_, isPending, err := t.client.TransactionByHash(context.Background(), txHash)
	if err != nil {
		log.Error("filed to get transaction by hash", "hash", txHash, "error", err)
		return false, err
	}
	return isPending, nil
}

func (t *Transactor) TransactionStatus(txHash common.Hash) (uint64, error) {
	receipt, err := t.client.TransactionReceipt(context.Background(), txHash)
	if err != nil {
		log.Error("failed to get transaction receipt", "hash", txHash, "error", err)
		return 0, err
	}
	return receipt.Status, nil
}
