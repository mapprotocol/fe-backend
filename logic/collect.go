package logic

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	localErr "github.com/mapprotocol/fe-backend/resource/err"
	"github.com/mapprotocol/fe-backend/utils"
	"strconv"
	"sync"
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
	CollectDoing       = 1
	CollectFinish      = 2
)

var (
	//PrevAdminOutPoint2        *PrevOutPoint = nil
	//MinPreAdminOutPointValue2               = int64(20000)
	MinBalanceInFeeAddress = int64(20000)
	NoMoreUTXO             = errors.New("no more utxo")
	MinUtxoAmount          = int64(100)
	LowBalanceHotWallet    = 11
	FullBalanceHotWallet   = 12
)

type PrevOutPoint struct {
	Outpoint *wire.OutPoint
	Value    int64
}
type CollectCfg struct {
	Testnet                 bool
	StrHotWalletFee1Privkey string
	StrHotWallet1Priv       string
	HotWalletFee1           btcutil.Address
	HotWallet1              btcutil.Address

	StrHotWalletFee2Privkey string
	StrHotWallet2Priv       string
	HotWalletFee2           btcutil.Address
	HotWallet2              btcutil.Address
	HotWallet2Line          int64

	StrHotWalletFee3Privkey string
	HotWalletFee3           btcutil.Address

	MaxTransferAmount int64
}

type OrderItem struct {
	OrderID uint64
	Sender  btcutil.Address
	Priv    *btcec.PrivateKey
	Amount  int64
}

type WithdrawOrder struct {
	OrderID  uint64
	Receiver btcutil.Address
	Amount   int64
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
		if unspentList[i].Output.Value < MinUtxoAmount {
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

func getOrders(limit int, network *chaincfg.Params) ([]*OrderItem, int64, error) {
	order := dao.Order{
		Action: dao.OrderActionToEVM,
		Stage:  dao.OrderStag2,
		Status: dao.OrderStatusConfirmed,
	}
	gotOrders, count, err := order.Find(nil, dao.Paginate(1, limit))
	if err != nil {
		return nil, 0, err
	}

	orders := make([]*OrderItem, 0, len(gotOrders))
	for _, o := range gotOrders {
		relayer, err := btcutil.DecodeAddress(o.Relayer, network)
		if err != nil {
			params := map[string]interface{}{
				"order_id": o.ID,
				"network":  network.Net.String(),
				"relayer":  o.Relayer,
				"error":    err,
			}
			log.Logger().WithFields(params).Error("decode relayer address failed")
			return nil, 0, err
		}

		privateKeyBytes, err := hex.DecodeString(o.RelayerPrivateKey)
		if err != nil {
			params := map[string]interface{}{
				"order_id": o.ID,
				"error":    err,
			}
			log.Logger().WithFields(params).Error("failed to decode private key")
			return nil, 0, err
		}
		privakeKey, _ := btcec.PrivKeyFromBytes(privateKeyBytes)

		amount, err := strconv.ParseInt(o.InAmount, 10, 64)
		if err != nil {
			params := map[string]interface{}{
				"order_id": o.ID,
				"amount":   o.InAmount,
				"error":    err,
			}
			log.Logger().WithFields(params).Error("failed to parse amount")
			return nil, 0, err
		}

		orders = append(orders, &OrderItem{
			OrderID: o.ID,
			Sender:  relayer,
			Priv:    privakeKey,
			Amount:  amount,
		})
	}

	return orders, count, nil
}

func getLatestCollectInfo() (*chainhash.Hash, []*OrderItem, error) {
	collect := &dao.Collect{
		Status: dao.CollectStatusPending,
	}
	gotCollects, _, err := collect.Find(nil, nil)
	if err != nil {
		return nil, nil, err
	}
	if len(gotCollects) == 0 {
		return nil, []*OrderItem{}, nil
	}

	orders := make([]*OrderItem, 0, len(gotCollects))
	for _, c := range gotCollects {
		orders = append(orders, &OrderItem{
			OrderID: c.OrderID,
		})
	}
	txHash, err := chainhash.NewHashFromStr(gotCollects[0].TxHash)
	if err != nil {
		return nil, nil, err
	}
	return txHash, orders, nil
}
func createLatestCollectInfo(txhash *chainhash.Hash, orders []*OrderItem) error {
	txHash := txhash.String()
	collects := make([]*dao.Collect, 0, len(orders))
	for _, o := range orders {
		collects = append(collects, &dao.Collect{
			OrderID: o.OrderID,
			TxHash:  txHash,
		})
	}

	if err := dao.NewCollect().BatchCreate(collects); err != nil {
		params := map[string]interface{}{
			"collects": utils.JSON(collects),
			"error":    err,
		}
		log.Logger().WithFields(params).Error("failed to create collects")
		return err
	}
	return nil
}

func setLatestCollectInfo(txhash *chainhash.Hash) error {
	collect := &dao.Collect{
		TxHash: txhash.String(),
	}
	update := &dao.Collect{
		Status: dao.CollectStatusConfirmed,
	}
	if err := collect.Updates(update); err != nil {
		params := map[string]interface{}{
			"tx_hash": txhash.String(),
			"status":  dao.CollectStatusConfirmed,
			"error":   err,
		}
		log.Logger().WithFields(params).Error("failed to update collect status")
		return err
	}
	return nil
}

func setOrders(ords []*OrderItem, status uint8) error {
	ids := make([]uint64, 0, len(ords))
	for _, o := range ords {
		ids = append(ids, o.OrderID)
	}

	update := &dao.Order{
		Stage:  dao.OrderStag3,
		Status: status,
	}
	if err := dao.NewOrder().UpdatesByIDs(ids, update); err != nil {
		params := map[string]interface{}{
			"ids":    utils.JSON(ids),
			"update": utils.JSON(update),
			"error":  err,
		}
		log.Logger().WithFields(params).Error("failed to update order status")
		return err
	}
	return nil
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
		return 20
	}
	resp, err := client.RecommendedFees()
	if err != nil {
		return 50
	}
	return resp.FastestFee
}
func waitTxOnChain(txhash *chainhash.Hash, client *mempool.MempoolClient) error {
	time.Sleep(30 * time.Second)
	fmt.Println("begin query....")
	for {
		resp, err := client.TransactionStatus(txhash)
		if err != nil {
			return err
		}
		if resp.Confirmed {
			return nil
		}
		fmt.Println("try query again....")
		time.Sleep(1 * time.Minute)
	}
}
func checkLatestTx(client *mempool.MempoolClient) error {
	txhash, itmes, err := getLatestCollectInfo()
	if err != nil {
		log.Logger().WithField("error", err).Error("get latest collect info failed")
		return err
	}
	if txhash == nil {
		return nil
	}
	err = waitTxOnChain(txhash, client)
	if err != nil {
		log.Logger().WithField("error", err).Error("wait tx on chain failed")
		return err
	} else {
		if err = setOrders(itmes, CollectFinish); err == nil {
			err = setLatestCollectInfo(txhash)
			if err != nil {
				log.Logger().WithField("error", err).Error("set latest collect info failed in check process")
			}
		} else {
			log.Logger().WithField("error", err).Error("setOrders finish failed in check process")
		}
	}
	return nil
}

// =============================================================================
// withdraw infos
func withdrawOrderToIds(items []*WithdrawOrder) []uint64 {
	ids := make([]uint64, 0)
	for _, item := range items {
		ids = append(ids, item.OrderID)
	}
	return ids
}
func withdrawOrdersInfos(items []*WithdrawOrder) string {
	str, all := "", int64(0)
	for _, item := range items {
		str += fmt.Sprintf("[id=%v,amount=%v]", item.OrderID, item.Amount)
		all += item.Amount
	}
	str0 := fmt.Sprintf("[ids=%d,all=%v] {", len(items), all)
	str0 = str0 + str + "}"
	return str0
}

func getWithdrawOrders(limit int, network *chaincfg.Params) ([]*WithdrawOrder, error) {
	order := dao.Order{
		Action: dao.OrderActionFromEVM,
		Stage:  dao.OrderStag2,
		Status: dao.OrderStatusConfirmed,
	}
	gotOrders, _, err := order.Find(nil, dao.Paginate(1, limit))
	if err != nil {
		return nil, err
	}

	orders := make([]*WithdrawOrder, 0, len(gotOrders))
	for _, o := range gotOrders {
		receiver, err := btcutil.DecodeAddress(o.Receiver, network)
		if err != nil {
			params := map[string]interface{}{
				"order_id": o.ID,
				"network":  network.Net.String(),
				"relayer":  o.Relayer,
				"error":    err,
			}
			log.Logger().WithFields(params).Error("decode relayer address failed")
			return nil, err
		}

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

		orders = append(orders, &WithdrawOrder{
			OrderID:  o.ID,
			Receiver: receiver,
			Amount:   amount,
		})
	}

	return orders, nil
}

// state = 1 | 2
func getInitedWithdrawOrders(state uint8, limit int) ([]*chainhash.Hash, error) {
	order := dao.Order{
		Action: dao.OrderActionFromEVM,
		Stage:  dao.OrderStag3,
		Status: state,
	}
	gotOrders, _, err := order.Find(nil, dao.Paginate(1, limit))
	if err != nil {
		return nil, err
	}

	hashes := make([]*chainhash.Hash, 0, len(gotOrders))
	for _, o := range gotOrders {
		txHash, err := chainhash.NewHashFromStr(o.OutTxHash)
		if err != nil {
			return nil, fmt.Errorf("failed to parse tx hash, %s", o.OutTxHash)
		}

		hashes = append(hashes, txHash)
	}
	return nil, nil
}

// state = 1 & txhash
func initWithdrawOrders(txhash *chainhash.Hash, ids []uint64) error {
	update := &dao.Order{
		Stage:     dao.OrderStag3,
		Status:    dao.WithdrawStateInit,
		OutTxHash: txhash.String(),
	}
	if err := dao.NewOrder().UpdatesByIDs(ids, update); err != nil {
		params := map[string]interface{}{
			"ids":    utils.JSON(ids),
			"update": utils.JSON(update),
			"error":  err,
		}
		log.Logger().WithFields(params).Error("failed to update order status")
		return err
	}
	return nil
}

// state=2
func updateWithdrawOrdersState(txhash *chainhash.Hash, state uint8) error {
	order := dao.Order{
		OutTxHash: txhash.String(),
	}
	update := &dao.Order{
		Status: state,
	}
	if err := order.Updates(update); err != nil {
		params := map[string]interface{}{
			"txHash": txhash.String(),
			"update": utils.JSON(update),
			"error":  err,
		}
		log.Logger().WithFields(params).Error("failed to update order status")
		return err
	}
	return nil
}

/*
tx_in_sender1 						tx_out_receiver
tx_in_sender2...       	--- >		tx_out_change (sender)
tx_in_fee1				--- >   	tx_out_change (fee)
tx_in_fee2...
*/
func makeWithdrawTx0(feerate int64, tipper, sender btcutil.Address, senderPriv, feePriv *btcec.PrivateKey, senderOutList,
	feeOutList []*PrevOutPoint, items []*WithdrawOrder, btcApiClient *mempool.MempoolClient) (*wire.MsgTx, error) {

	commitTx := wire.NewMsgTx(wire.TxVersion)
	totalSenderAmount, totalAmount := btcutil.Amount(0), btcutil.Amount(0)
	TxPrevOutputFetcher := txscript.NewMultiPrevOutFetcher(nil)
	pos := 0
	tmpPrivs := make(map[int]*btcec.PrivateKey)

	// handle the sender's utxo
	for _, out := range senderOutList {
		txOut, err := getTxOutByOutPoint2(out.Outpoint, btcApiClient)
		if err != nil {
			return nil, err
		}
		TxPrevOutputFetcher.AddPrevOut(*out.Outpoint, txOut)
		in := wire.NewTxIn(out.Outpoint, nil, nil)
		in.Sequence = defaultSequenceNum
		commitTx.AddTxIn(in)
		tmpPrivs[pos] = senderPriv
		pos++
		totalSenderAmount += btcutil.Amount(out.Value)
		totalAmount += btcutil.Amount(out.Value)
	}
	time.Sleep(1 * time.Second) // limit rate
	// handle the fee's utxo
	for _, out := range feeOutList {
		txOut, err := getTxOutByOutPoint2(out.Outpoint, btcApiClient)
		if err != nil {
			return nil, err
		}
		TxPrevOutputFetcher.AddPrevOut(*out.Outpoint, txOut)
		in := wire.NewTxIn(out.Outpoint, nil, nil)
		in.Sequence = defaultSequenceNum
		commitTx.AddTxIn(in)
		tmpPrivs[pos] = feePriv
		pos++
		totalAmount += btcutil.Amount(out.Value)
	}

	// handle the tx output
	outAmount := int64(0)
	for _, item := range items {
		PkScript0, err := txscript.PayToAddrScript(item.Receiver)
		if err != nil {
			return nil, err
		}
		commitTx.AddTxOut(&wire.TxOut{
			PkScript: PkScript0,
			Value:    item.Amount,
		})
		outAmount += item.Amount
	}
	if int64(totalSenderAmount) < outAmount {
		log.Logger().WithField("hotwallet2", totalSenderAmount).WithField("need", outAmount).
			Error("low balance")
		alarm.Slack(context.Background(), fmt.Sprintf("[hot-wallet2=%v,need=%v]:low balance", totalAmount, outAmount))
		return nil, localErr.LowBalanceInHotWallet2
	}

	PkScript1, err := txscript.PayToAddrScript(sender)
	if err != nil {
		return nil, err
	}
	commitTx.AddTxOut(&wire.TxOut{
		PkScript: PkScript1,
		Value:    int64(totalSenderAmount) - outAmount,
	})
	changePkScript, err := txscript.PayToAddrScript(tipper)
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
		log.Logger().Error(localErr.LowFeeInHotWalletFee2)
		alarm.Slack(context.Background(), localErr.LowFeeInHotWalletFee2.Error())

		return nil, localErr.LowFeeInHotWalletFee2
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

func makeWithdrawTx1(feerate int64, tipper, sender btcutil.Address, senderPriv, feePriv *btcec.PrivateKey,
	items []*WithdrawOrder, btcApiClient *mempool.MempoolClient) (*wire.MsgTx, error) {

	// get the sender_address utxo
	senderOutlist, err := gatherUTXO3(sender, btcApiClient)
	if err != nil {
		return nil, err
	}
	// get the fee_address utxo
	feeOutlist, err := gatherUTXO3(tipper, btcApiClient)
	if err != nil {
		return nil, err
	}

	tx, err := makeWithdrawTx0(feerate, tipper, sender, senderPriv, feePriv,
		senderOutlist, feeOutlist, items, btcApiClient)

	return tx, err
}

func checkWithdrawTxsState(cfg *CollectCfg) {
	network := &chaincfg.MainNetParams
	if cfg.Testnet {
		network = &chaincfg.TestNet3Params
	}
	client := mempool.NewClient(network)

	for {
		log.Logger().Info("begin withdraw tx state check...")
		hashs, err := getInitedWithdrawOrders(dao.WithdrawStateSend, 20) // todo update limit
		if err != nil {
			log.Logger().WithField("error", err).Error("getInitedWithdrawOrders in check state failed")
		} else {
			for _, h := range hashs {
				err = waitTxOnChain(h, client)
				if err != nil {
					log.Logger().WithField("error", err).WithField("hash", h.String()).
						Error("wait on chain failed [check state]")
				} else {
					log.Logger().WithField("hash", h.String()).Info("the hash was on chain")
					err = updateWithdrawOrdersState(h, dao.WithdrawStateFinish)
					if err != nil {
						log.Logger().WithField("hash", h.String()).WithField("error", err).Error("update the hash to finish state failed")
					}
				}
			}
		}
		time.Sleep(5 * time.Minute)
	}
}
func checkHotwallet2Balance(cfg *CollectCfg, hotwallet2 btcutil.Address, client *mempool.MempoolClient) (bool, error) {
	outs, err := gatherUTXO3(hotwallet2, client)
	if err != nil {
		return false, err
	}
	all := int64(0)
	for _, out := range outs {
		all += out.Value
	}
	if all < cfg.HotWallet2Line {
		return true, nil
	}
	return false, nil
}

// =============================================================================
func makeSimpleTx0(feerate, amount int64, sender, receiver, tipper btcutil.Address, senderPriv,
	feePriv *btcec.PrivateKey, senderOutList, feeOutList []*PrevOutPoint, btcApiClient *mempool.MempoolClient) (*wire.MsgTx, error) {

	commitTx := wire.NewMsgTx(wire.TxVersion)
	totalSenderAmount, totalAmount := btcutil.Amount(0), btcutil.Amount(0)
	TxPrevOutputFetcher := txscript.NewMultiPrevOutFetcher(nil)
	tmpPrivs, pos := make(map[int]*btcec.PrivateKey), 0

	// handle the sender's utxo
	for _, out := range senderOutList {
		txOut, err := getTxOutByOutPoint2(out.Outpoint, btcApiClient)
		if err != nil {
			return nil, err
		}
		TxPrevOutputFetcher.AddPrevOut(*out.Outpoint, txOut)
		in := wire.NewTxIn(out.Outpoint, nil, nil)
		in.Sequence = defaultSequenceNum
		commitTx.AddTxIn(in)
		tmpPrivs[pos] = senderPriv
		pos++
		totalSenderAmount += btcutil.Amount(out.Value)
		totalAmount += btcutil.Amount(out.Value)
	}
	time.Sleep(1 * time.Second) // limit rate
	// handle the fee's utxo
	for _, out := range feeOutList {
		txOut, err := getTxOutByOutPoint2(out.Outpoint, btcApiClient)
		if err != nil {
			return nil, err
		}
		TxPrevOutputFetcher.AddPrevOut(*out.Outpoint, txOut)
		in := wire.NewTxIn(out.Outpoint, nil, nil)
		in.Sequence = defaultSequenceNum
		commitTx.AddTxIn(in)
		tmpPrivs[pos] = feePriv
		pos++
		totalAmount += btcutil.Amount(out.Value)
	}

	// handle the tx output
	PkScript0, err := txscript.PayToAddrScript(receiver)
	if err != nil {
		return nil, err
	}
	commitTx.AddTxOut(&wire.TxOut{
		PkScript: PkScript0,
		Value:    amount,
	})
	if int64(totalSenderAmount) < amount {
		return nil, localErr.LowBalanceInHotWallet1
	}
	// handle the sender change
	PkScript1, err := txscript.PayToAddrScript(sender)
	if err != nil {
		return nil, err
	}
	commitTx.AddTxOut(&wire.TxOut{
		PkScript: PkScript1,
		Value:    int64(totalSenderAmount) - amount,
	})
	changePkScript, err := txscript.PayToAddrScript(tipper)
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
		return nil, localErr.LowFeeInHotWalletFee3
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

func makeHotWallet1ToHotWallet2Tx(feerate int64, cfg *CollectCfg, btcApiClient *mempool.MempoolClient) (*wire.MsgTx, error) {
	privateKeyBytes, err := hex.DecodeString(cfg.StrHotWallet1Priv)
	if err != nil {
		panic(err)
	}
	senderPriv, _ := btcec.PrivKeyFromBytes(privateKeyBytes)

	privateKeyBytes, err = hex.DecodeString(cfg.StrHotWalletFee3Privkey)
	if err != nil {
		panic(err)
	}
	feePriv, _ := btcec.PrivKeyFromBytes(privateKeyBytes)

	// get the sender_address utxo
	senderOutlist, err := gatherUTXO3(cfg.HotWallet1, btcApiClient)
	if err != nil {
		return nil, err
	}
	// get the fee_address utxo
	feeOutlist, err := gatherUTXO3(cfg.HotWalletFee3, btcApiClient)
	if err != nil {
		return nil, err
	}

	amount := cfg.MaxTransferAmount
	tx, err := makeSimpleTx0(feerate, amount, cfg.HotWallet1, cfg.HotWallet2, cfg.HotWalletFee3, senderPriv, feePriv,
		senderOutlist, feeOutlist, btcApiClient)

	return tx, err
}

// =============================================================================
func Run(cfg *CollectCfg) error {
	wg := &sync.WaitGroup{}
	wg.Add(4)
	chBalanceLow, chBalanceHigh := make(chan int), make(chan int)

	go func(cfg *CollectCfg) {
		defer wg.Done()
		err := RunCollect(cfg)
		log.Logger().WithField("error", err).Error("collect process finish")
	}(cfg)

	go func(cfg *CollectCfg) {
		defer wg.Done()
		err := RunBtcWithdraw(cfg, chBalanceLow, chBalanceHigh)
		log.Logger().WithField("error", err).Error("withdraw process finish")
	}(cfg)

	go func(cfg *CollectCfg) {
		defer wg.Done()
		err := RunHotWalletBalance(cfg, chBalanceLow, chBalanceHigh)
		log.Logger().WithField("error", err).Error("hot wallet balance process finish")
	}(cfg)

	go func(cfg *CollectCfg) {
		defer wg.Done()
		checkWithdrawTxsState(cfg)
	}(cfg)

	wg.Wait()
	close(chBalanceLow)
	close(chBalanceHigh)
	log.Logger().Info("......finish......")
	return nil
}
func RunCollect(cfg *CollectCfg) error {
	privateKeyBytes, err := hex.DecodeString(cfg.StrHotWalletFee1Privkey)
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
		fmt.Println("begin collecting, check the latest collect info.....")
		err := checkLatestTx(client)
		if err != nil {
			log.Logger().Info("check latest failed... will be retry")
			time.Sleep(3 * time.Minute)
			continue
		}
		ords, _, err := getOrders(10, network)
		if err != nil {
			log.Logger().WithField("error", err).Error("failed to get orders")
			alarm.Slack(context.Background(), "failed to get orders")
			return err
		}
		strOrder, allAmount := orderInfos(ords)
		feerate := getFeeRate(cfg.Testnet, client)

		log.Logger().WithField("orders", strOrder).Info("collect the order")
		log.Logger().WithField("all amount", allAmount).WithField("feerate", feerate).Info("collect the order")

		enough, err := checkFeeAddress(cfg.HotWalletFee1, client)
		if err != nil {
			log.Logger().WithField("error", err).Error("failed to checkFeeAddress")
			alarm.Slack(context.Background(), "check fee address balance failed")
			return err
		}
		if !enough {
			e := errors.New("the fee address has not enough balance")
			log.Logger().WithField("error", e).Error("low balance in fee address")
			alarm.Slack(context.Background(), "low balance in fee address")
		}

		if len(ords) > 0 {
			tx, err := makeCollectTx1(feerate, cfg.HotWallet1, cfg.HotWalletFee1, feePriv, ords, client)
			if err != nil {
				//fmt.Println(err)
				log.Logger().WithField("error", err).Info("make collect tx")
				alarm.Slack(context.Background(), "failed to make collect tx")
				return err
			}
			txHash, err := client.BroadcastTx(tx)
			if err != nil {
				log.Logger().WithField("error", err).Error("failed to broadcast tx")
				alarm.Slack(context.Background(), "failed to broadcast tx")
				return err
			}
			log.Logger().WithField("txhash", txHash.String()).Info("broadcast the collect tx")
			err = createLatestCollectInfo(txHash, ords)
			if err != nil {
				log.Logger().WithField("error", err).Error("create latest collect info failed")
				return err
			}
			log.Logger().Info("create latest collect info success")
			err = waitTxOnChain(txHash, client)
			if err != nil {
				//fmt.Println("the collect tx on chain failed", err)
				log.Logger().WithField("error", err).Info("the collect tx on chain failed")
				alarm.Slack(context.Background(), "the collect tx on chain failed")
				return err
			} else {
				err = setOrders(ords, CollectFinish)
				if err != nil {
					log.Logger().WithField("error", err).Info("set orders state failed")
				}
				setLatestCollectInfo(txHash)
			}
		}
		time.Sleep(30 * time.Minute)
	}
	return nil
}

func btcWithdrawTxTransfer(cfg *CollectCfg, tipper, sender btcutil.Address, tipperPriv, senderPriv *btcec.PrivateKey,
	orders []*WithdrawOrder, client *mempool.MempoolClient) error {

	log.Logger().Info("begin withdrawal tx process.....")
	feerate := getFeeRate(cfg.Testnet, client)

	if len(orders) > 0 {
		log.Logger().WithField("key", withdrawOrdersInfos(orders)).Info("user withdraw...")

		tx, err := makeWithdrawTx1(feerate, tipper, sender, senderPriv, tipperPriv, orders, client)
		if err != nil {
			log.Logger().WithField("error", err).Info("make withdraw tx failed")
			alarm.Slack(context.Background(), "failed to make withdraw tx")
			return err
		}
		// 1 init the withdraw state
		txhash1, ids := tx.TxHash(), withdrawOrderToIds(orders)
		log.Logger().WithField("txhash", txhash1.String()).Info("init user withdraw order")

		err = initWithdrawOrders(&txhash1, ids)
		if err != nil {
			log.Logger().WithField("error", err).WithField("txhash", txhash1.String()).
				Error("init user withdraw order failed")
			return err
		} else {
			txHash, err := client.BroadcastTx(tx)
			if err != nil {
				log.Logger().WithField("error", err).Error("failed to broadcast tx")
				alarm.Slack(context.Background(), "failed to broadcast tx")
				return nil
			}
			log.Logger().WithField("txhash", txHash.String()).Info("broadcast the withdraw tx")
			//  2. update orders state to WithdrawStateSend
			err = updateWithdrawOrdersState(&txhash1, dao.WithdrawStateSend)
			if err != nil {
				log.Logger().WithField("error", err).WithField("setstate", dao.WithdrawStateSend).Error("update state failed")
				return err
			}
		}
	}
	return nil
}
func RunBtcWithdraw(cfg *CollectCfg, ch1, ch2 chan int) error {
	privateKeyBytes, err := hex.DecodeString(cfg.StrHotWalletFee2Privkey)
	if err != nil {
		panic(err)
	}
	tipperPriv, _ := btcec.PrivKeyFromBytes(privateKeyBytes)
	privBytes1, err := hex.DecodeString(cfg.StrHotWallet2Priv)
	if err != nil {
		panic(err)
	}
	senderPriv, _ := btcec.PrivKeyFromBytes(privBytes1)
	tipper, sender := cfg.HotWalletFee2, cfg.HotWallet2

	network := &chaincfg.MainNetParams
	if cfg.Testnet {
		network = &chaincfg.TestNet3Params
	}
	client := mempool.NewClient(network)

	stopWithdraw := false
	ticker := time.NewTicker(2 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case msg := <-ch2:
			if msg == FullBalanceHotWallet {
				stopWithdraw = false
			}
		case <-ticker.C:
			if !stopWithdraw {
				orders, err := getWithdrawOrders(20, network)
				if err != nil {
					log.Logger().WithField("error", err).Error("failed to get orders")
					alarm.Slack(context.Background(), "failed to get withdraw orders")
				} else {
					err = btcWithdrawTxTransfer(cfg, tipper, sender, tipperPriv, senderPriv, orders, client)
					if err != nil {
						if err == localErr.LowBalanceInHotWallet2 {
							stopWithdraw = true
							ch1 <- LowBalanceHotWallet
						}
					}
				}
			}
		default:
		}
	}
	return nil
}

func HotWalletBalanceTransfer(cfg *CollectCfg, client *mempool.MempoolClient) error {
	// get the withdrawal order
	log.Logger().Info("begin hotwallet1 to hotwallet2 process.....")
	feerate := getFeeRate(cfg.Testnet, client)

	tx, err := makeHotWallet1ToHotWallet2Tx(feerate, cfg, client)
	if err != nil {
		log.Logger().WithField("error", err).Info("wallet1 to wallet2 failed")
		return err
	}

	txHash, err := client.BroadcastTx(tx)
	if err != nil {
		log.Logger().WithField("error", err).Error("failed to broadcast tx")
		return err
	}
	log.Logger().WithField("txhash", txHash.String()).Info("broadcast the wallet1 to wallet2 tx")
	err = waitTxOnChain(txHash, client)
	if err != nil {
		log.Logger().WithField("error", err).Error("wait the tx on chain failed")
	} else {
		log.Logger().WithField("txhash", txHash.String()).Info("the tx was on chain")
	}
	return err
}

func RunHotWalletBalance(cfg *CollectCfg, ch1, ch2 chan int) error {
	network := &chaincfg.MainNetParams
	if cfg.Testnet {
		network = &chaincfg.TestNet3Params
	}
	client := mempool.NewClient(network)

	ticker := time.NewTicker(30 * time.Minute)
	defer ticker.Stop()
	transfer, hotwallet2 := false, cfg.HotWallet2

	for {
		select {
		case <-ticker.C:
			// check the water line
			t, err := checkHotwallet2Balance(cfg, hotwallet2, client)
			if transfer || (t && err == nil) {
				err := HotWalletBalanceTransfer(cfg, client)
				if err != nil {
					log.Logger().WithField("error", err).Info("wallet1 to wallet2 failed")
					alarm.Slack(context.Background(), "failed to makeHotWallet1ToHotWallet2Tx")
				}
			}
		case msg := <-ch1:
			if msg == LowBalanceHotWallet {
				transfer = true
				err := HotWalletBalanceTransfer(cfg, client)
				if err != nil {
					log.Logger().WithField("error", err).Info("wallet1 to wallet2 failed")
					alarm.Slack(context.Background(), "failed to makeHotWallet1ToHotWallet2Tx")
				} else {
					ch2 <- FullBalanceHotWallet
					transfer = false
				}
			}
		default:
		}
	}
	return nil
}
