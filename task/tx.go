package task

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	btcmempool "github.com/btcsuite/btcd/mempool"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/mapprotocol/fe-backend/third-party/mempool"
	"testing"
)

const (
	defaultSequenceNum = wire.MaxTxInSequenceNum - 10
)

var (
	PrevAdminOutPoint2        *PrevOutPoint = nil
	MinPreAdminOutPointValue2               = int64(20000)
)

type PrevOutPoint struct {
	Outpoint *wire.OutPoint
	Value    int64
}

func makeTpAddress(privKey *btcec.PrivateKey, network *chaincfg.Params) (btcutil.Address, error) {
	tapKey := txscript.ComputeTaprootKeyNoScript(privKey.PubKey())

	address, err := btcutil.NewAddressTaproot(
		schnorr.SerializePubKey(tapKey),
		network,
	)
	if err != nil {
		return nil, err
	}
	return address, nil
}
func setPrevOutPoint(outpoint *wire.OutPoint, val int64) {
	tmp := &PrevOutPoint{
		Outpoint: outpoint,
		Value:    val,
	}
	PrevAdminOutPoint2 = tmp
}

func gatherUtxo(client *mempool.MempoolClient, sender btcutil.Address) ([]*PrevOutPoint, error) {
	outPointList := make([]*PrevOutPoint, 0)

	if PrevAdminOutPoint2 != nil {
		if PrevAdminOutPoint2.Value > MinPreAdminOutPointValue2 {
			outPointList = append(outPointList, PrevAdminOutPoint2)
			return outPointList, nil
		}
	}

	unspentList, err := client.ListUnspent(sender)
	if err != nil {
		return nil, err
	}

	if len(unspentList) == 0 {
		err = fmt.Errorf("no utxo for %s", sender)
		return nil, err
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

func getTxOutByOutPoint(outPoint *wire.OutPoint, btcClient *mempool.MempoolClient) (*wire.TxOut, error) {
	tx, err := btcClient.GetRawTransaction(&outPoint.Hash)
	if err != nil {
		return nil, err
	}
	if int(outPoint.Index) >= len(tx.TxOut) {
		return nil, errors.New("err out point")
	}
	return tx.TxOut[outPoint.Index], nil
}

func makeTransaction(client *mempool.MempoolClient, feeRate, amount int64, sender, receiver btcutil.Address, outList []*PrevOutPoint,
	senderPriv *btcec.PrivateKey) (*wire.MsgTx, error) {

	commitTx := wire.NewMsgTx(wire.TxVersion)
	totalSenderAmount := btcutil.Amount(0)
	TxPrevOutputFetcher := txscript.NewMultiPrevOutFetcher(nil)

	for _, out := range outList {
		txOut, err := getTxOutByOutPoint(out.Outpoint, client)
		if err != nil {
			return nil, err
		}
		TxPrevOutputFetcher.AddPrevOut(*out.Outpoint, txOut)
		in := wire.NewTxIn(out.Outpoint, nil, nil)
		in.Sequence = defaultSequenceNum
		commitTx.AddTxIn(in)

		totalSenderAmount += btcutil.Amount(out.Value)
	}

	PkScript0, err := txscript.PayToAddrScript(receiver)
	if err != nil {
		return nil, err
	}

	commitTx.AddTxOut(&wire.TxOut{
		PkScript: PkScript0,
		Value:    amount,
	})
	changePkScript, err := txscript.PayToAddrScript(sender)
	if err != nil {
		return nil, err
	}
	// make the change
	commitTx.AddTxOut(wire.NewTxOut(0, changePkScript))
	txsize := btcmempool.GetTxVirtualSize(btcutil.NewTx(commitTx))
	fee := btcutil.Amount(txsize) * btcutil.Amount(feeRate)
	changeAmount := totalSenderAmount - fee - btcutil.Amount(amount)

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
			i, txOut.Value, txOut.PkScript, txscript.SigHashDefault, senderPriv)
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

func SendTransaction(client *mempool.MempoolClient, feeRate, amount int64, sender, receiver btcutil.Address, senderPrivateKey *btcec.PrivateKey) (*chainhash.Hash, error) {
	outPointList, err := gatherUtxo(client, sender)
	if err != nil {
		return nil, err
	}
	commitTx, err := makeTransaction(client, feeRate, amount, sender, receiver, outPointList, senderPrivateKey)
	if err != nil {
		return nil, err
	}

	txHash, err := client.BroadcastTx(commitTx)
	if err != nil {
		return nil, err
	}

	if len(commitTx.TxOut) > 0 {
		vout := len(commitTx.TxOut) - 1
		setPrevOutPoint(&wire.OutPoint{
			Hash:  *txHash,
			Index: uint32(vout),
		}, commitTx.TxOut[vout].Value)
	}
	return txHash, nil
}

func Test_makeTx(t *testing.T) {
	network := &chaincfg.MainNetParams
	client := mempool.NewClient(network)
	privateKeyBytes, err := hex.DecodeString("")
	if err != nil {
		panic(err)
	}
	senderKey, _ := btcec.PrivKeyFromBytes(privateKeyBytes)
	sender, err := makeTpAddress(senderKey, network)
	if err != nil {
		panic(err)
	}
	sendAmount, feerate := int64(10000), int64(50)
	OutPointList, err := gatherUtxo(client, sender)
	if err != nil {
		panic(err)
	}

	commitTx, err := makeTransaction(client, feerate, sendAmount, sender, sender, OutPointList, senderKey)
	if err != nil {
		panic(err)
	}

	txHash, err := client.BroadcastTx(commitTx)
	if err != nil {
		panic(err)
	}

	if len(commitTx.TxOut) > 0 {
		vout := len(commitTx.TxOut) - 1
		setPrevOutPoint(&wire.OutPoint{
			Hash:  *txHash,
			Index: uint32(vout),
		}, commitTx.TxOut[vout].Value)
	}
}
