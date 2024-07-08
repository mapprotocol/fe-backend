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
	"sync"
	"testing"
	"time"
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
	privateKeyBytes, err := hex.DecodeString("8b04a45a7f66395aa3f61fbd2bd1172b0a5f4e64891729dc9e49a9a9a6eb05fc")
	if err != nil {
		panic(err)
	}
	senderPriv, _ := btcec.PrivKeyFromBytes(privateKeyBytes)
	sender, _ := btcutil.DecodeAddress("tb1p23dgrhckt9vr24yuqdl3yu2xwj8em3wmn40ly0dtuf0lk0kk80jqesjhk4", network)

	feerate := int64(5)
	addrCount, amount := 5, int64(100)

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

	err = waitTxOnChain(txHash, client)
	if err != nil {
		fmt.Println("get tx state failed", err)
		return
	}
	fmt.Println("txhash on chain")
}

func Test_03(t *testing.T) {
	// collect test
	privStrs := []string{
		"c2a945ac319edcc497a2237acbc7077398b3f906efff9707bbd1c403928e6ada",
		"6a4b301b961c50dfd56d84cd0c28b33b87a33669c811840fadaa83ef7d00e66f",
		"b7f3592b75f5894ede91d408c4abffef21d0ca5e3f9d9fbec8ac9384b8595331",
		"6f4dbb1e415761d97e008d8cee77abb1475fdba033547e8486cb17b436e959c3",
		"253b06999cc7b531d6f53de8e43c1fc77e2bd238516c1d6d61e8954f7d32d34a",
	}
	addrStrs := []string{
		"tb1p42xadanfhg82s8wm4yw59ys4vuunsyyvzteacdvta2z4p4vrs29satct4l",
		"tb1pfx50n7wkdnha0rh3j70363phkv50s4pafrg6s9cfhhtxtxusf4xs46w8ve",
		"tb1pew99gkv36gerrs7shy4tpr952250n02flnz30ezjy2qmz2rd7h6qn4e676",
		"tb1psud7xj9sncur40xe4a3y72ngld2aq6pw36rjcn79ncq4ga656mjq8jnwgw",
		"tb1pn4ammcs3dyzyfg3tk39ss8ly6d5ndpu9z75c9fg693c2gads7v9q4l0yys",
	}
	network := &chaincfg.MainNetParams
	if testnet {
		network = &chaincfg.TestNet3Params
	}
	client := mempool.NewClient(network)
	privateKeyBytes, err := hex.DecodeString("8b04a45a7f66395aa3f61fbd2bd1172b0a5f4e64891729dc9e49a9a9a6eb05fc")
	if err != nil {
		panic(err)
	}
	feePriv, _ := btcec.PrivKeyFromBytes(privateKeyBytes)
	feeSender, _ := btcutil.DecodeAddress("tb1p23dgrhckt9vr24yuqdl3yu2xwj8em3wmn40ly0dtuf0lk0kk80jqesjhk4", network)
	receiver, _ := btcutil.DecodeAddress("tb1pwf8u8g9pxnnm3kleec2wwk790y0g7nuvm7qyvu7xl8752c9cqe7swdaakj", network)

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
			OrderID: uint64(i + 1),
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

func Test_fullCollection(t *testing.T) {
	network := &chaincfg.MainNetParams
	if testnet {
		network = &chaincfg.TestNet3Params
	}
	client := mempool.NewClient(network)

	privateKeyBytes, err := hex.DecodeString("8b04a45a7f66395aa3f61fbd2bd1172b0a5f4e64891729dc9e49a9a9a6eb05fc")
	if err != nil {
		panic(err)
	}
	feePriv, _ := btcec.PrivKeyFromBytes(privateKeyBytes)
	feeAddress, _ := btcutil.DecodeAddress("tb1p23dgrhckt9vr24yuqdl3yu2xwj8em3wmn40ly0dtuf0lk0kk80jqesjhk4", network)
	receiver, _ := btcutil.DecodeAddress("tb1pwf8u8g9pxnnm3kleec2wwk790y0g7nuvm7qyvu7xl8752c9cqe7swdaakj", network)

	addrCount, amount, feerate := 10, int64(100), int64(5)

	addrs, privs := make([]btcutil.Address, 0), make([]*btcec.PrivateKey, 0)
	items := make([]*OrderItem, 0)
	for i := 0; i < addrCount; i++ {
		priv, addr, err := makeNewTpAddress()
		if err != nil {
			panic(err)
		}
		addrs = append(addrs, addr)
		privs = append(privs, priv)
		item := &OrderItem{
			OrderID: uint64(i + 1),
			Sender:  addrs[i],
			Priv:    privs[i],
		}
		items = append(items, item)
		fmt.Println("priv:", hex.EncodeToString(priv.Serialize()), "addr:", addr.String())
	}
	outlist, err := gatherUTXO3(feeAddress, client)
	if err != nil {
		panic(err)
	}

	dTx, err := makeMultiAddressTx(feerate, amount, outlist, feeAddress, feePriv, addrs, client)
	if err != nil {
		panic(err)
	}
	txHash, err := client.BroadcastTx(dTx)
	if err != nil {
		panic(err)
	}
	fmt.Println("txhash:", txHash.String())

	err = waitTxOnChain(txHash, client)
	if err != nil {
		fmt.Println("get tx state failed", err)
		return
	}
	fmt.Println("txhash on chain")

	// collect process
	fmt.Println("make the collect tx.....")
	cTx, err := makeCollectTx1(feerate, receiver, feeAddress, feePriv, items, client)
	if err != nil {
		fmt.Println(err)
		return
	}
	txHash1, err := client.BroadcastTx(cTx)
	if err != nil {
		panic(err)
	}
	fmt.Println("collect the order...")
	fmt.Println("collect the txhash", txHash1.String())
	fmt.Println("wait the tx on the chain")
	waitTxOnChain(txHash1, client)
	fmt.Println("finish")
}

func Test_btc(t *testing.T) {
	amount, err := btcutil.NewAmount(1.2)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(int64(amount), amount)
}
func Test_makeTxhashAndSimpleTransfer(t *testing.T) {
	network := &chaincfg.MainNetParams
	if testnet {
		network = &chaincfg.TestNet3Params
	}
	client := mempool.NewClient(network)
	privateKeyBytes, err := hex.DecodeString("8b04a45a7f66395aa3f61fbd2bd1172b0a5f4e64891729dc9e49a9a9a6eb05fc")
	if err != nil {
		panic(err)
	}
	senderPriv, _ := btcec.PrivKeyFromBytes(privateKeyBytes)
	sender, _ := btcutil.DecodeAddress("tb1p23dgrhckt9vr24yuqdl3yu2xwj8em3wmn40ly0dtuf0lk0kk80jqesjhk4", network)

	privateKeyBytes, err = hex.DecodeString("5a538bd636712fdd4855c74d8f4b3438b06bee3a1240adc4ee5f2f6b8393dfb9")
	if err != nil {
		panic(err)
	}
	tipperPriv, _ := btcec.PrivKeyFromBytes(privateKeyBytes)
	tipper, _ := btcutil.DecodeAddress("tb1pj3jq046hzh9k3czd55v6ldvey3quhqagzm7wnueszjukjjah950s4rwhdu", network)

	receiver, _ := btcutil.DecodeAddress("tb1pwf8u8g9pxnnm3kleec2wwk790y0g7nuvm7qyvu7xl8752c9cqe7swdaakj", network)

	senderOutlist, err := gatherUTXO3(sender, client)
	if err != nil {
		fmt.Println(err)
		return
	}
	feeOutlist, err := gatherUTXO3(tipper, client)
	if err != nil {
		fmt.Println(err)
		return
	}

	feerate := int64(25)

	tx, err := makeSimpleTx0(feerate, 500, sender, receiver, tipper, senderPriv, tipperPriv, senderOutlist, feeOutlist, client)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("txhash1", tx.TxHash().String())
	txHash, err := client.BroadcastTx(tx)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("txhash2", txHash.String())
	err = waitTxOnChain(txHash, client)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("tx was on chain")
}

func Test_04(t *testing.T) {
}

func Test_channe01(t *testing.T) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	index := 0
	for {
		select {
		case <-ticker.C:
			index++
			fmt.Println("ticker index", index, time.Now())
			if index > 10 {
				return
			}
		default:
		}
	}
}
func Test_channe02(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(4)
	chBalanceLow, chBalanceHigh := make(chan int), make(chan int)
	ticker1 := time.NewTicker(1 * time.Second)

	go func() {
		defer wg.Done()
		action, pos := true, 0

		for {
			select {
			case <-ticker1.C:
				i := pos % 5
				if action {
					if i == 0 {
						action = false
						chBalanceLow <- 1
						fmt.Println("action=false,chLow set...")
					}
					fmt.Println("set withdraw state....")
				}
				pos++
			case msg := <-chBalanceHigh:
				if msg == 2 {
					action = true
					fmt.Println("set action=true,begin the channel....[withdraw]")
				}
			}
		}
	}()

	go func() {
		defer wg.Done()
		pos := 0

		for {
			select {
			case <-ticker1.C:
				i := pos % 50
				if i == 0 {
					fmt.Println("set ch2 state....")
				}
				pos++
			case msg := <-chBalanceLow:
				if msg == 1 {
					fmt.Println("begin the balance channel,will set action=true....")
					time.Sleep(5 * time.Second)
					chBalanceHigh <- 2
					fmt.Println("end the balance channel....")
				}
			}
		}
	}()

	wg.Wait()
	close(chBalanceLow)
	close(chBalanceHigh)
}
