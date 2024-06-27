package logic

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	btcmempool "github.com/btcsuite/btcd/mempool"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/mapprotocol/fe-backend/dao"
	"github.com/mapprotocol/fe-backend/resource/log"
	"github.com/mapprotocol/fe-backend/third-party/mempool"
	"github.com/mapprotocol/fe-backend/utils/alarm"
)

const (
	defaultSequenceNum = wire.MaxTxInSequenceNum - 10
	CollectFinish      = 1
)

var (
	//PrevAdminOutPoint2        *PrevOutPoint = nil
	//MinPreAdminOutPointValue2               = int64(20000)
	MinBalanceInFeeAddress = int64(20000)
	NoMoreUTXO             = errors.New("no more utxo")
)

type PrevOutPoint struct {
	Outpoint *wire.OutPoint
	Value    int64
}
type CollectCfg struct {
	Testnet       bool
	StrFeePrivkey string
	FeeAddress    btcutil.Address
	Receiver      btcutil.Address
}
type OrderItem struct {
	OrderID uint64
	Sender  btcutil.Address
	Priv    *btcec.PrivateKey
	Amount  int64
}

func gatherUTXOForItem(sender btcutil.Address, client *mempool.MempoolClient) ([]*PrevOutPoint, error) {
	outPointList := make([]*PrevOutPoint, 0)
	unspentList, err := client.ListUnspent(sender)
	if err != nil {
		return nil, err
	}
	if len(unspentList) == 0 {
		return nil, NoMoreUTXO
	}
	for i := range unspentList {
		outPointList = append(outPointList, &PrevOutPoint{
			Outpoint: unspentList[i].Outpoint,
			Value:    unspentList[i].Output.Value,
		})
	}
	return outPointList, nil
}
func gatherUTXO3(sender btcutil.Address, client *mempool.MempoolClient) ([]*PrevOutPoint, error) {
	outPointList := make([]*PrevOutPoint, 0)
	unspentList, err := client.ListUnspent(sender)
	if err != nil {
		return nil, err
	}

	if len(unspentList) == 0 {
		return nil, NoMoreUTXO
	}

	for i := range unspentList {
		if unspentList[i].Output.Value < 5000 {
			continue
		}
		outPointList = append(outPointList, &PrevOutPoint{
			Outpoint: unspentList[i].Outpoint,
			Value:    unspentList[i].Output.Value,
		})
	}
	return outPointList, nil
}
func getTxOutByOutPoint2(outPoint *wire.OutPoint, btcClient *mempool.MempoolClient) (*wire.TxOut, error) {
	tx, err := btcClient.GetRawTransaction(&outPoint.Hash)
	if err != nil {
		return nil, err
	}
	if int(outPoint.Index) >= len(tx.TxOut) {
		return nil, errors.New("err out point")
	}
	return tx.TxOut[outPoint.Index], nil
}

// the last item is the fee utxos and private key in the privs and outlists
func makeCollectTx0(feerate int64, receiverAddress, feeAddress btcutil.Address, outLists [][]*PrevOutPoint,
	privs []*btcec.PrivateKey, btcApiClient *mempool.MempoolClient) (*wire.MsgTx, error) {

	commitTx := wire.NewMsgTx(wire.TxVersion)
	totalSenderAmount, totalAmount := btcutil.Amount(0), btcutil.Amount(0)
	TxPrevOutputFetcher := txscript.NewMultiPrevOutFetcher(nil)
	count, pos := len(outLists), 0
	tmpPrivs := make(map[int]*btcec.PrivateKey)
	// handle the every address's utxo
	for i, outs := range outLists {
		for _, out := range outs {
			txOut, err := getTxOutByOutPoint2(out.Outpoint, btcApiClient)
			if err != nil {
				return nil, err
			}
			TxPrevOutputFetcher.AddPrevOut(*out.Outpoint, txOut)
			in := wire.NewTxIn(out.Outpoint, nil, nil)
			in.Sequence = defaultSequenceNum
			commitTx.AddTxIn(in)
			tmpPrivs[pos] = privs[i]
			pos++
			if i < count-1 { // the last uxto is fee item
				totalSenderAmount += btcutil.Amount(out.Value)
			}
			totalAmount += btcutil.Amount(out.Value)
		}
		time.Sleep(1 * time.Second) // limit rate
	}

	PkScript0, err := txscript.PayToAddrScript(receiverAddress)
	if err != nil {
		return nil, err
	}

	commitTx.AddTxOut(&wire.TxOut{
		PkScript: PkScript0,
		Value:    int64(totalSenderAmount),
	})
	changePkScript, err := txscript.PayToAddrScript(feeAddress)
	if err != nil {
		return nil, err
	}
	// make the change
	commitTx.AddTxOut(wire.NewTxOut(0, changePkScript))
	txsize := btcmempool.GetTxVirtualSize(btcutil.NewTx(commitTx))
	fee := btcutil.Amount(txsize) * btcutil.Amount(feerate)
	changeAmount := totalAmount - fee - totalSenderAmount

	if changeAmount > 0 {
		commitTx.TxOut[len(commitTx.TxOut)-1].Value = int64(changeAmount)
	} else {
		return nil, errors.New("not enough fees")
	}
	// make the signature
	witnessList := make([]wire.TxWitness, len(commitTx.TxIn))
	for i := range commitTx.TxIn {
		txOut := TxPrevOutputFetcher.FetchPrevOutput(commitTx.TxIn[i].PreviousOutPoint)
		witness, err := txscript.TaprootWitnessSignature(commitTx, txscript.NewTxSigHashes(commitTx, TxPrevOutputFetcher),
			i, txOut.Value, txOut.PkScript, txscript.SigHashDefault, tmpPrivs[i])
		if err != nil {
			return nil, err
		}
		witnessList[i] = witness
	}
	for i := range witnessList {
		commitTx.TxIn[i].Witness = witnessList[i]
	}
	return commitTx, nil
}

// make the collect tx
func makeCollectTx1(feerate int64, receiverAddress, feeAddress btcutil.Address, feePriv *btcec.PrivateKey,
	items []*OrderItem, btcApiClient *mempool.MempoolClient) (*wire.MsgTx, error) {

	privs, outlists, err := getUtxoFromOrders(items, btcApiClient)
	if err != nil {
		return nil, err
	}
	// get the fee_address utxo
	feeOutlist, err := gatherUTXO3(feeAddress, btcApiClient)
	if err != nil {
		return nil, err
	}
	privs = append(privs, feePriv)
	outlists = append(outlists, feeOutlist)

	tx, err := makeCollectTx0(feerate, receiverAddress, feeAddress, outlists, privs, btcApiClient)
	return tx, err
}

func getUtxoFromOrders(items []*OrderItem, btcApiClient *mempool.MempoolClient) ([]*btcec.PrivateKey, [][]*PrevOutPoint, error) {
	privs, outlists := make([]*btcec.PrivateKey, 0), make([][]*PrevOutPoint, 0)
	for _, item := range items {
		outlist, err := gatherUTXOForItem(item.Sender, btcApiClient)
		if err != nil && err != NoMoreUTXO {
			return nil, nil, err
		}
		privs = append(privs, item.Priv)
		outlists = append(outlists, outlist)
	}
	return privs, outlists, nil
}

func getOrders() ([]*OrderItem, error) {
	order := dao.Order{
		Action: dao.OrderActionToEVM,
		Stage:  dao.OrderStag2,
		Status: dao.OrderStatusConfirmed,
	}
	gotOrders, _, err := order.Find(nil, nil)
	if err != nil {
		return nil, err
	}

	orders := make([]*OrderItem, 0, len(gotOrders))
	for _, o := range gotOrders {
		amount, err := strconv.ParseInt(o.InAmount, 10, 64)
		if err != nil {
			params := map[string]interface{}{
				"order_id": o.ID,
				"amount":   o.InAmount,
				"error":    err,
			}
			log.Logger().WithFields(params).Error("failed to parse amount")
			return nil, err
		}
		orders = append(orders, &OrderItem{
			OrderID: o.ID,
			//Sender:  , // todo
			//Priv:    , // todo
			Amount: amount,
		})
	}

	return []*OrderItem{}, nil
}

func setOrders(ords []*OrderItem) error {
	ids := make([]uint64, 0, len(ords))
	for _, o := range ords {
		ids = append(ids, o.OrderID)
	}

	update := &dao.Order{
		Stage:  dao.OrderStag3,
		Status: dao.OrderStatusCollected,
	}
	return dao.NewOrder().UpdatesByIDs(ids, update)
}
func orderInfos(items []*OrderItem) (string, int64) {
	str, all := "", int64(0)
	for _, item := range items {
		all = all + item.Amount
		s := fmt.Sprintf("[orderID=%v,amount=%v]", item.OrderID, item.Amount)
		str = str + s + "\n"
	}
	return str, all
}
func checkFeeAddress(addr btcutil.Address, client *mempool.MempoolClient) (bool, error) {
	all := int64(0)
	outs, err := gatherUTXO3(addr, client)
	if err != nil {
		return false, err
	}
	for _, out := range outs {
		all = all + out.Value
	}
	if all < MinBalanceInFeeAddress {
		return false, nil
	}
	return true, nil
}
func getFeeRate(test bool, client *mempool.MempoolClient) int64 {
	if test {
		return 5
	}
	resp, err := client.RecommendedFees()
	if err != nil {
		return 50
	}
	return resp.FastestFee
}
func waitTxOnChain(txhash *chainhash.Hash, client *mempool.MempoolClient) (bool, error) {
	time.Sleep(30 * time.Second)
	fmt.Println("begin query....")
	for {
		resp, err := client.TransactionStatus(txhash)
		if err != nil {
			return false, err
		}
		if resp.Confirmed {
			return true, nil
		}
		fmt.Println("try query again....")
		time.Sleep(1 * time.Minute)
	}
}

// =============================================================================
func RunCollect(cfg *CollectCfg) error {
	privateKeyBytes, err := hex.DecodeString(cfg.StrFeePrivkey)
	if err != nil {
		panic(err)
	}
	feePriv, _ := btcec.PrivKeyFromBytes(privateKeyBytes)
	network := &chaincfg.MainNetParams
	if cfg.Testnet {
		network = &chaincfg.TestNet3Params
	}
	client := mempool.NewClient(network)

	for {
		// get the orders
		fmt.Println("begin collecting.....")
		ords, err := getOrders()
		if err != nil {
			log.Logger().WithField("error", err).Error("failed to get orders")
			alarm.Slack(context.Background(), "failed to get orders")
			return err // todo
		}
		strOrder, allAmount := orderInfos(ords)
		feerate := getFeeRate(cfg.Testnet, client)

		fmt.Println("orders:", strOrder)
		fmt.Println("all amount", allAmount, "feerate", feerate)

		enough, err := checkFeeAddress(cfg.FeeAddress, client)
		if err != nil {
			fmt.Println("checkFeeAddress failed", err)
			return err
		}
		if !enough {
			fmt.Println("the fee address has not enough balance")
			return errors.New("the fee address has not enough balance")
		}
		if len(ords) > 0 {
			tx, err := makeCollectTx1(feerate, cfg.Receiver, cfg.FeeAddress, feePriv, ords, client)
			if err != nil {
				//fmt.Println(err)
				log.Logger().WithField("error", err).Info("collect the order")
				alarm.Slack(context.Background(), "failed to make collect tx")
				return err // todo
			}
			txHash, err := client.BroadcastTx(tx)
			if err != nil {
				log.Logger().WithField("error", err).Error("failed to broadcast tx")
				alarm.Slack(context.Background(), "failed to broadcast tx")
				panic(err) // todo
			}
			//fmt.Println("collect the order...")
			//fmt.Println("collect the txhash", txHash.String())
			log.Logger().WithField("txhash", txHash.String()).Info("collect the order")
			onChain, err := waitTxOnChain(txHash, client)
			if err != nil {
				//fmt.Println("the collect tx on chain failed", err)
				log.Logger().WithField("error", err).Info("the collect tx on chain failed")
				alarm.Slack(context.Background(), "the collect tx on chain failed")
				return err
			}
			if onChain {
				err = setOrders(ords)
				if err != nil {
					fmt.Println("set order failed", err)
				}
			}
		}
		time.Sleep(30 * time.Minute)
	}
	return nil
}
