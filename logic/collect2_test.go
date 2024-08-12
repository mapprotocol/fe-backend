package logic

import (
	"errors"
	"github.com/btcsuite/btcd/btcutil"
	btcmempool "github.com/btcsuite/btcd/mempool"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/mapprotocol/fe-backend/third-party/mempool"
	"testing"
)

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

func makeSenderTx(feerate, sendAmount int64, sender, receiver, feeAddress btcutil.Address,
	outLists []*PrevOutPoint, btcApiClient *mempool.MempoolClient) (*wire.MsgTx, error) {

	commitTx := wire.NewMsgTx(wire.TxVersion)
	feeAmount, totalAmount := btcutil.Amount(0), btcutil.Amount(0)
	TxPrevOutputFetcher := txscript.NewMultiPrevOutFetcher(nil)

	// handle the every address's utxo
	for _, out := range outLists {
		txOut, err := getTxOutByOutPoint(out.Outpoint, btcApiClient)
		if err != nil {
			return nil, err
		}
		TxPrevOutputFetcher.AddPrevOut(*out.Outpoint, txOut)
		in := wire.NewTxIn(out.Outpoint, nil, nil)
		in.Sequence = defaultSequenceNum
		commitTx.AddTxIn(in)
		totalAmount += btcutil.Amount(out.Value)
	}
	// out0
	PkScript0, err := txscript.PayToAddrScript(receiver)
	if err != nil {
		return nil, err
	}
	commitTx.AddTxOut(&wire.TxOut{
		PkScript: PkScript0,
		Value:    int64(sendAmount),
	})
	// out1
	PkScript1, err := txscript.PayToAddrScript(feeAddress)
	if err != nil {
		return nil, err
	}
	commitTx.AddTxOut(&wire.TxOut{
		PkScript: PkScript1,
		Value:    int64(feeAmount),
	})

	changePkScript, err := txscript.PayToAddrScript(sender)
	if err != nil {
		return nil, err
	}
	// make the change
	commitTx.AddTxOut(wire.NewTxOut(0, changePkScript))
	txsize := btcmempool.GetTxVirtualSize(btcutil.NewTx(commitTx))
	fee := btcutil.Amount(txsize) * btcutil.Amount(feerate)
	changeAmount := totalAmount - fee - btcutil.Amount(sendAmount)

	if changeAmount > 0 {
		commitTx.TxOut[len(commitTx.TxOut)-1].Value = int64(changeAmount)
	} else {
		return nil, errors.New("not enough fees")
	}
	// make the signature
	//witnessList := make([]wire.TxWitness, len(commitTx.TxIn))
	//for i := range commitTx.TxIn {
	//	txOut := TxPrevOutputFetcher.FetchPrevOutput(commitTx.TxIn[i].PreviousOutPoint)
	//	witness, err := txscript.TaprootWitnessSignature(commitTx, txscript.NewTxSigHashes(commitTx, TxPrevOutputFetcher),
	//		i, txOut.Value, txOut.PkScript, txscript.SigHashDefault, tmpPrivs[i])
	//	if err != nil {
	//		return nil, err
	//	}
	//	witnessList[i] = witness
	//}
	//for i := range witnessList {
	//	commitTx.TxIn[i].Witness = witnessList[i]
	//}
	return commitTx, nil
}

func Test_makeSenderTx(t *testing.T) {
	//proxyAddressStr := ""
	//senderStr := ""
	//rate := 100 / 10000
	//network := &chaincfg.MainNetParams
	//if testnet {
	//	network = &chaincfg.TestNet3Params
	//}
	//proxyAddress, _ := btcutil.DecodeAddress(proxyAddressStr, network)
	//sender, _ := btcutil.DecodeAddress(proxyAddressStr, network)

}
