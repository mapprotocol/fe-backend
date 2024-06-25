package logic

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	btcmempool "github.com/btcsuite/btcd/mempool"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/mapprotocol/fe-backend/third-party/mempool"
	"testing"
)

var (
	testnet = false
	priv1   = ""
)

func makeTpAddress(privKey *btcec.PrivateKey, netParams *chaincfg.Params) (btcutil.Address, error) {
	tapKey := txscript.ComputeTaprootKeyNoScript(privKey.PubKey())

	address, err := btcutil.NewAddressTaproot(
		schnorr.SerializePubKey(tapKey),
		netParams,
	)
	if err != nil {
		return nil, err
	}
	return address, nil
}

func makeMultiAddressTx(feerate, amount int64, outList []*PrevOutPoint, sender btcutil.Address,
	privs *btcec.PrivateKey, receivers []btcutil.Address, btcApiClient *mempool.MempoolClient) (*wire.MsgTx, error) {

	commitTx := wire.NewMsgTx(wire.TxVersion)
	totalAmount := btcutil.Amount(0)
	TxPrevOutputFetcher := txscript.NewMultiPrevOutFetcher(nil)
	count := len(receivers)

	// handle the every address's utxo
	for _, out := range outList {
		txOut, err := getTxOutByOutPoint2(out.Outpoint, btcApiClient)
		if err != nil {
			return nil, err
		}
		TxPrevOutputFetcher.AddPrevOut(*out.Outpoint, txOut)
		in := wire.NewTxIn(out.Outpoint, nil, nil)
		in.Sequence = defaultSequenceNum
		commitTx.AddTxIn(in)
		totalAmount += btcutil.Amount(out.Value)
	}
	for i := 0; i < count; i++ {
		PkScript0, err := txscript.PayToAddrScript(receivers[i])
		if err != nil {
			return nil, err
		}

		commitTx.AddTxOut(&wire.TxOut{
			PkScript: PkScript0,
			Value:    int64(amount),
		})
	}

	changePkScript, err := txscript.PayToAddrScript(sender)
	if err != nil {
		return nil, err
	}
	// make the change
	commitTx.AddTxOut(wire.NewTxOut(0, changePkScript))
	txsize := btcmempool.GetTxVirtualSize(btcutil.NewTx(commitTx))
	fee := btcutil.Amount(txsize) * btcutil.Amount(feerate)
	changeAmount := totalAmount - fee - btcutil.Amount(int64(count)*amount)

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
			i, txOut.Value, txOut.PkScript, txscript.SigHashDefault, privs)
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

func TestGeneratePrivateKey(t *testing.T) {
	privateKey, err := btcec.NewPrivateKey()
	if err != nil {
		t.Fatal(err)
	}
	privateKeyBytes := privateKey.Serialize()
	privateKeyString := hex.EncodeToString(privateKeyBytes)
	t.Logf("private key: %s", privateKeyString)

	privateKeyBytes, err = hex.DecodeString(privateKeyString)
	if err != nil {
		t.Fatal(err)
	}
	privateKey, _ = btcec.PrivKeyFromBytes(privateKeyBytes)
	privateKeyBytes = privateKey.Serialize()
	privateKeyString = hex.EncodeToString(privateKeyBytes)
	t.Logf("private key: %s", privateKeyString)

	network := &chaincfg.MainNetParams
	if testnet {
		network = &chaincfg.TestNet3Params
	}

	sender, err := makeTpAddress(privateKey, network)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("sender: ", sender.String())
}
func Test_getFeerate(t *testing.T) {
	network := &chaincfg.MainNetParams
	if testnet {
		network = &chaincfg.TestNet3Params
	}
	client := mempool.NewClient(network)
	fees, err := client.RecommendedFees()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("FastestFee:", fees.FastestFee, "HalfHourFee", fees.HalfHourFee)
}
func Test_01(t *testing.T) {
	//privateKeyBytes, err := hex.DecodeString(priv1)
	//if err != nil {
	//	panic(err)
	//}
	//senderPriv, _ := btcec.PrivKeyFromBytes(privateKeyBytes)
	//network := &chaincfg.MainNetParams
	//if testnet {
	//	network = &chaincfg.TestNet3Params
	//}
	//client := mempool.NewClient(network)
	//feerate := 50

	//txHash, err := client.BroadcastTx(tx)
	//if err != nil {
	//	panic(err)
	//}
}
