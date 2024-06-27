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
	testnet = true
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

func makeNewTpAddress() (*btcec.PrivateKey, btcutil.Address, error) {
	privateKey, err := btcec.NewPrivateKey()
	if err != nil {
		return nil, nil, err
	}
	network := &chaincfg.MainNetParams
	if testnet {
		network = &chaincfg.TestNet3Params
	}

	addr, err := makeTpAddress(privateKey, network)
	if err != nil {
		return nil, nil, err
	}
	return privateKey, addr, nil
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
	network := &chaincfg.MainNetParams
	if testnet {
		network = &chaincfg.TestNet3Params
	}
	client := mempool.NewClient(network)
	height, err := client.BlockHeight()
	if err != nil {
		panic(err)
	}
	fmt.Println("height", height.String())
}
func Test_02(t *testing.T) {
	network := &chaincfg.MainNetParams
	if testnet {
		network = &chaincfg.TestNet3Params
	}
	client := mempool.NewClient(network)
	privateKeyBytes, err := hex.DecodeString(priv1)
	if err != nil {
		panic(err)
	}
	senderPriv, _ := btcec.PrivKeyFromBytes(privateKeyBytes)
	sender, _ := btcutil.DecodeAddress("", network)

	feerate := int64(5)
	addrCount, amount := 3, int64(100)

	addrs := make([]btcutil.Address, 0)
	for i := 0; i < addrCount; i++ {
		priv, addr, err := makeNewTpAddress()
		if err != nil {
			panic(err)
		}
		addrs = append(addrs, addr)
		fmt.Println("priv:", hex.EncodeToString(priv.Serialize()), "addr:", addr.String())
	}
	outlist, err := gatherUTXO3(sender, client)
	if err != nil {
		panic(err)
	}

	tx, err := makeMultiAddressTx(feerate, amount, outlist, sender, senderPriv, addrs, client)
	if err != nil {
		panic(err)
	}
	txHash, err := client.BroadcastTx(tx)
	if err != nil {
		panic(err)
	}
	fmt.Println("txhash:", txHash.String())

	onChain, err := waitTxOnChain(txHash, client)
	if err != nil {
		fmt.Println("get tx state failed", err)
		return
	}
	fmt.Println("txhash on chain", onChain)
}

func Test_03(t *testing.T) {
	// collect test
	privStrs := []string{
		"",
		"",
	}
	addrStrs := []string{
		"",
		"",
	}
	network := &chaincfg.MainNetParams
	if testnet {
		network = &chaincfg.TestNet3Params
	}
	client := mempool.NewClient(network)
	privateKeyBytes, err := hex.DecodeString(priv1)
	if err != nil {
		panic(err)
	}
	feePriv, _ := btcec.PrivKeyFromBytes(privateKeyBytes)
	feeSender, _ := btcutil.DecodeAddress("", network)
	receiver, _ := btcutil.DecodeAddress("", network)

	feerate, addrCount := int64(5), len(addrStrs)
	privs, addrs := make([]*btcec.PrivateKey, 0), make([]btcutil.Address, 0)
	for _, astr := range addrStrs {
		addr, _ := btcutil.DecodeAddress(astr, network)
		addrs = append(addrs, addr)
	}
	for _, astr := range privStrs {
		pbytes, err := hex.DecodeString(astr)
		if err != nil {
			panic(err)
		}
		priv, _ := btcec.PrivKeyFromBytes(pbytes)
		privs = append(privs, priv)
	}
	// make the tmp orders
	items := make([]*OrderItem, 0)
	for i := 0; i < addrCount; i++ {
		item := &OrderItem{
			OrderID: int64(i + 1),
			Sender:  addrs[i],
			Priv:    privs[i],
		}
		items = append(items, item)
	}

	tx, err := makeCollectTx1(feerate, receiver, feeSender, feePriv, items, client)
	if err != nil {
		fmt.Println(err)
		return
	}
	txHash, err := client.BroadcastTx(tx)
	if err != nil {
		panic(err)
	}
	fmt.Println("collect the order...")
	fmt.Println("collect the txhash", txHash.String())
	fmt.Println("wait the tx on the chain")
	waitTxOnChain(txHash, client)
	fmt.Println("finish")
}
