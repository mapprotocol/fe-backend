package task

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"testing"

	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/liteclient"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/ton"
	"github.com/xssnick/tonutils-go/ton/wallet"
	"github.com/xssnick/tonutils-go/tvm/cell"

	"github.com/mapprotocol/fe-backend/resource/tonclient"
)

func TestTonFilterEvent(t *testing.T) {
	client := liteclient.NewConnectionPool()

	cfg, err := liteclient.GetConfigFromUrl(context.Background(), "https://ton.org/global.config.json")
	//cfg, err := liteclient.GetConfigFromUrl(context.Background(), "https://ton.org/testnet-global.config.json")
	if err != nil {
		log.Fatalln("get config err: ", err.Error())
		return
	}

	// connect to mainnet lite servers
	err = client.AddConnectionsFromConfig(context.Background(), cfg)
	if err != nil {
		log.Fatalln("connection err: ", err.Error())
		return
	}

	// initialize ton api lite connection wrapper with full proof checks
	//api := ton.NewAPIClient(client, ton.ProofCheckPolicySecure).WithRetry() // todo
	api := ton.NewAPIClient(client).WithRetry()
	api.SetTrustedBlockFromConfig(cfg)

	log.Println("fetching and checking proofs since config init block, it may take near a minute...")
	master, err := api.CurrentMasterchainInfo(context.Background()) // we fetch block just to trigger chain proof check
	if err != nil {
		log.Fatalln("get masterchain info err: ", err.Error())
		return
	}
	log.Println("master proof checks are completed successfully, now communication is 100% safe!, master: ", master)

	// address on which we are accepting payments
	//treasuryAddress := address.MustParseAddr("EQCcRLZA4yLYKJOQExatRMokKXv92kHyeZF7SKjG8TQO4kAI")
	treasuryAddress := address.MustParseAddr("EQAe8Lq5uwe_OTE3qh0iircglVnsETikIPK6Qp4DoKa9sWQ6")
	//treasuryAddress := address.MustParseAddr("EQCjQ5CNGkoyp8xW-V75g1v1-CmClYCfJLTGXOmaKxQTDuIX")
	//treasuryAddress := address.MustParseAddr("EQDYMWIoZEncsefn-NS40_L6hecmGI9JvssQu7jN_QHJH0p6")
	//treasuryAddress := address.MustParseAddr("EQBykWvdmoyFD2kB4BJbdrxRqyWKzSLDn-SgwNvznRfbsaKv")
	// https://tonscan.org/tx/4SzFuwn%2FDdMMteesZcAAIzk0GTCsUV5SXcGHfRJ3hdw=
	// https://tonscan.org/tx/XjLl%2FKcRed0gX7umFXsX1qqdxL1NR%2FUT%2FWhbBD9uT+E=

	acc, err := api.GetAccount(context.Background(), master, treasuryAddress)
	if err != nil {
		log.Fatalln("get master chain info err: ", err.Error())
		return
	}

	// Cursor of processed transaction, save it to your db
	// We start from last transaction, will not process transactions older than we started from.
	// After each processed transaction, save lt to your db, to continue after restart
	log.Println("last tx lt: ", acc.LastTxLT)
	//lastProcessedLT := acc.LastTxLT
	lastProcessedLT := uint64(47839075000001)
	// channel with new transactions
	transactions := make(chan *tlb.Transaction)

	// it is a blocking call, so we start it asynchronously
	go api.SubscribeOnTransactions(context.Background(), treasuryAddress, lastProcessedLT, transactions)

	log.Println("waiting for transfers...")

	// listen for new transactions from channel
	for tx := range transactions {
		// process transaction here
		log.Println("tx: ", tx.String())
		out := tx.IO.Out
		if out == nil {
			log.Println("no out messages")
			continue
		}
		messages, err := out.ToSlice()
		if err != nil {
			log.Println("get out messages err: ", err.Error())
			continue
		}
		for _, msg := range messages {
			log.Println("msg type: ", msg.MsgType)
			if msg.MsgType != tlb.MsgTypeExternalOut {
				continue
			}

			out := msg.AsExternalOut()
			log.Println("src: ", out.SrcAddr)
			log.Println("dst: ", out.DstAddr)
			log.Println("dst: ", out.CreatedLT)
			log.Println("dst: ", out.CreatedAt)

			payload := msg.AsExternalOut().Payload()
			log.Printf("payload: %+v\n", payload)

			//slice := payload.BeginParse()
			//orderID := slice.MustLoadBigUInt(256)
			//sender := slice.MustLoadAddr()
			//srcToken := slice.MustLoadAddr()
			//srcAmount := slice.MustLoadBigUInt(256)
			//log.Println("orderID", orderID)
			//log.Println("sender: ", sender)
			//log.Println("srcToken: ", srcToken)
			//log.Println("srcAmount: ", srcAmount)

			//t.Log("", slice.MustLoadUInt(4))
			//t.Log("", slice.MustLoadUInt(2))
			//t.Log("", slice.MustLoadUInt(9))
			//t.Log("event id: ", slice.MustLoadUInt(256))
			//t.Log("", slice.MustLoadUInt(64+32+2))

			//slice := payload.BeginParse()
			//t.Log("data 1: ", slice.MustLoadSlice(32))
			//t.Log("data 2: ", slice.MustLoadSlice(32))

			//slice := payload.BeginParse()
			//t.Log("data: ", slice.MustLoadUInt(64))
			//t.Log("data: ", slice.MustLoadUInt(32))
			//t.Log("data: ", slice.MustLoadAddr())
			//t.Log("data: ", slice.MustLoadCoins()) // 8057622

			slice := payload.BeginParse()
			t.Log("orderID: ", slice.MustLoadUInt(64))
			from := slice.MustLoadRef()
			to := slice.MustLoadRef()

			t.Log("sender: ", from.MustLoadAddr())
			t.Log("srcToken: ", from.MustLoadAddr())
			t.Log("amountIn: ", from.MustLoadUInt(32))

			t.Log("toChainID: ", to.MustLoadUInt(32))
			t.Log("receiver =", to.MustLoadUInt(160)) // 0xce83e2e5ea8a6fbe
			t.Log("tokenOutAddress =", to.MustLoadUInt(160))
			t.Log("receiver: ", "0x"+hex.EncodeToString(new(big.Int).SetUint64(to.MustLoadUInt(160)).Bytes())) // 0xce83e2e5ea8a6fbe
			t.Log("tokenOutAddress: ", "0x"+hex.EncodeToString(new(big.Int).SetUint64(to.MustLoadUInt(160)).Bytes()))
			t.Log("jetton amount: ", slice.MustLoadUInt(32))

			// receiver = '14880987070872580030'
			// BigInt(receiver).toString(16)
			//amountOut = 2449135357229387394

		}
		// 0 0 0 0 0 1 8 171  BigInt(tokenOutAddress).toString(16)
		// 0 0 0 137
		// 128 26 52 231 7 51 185 242 188 246 17 52 137 24 208 83 246 112 155 184 9 248 11 122 202 189 48 0 72 94 186 111 94
		// 111 94 98 220
		// update last processed lt and save it in db
		lastProcessedLT = tx.LT
	}

	// it can happen due to none of available liteservers know old enough state for our address
	// (when our unprocessed transactions are too old)
	log.Println("something went wrong, transaction listening unexpectedly finished")
}
func TestTonFilterEvents(t *testing.T) {
	client := liteclient.NewConnectionPool()

	cfg, err := liteclient.GetConfigFromUrl(context.Background(), "https://ton.org/global.config.json")
	//cfg, err := liteclient.GetConfigFromUrl(context.Background(), "https://ton.org/testnet-global.config.json")
	if err != nil {
		log.Fatalln("get config err: ", err.Error())
		return
	}

	// connect to mainnet lite servers
	err = client.AddConnectionsFromConfig(context.Background(), cfg)
	if err != nil {
		log.Fatalln("connection err: ", err.Error())
		return
	}

	// initialize ton api lite connection wrapper with full proof checks
	//api := ton.NewAPIClient(client, ton.ProofCheckPolicySecure).WithRetry() // todo
	api := ton.NewAPIClient(client).WithRetry()
	api.SetTrustedBlockFromConfig(cfg)

	log.Println("fetching and checking proofs since config init block, it may take near a minute...")
	master, err := api.CurrentMasterchainInfo(context.Background()) // we fetch block just to trigger chain proof check
	if err != nil {
		log.Fatalln("get masterchain info err: ", err.Error())
		return
	}
	log.Println("master proof checks are completed successfully, now communication is 100% safe!, master: ", master)

	// address on which we are accepting payments
	//treasuryAddress := address.MustParseAddr("EQCcRLZA4yLYKJOQExatRMokKXv92kHyeZF7SKjG8TQO4kAI")
	treasuryAddress := address.MustParseAddr("EQBM9zTT58eBMeArcb3dX8xtH1kBcPVZg5D_Ef6PRCZ6lMve")

	acc, err := api.GetAccount(context.Background(), master, treasuryAddress)
	if err != nil {
		log.Fatalln("get master chain info err: ", err.Error())
		return
	}

	// Cursor of processed transaction, save it to your db
	// We start from last transaction, will not process transactions older than we started from.
	// After each processed transaction, save lt to your db, to continue after restart
	log.Println("last tx lt: ", acc.LastTxLT)
	//lastProcessedLT := acc.LastTxLT
	lastProcessedLT := uint64(48112392000001)
	// channel with new transactions
	transactions := make(chan *tlb.Transaction)

	// it is a blocking call, so we start it asynchronously
	go api.SubscribeOnTransactions(context.Background(), treasuryAddress, lastProcessedLT, transactions)

	log.Println("waiting for transfers...")

	// listen for new transactions from channel
	for tx := range transactions {
		// process transaction here
		log.Println("tx: ", tx.String())
		out := tx.IO.Out
		if out == nil {
			continue
		}
		messages, err := out.ToSlice()
		if err != nil {
			log.Println("get out messages err: ", err.Error())
			continue
		}
		for _, msg := range messages {
			if msg.MsgType != tlb.MsgTypeExternalOut {
				continue
			}
			//

			externalOut := msg.AsExternalOut()
			log.Println("src: ", externalOut.SrcAddr)
			log.Println("dst: ", externalOut.DstAddr) // EXT:110000010000000000000000000000000000000000000000000000000000000000c0470ccf
			// todo 判断 event id

			payload := externalOut.Payload()
			if payload == nil {
				continue
			}
			data, err := payload.MarshalJSON()
			if err != nil {
				t.Fatal(err)
			}
			t.Log("data: ", hex.EncodeToString(data))

			slice := payload.BeginParse()
			if slice == nil {
				continue
			}
			t.Log("orderID: ", slice.MustLoadUInt(64))
			from := slice.MustLoadRef()
			to := slice.MustLoadRef()
			t.Log("sender: ", from.MustLoadAddr())
			t.Log("srcToken: ", from.MustLoadAddr())
			t.Log("amountIn: ", from.MustLoadUInt(32))

			t.Log("toChainID: ", to.MustLoadUInt(32))
			t.Log("receiver: ", "0x"+hex.EncodeToString(to.MustLoadBigUInt(160).Bytes()))
			t.Log("tokenOutAddress: ", "0x"+hex.EncodeToString(to.MustLoadBigUInt(160).Bytes()))
			t.Log("jetton amount: ", slice.MustLoadUInt(32)) // 0x0e367ce43859b170ef7dc147ce83e2e5ea8a6fbe
		}
		lastProcessedLT = tx.LT
	}

	// it can happen due to none of available liteservers know old enough state for our address   e367ce43859b180000000000000000000000000 0e367ce43859B170EF7DC147CE83e2e5ea8A6fbe
	// (when our unprocessed transactions are too old)
	log.Println("something went wrong, transaction listening unexpectedly finished")
}

func TestStoreAndLoadCell(t *testing.T) {
	ref1 := cell.BeginCell().
		MustStoreStringSnake("0xe0dc8d7f134d0a79019bef9c2fd4b2013a64fcd6").
		MustStoreStringSnake("0x624e6f327c4f91f1fa6285711245c215de264d49").
		MustStoreStringSnake("0x624e6f327c4f91f1fa6285711245c215de264d49").
		MustStoreStringSnake("0x624e6f327c4f91f1fa6285711245c215de264d49").
		EndCell()

	c1 := cell.BeginCell().
		MustStoreBigUInt(big.NewInt(1234567890), 64).
		MustStoreBigUInt(big.NewInt(2000000000000000000), 256).
		MustStoreAddr(address.MustParseAddr("EQBykWvdmoyFD2kB4BJbdrxRqyWKzSLDn-SgwNvznRfbsaKv")).
		MustStoreRef(ref1).
		EndCell()

	c2 := c1.BeginParse()

	t.Logf("c1: %+v\n", c1)
	t.Logf("ref: %+v\n", c1.MustPeekRef(0))
	t.Log("orderID", c2.MustLoadBigUInt(64))
	t.Log("amount", c2.MustLoadBigUInt(256))
	t.Log("address", c2.MustLoadAddr())
	//t.Log("sender", c2.MustLoadStringSnake())
	//t.Log("receiver", c2.MustLoadStringSnake())
	ref2 := c2.MustLoadRef()
	t.Log("string: ", ref2.MustLoadStringSnake())
}

func TestNewSeedWithPassword(t *testing.T) {
	t.Log("seed: ", wallet.NewSeedWithPassword("CeD89#0F17b5+kcT3b"))
}

func TestSendTransaction(t *testing.T) {
	words := "cushion bean assault oven hybrid account lunch festival valid soap history grant horn good decline tourist shadow very eye language person venture term shove"
	password := "J%A7k7sGXe4i58G#fN"
	tonclient.Init(words, password)

	log.Println("wallet address:", tonclient.Wallet().WalletAddress())

	// ================================

	block, err := tonclient.Client().CurrentMasterchainInfo(context.Background())
	if err != nil {
		log.Fatalln("CurrentMasterchainInfo err:", err.Error())
		return
	}

	balance, err := tonclient.Wallet().GetBalance(context.Background(), block)
	if err != nil {
		log.Fatalln("GetBalance err:", err.Error())
		return
	}

	if balance.Nano().Uint64() >= 3000000 {
		// data.txParams.to
		to := address.MustParseAddr("EQDa4VOnTYlLvDJ0gZjNYm5PXfSmmtL6Vs6A_CZEtXCNICq_")
		// data.txParams.value
		//value, err := tlb.FromNanoTONStr("10100000") // 0.01
		//value, err := tlb.FromNanoTONStr("110010000") // 0.01
		//value, err := tlb.FromNanoTONStr("101001000") // 0.001
		//value, err := tlb.FromNanoTONStr("100100100") // 0.0001
		value, err := tlb.FromNanoTONStr("600500000") // 0.01
		if err != nil {
			t.Fatal(err)
		}
		// data.txParams.data
		data := "te6cckEBAwEA4AABZeoGGF0AAAAAAAAAAEHc1lAIAHy/+VG7+ebYbZP+jeYqxVVqNzIpCLWthNl7zJOrFKsQBAEBS2afG6+AA94XVzdg9+cmJvVDpFFW5BKrPYInFIQeV0hTwHQU17YlAgD982HvywAAAAAAAATTgARVH/VZdHbitveNDGpmZ/tXNshmA1J1CaG4P91P5R6zkAB7wurm7B785MTeqHSKKtyCVWewROKQg8rpCngOgpr2xHc1lAAAAADgOsWqc/fjjpEce8ZDuY95aJPOHh8ITLQXTHJFKh8ZhHBB0iusEtY6Pi9CfRs="
		body := &cell.Cell{}
		if err := json.Unmarshal([]byte(fmt.Sprintf(`"%s"`, data)), &body); err != nil {
			t.Log(err)
		}
		t.Logf("body: %+v\n", body)

		log.Println("sending transaction and waiting for confirmation...")

		tx, block, err := tonclient.Wallet().SendWaitTransaction(context.Background(), &wallet.Message{
			Mode: wallet.PayGasSeparately, // pay fees separately (from balance, not from amount)
			InternalMessage: &tlb.InternalMessage{
				Bounce:  true, // return amount in case of processing error
				DstAddr: to,
				Amount:  value,
				Body:    body,
			},
		})
		if err != nil {
			log.Fatalln("Send err:", err.Error())
			return
		}

		log.Println("transaction sent, confirmed at block, hash:", base64.StdEncoding.EncodeToString(tx.Hash))

		balance, err = tonclient.Wallet().GetBalance(context.Background(), block)
		if err != nil {
			log.Fatalln("GetBalance err:", err.Error())
			return
		}

		log.Println("balance left:", balance.String())
		if balance.Nano().Uint64() < 3000000 {
			log.Println("ton account not enough balance:", balance.String())
		}

		return
	}

	log.Println("not enough balance:", balance.String())
}

func TestBase64Decode(t *testing.T) {
	data := "te6cckEBAwEA4AABZeoGGF0AAAAAAAAAAEBfXhAIAHy/+VG7+ebYbZP+jeYqxVVqNzIpCLWthNl7zJOrFKsQBAEBS2ad8M2ACxdsw9znlhEWQFW/HPuBC1xnxDkW0cbDUzylgznadUgFAgD982HvywAAAAAAAATTgARVH/VZdHbitveNDGpmZ/tXNshmA1J1CaG4P91P5R6zkAFi7Zh7nPLCIsgKt+OfcCFrjPiHIto42GpnlLBnO06pABfXhAAAAADgOsWqc/fjjpEce8ZDuY95aJPOHh8ITLQXTHJFKh8ZhHBB0iusEtY6PmIlJ0o="

	body := &cell.Cell{}
	if err := json.Unmarshal([]byte(fmt.Sprintf(`"%s"`, data)), &body); err != nil {
		t.Log(err)
	}

	t.Log(*body)

	//c1 := cell.BeginCell().
	//	MustStoreBigUInt(big.NewInt(1234567890), 64).
	//	MustStoreBigUInt(big.NewInt(2000000000000000000), 256).
	//	EndCell()
	//t.Logf("c1: %+v\n", c1)
	//
	//marshal, err := json.Marshal(c1)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//t.Log("str: ", string(marshal))
	//
	//c2 := cell.Cell{}
	//err = json.Unmarshal(marshal, &c2)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//t.Logf("c2: %+v\n", c2)
	//
	//if !reflect.DeepEqual(c1, c2) {
	//	t.Error("not equal")
	//}
}

func Test12(t *testing.T) {
	t.Log(tlb.FromNanoTONStr("600500000"))
	t.Log(tlb.FromNanoTONStr("100000000"))
	t.Log(tlb.MustFromNano(big.NewInt(3314083), 6))
}

func TestName(t *testing.T) {
	amountFloat, ok := new(big.Float).SetString("3489497")
	if !ok {
		t.Fatal("convert failed")
	}
	amountFloat = new(big.Float).Quo(amountFloat, big.NewFloat(1e6))
	t.Log(amountFloat)
}

//func TestParseTONToEVMEvent(t *testing.T) {
//	data := "2274653663636b4542417745416541414347455941414264496475674441464e4c6267454341476341424e5543414141414159414a315a73635866366e306a543454506833456446497a54775246714f54474b3056557a6274352b4f7469654141414141536f46386741417945414741414141414141414141695134446a366257657a3264486a4b3355595458302f4c6a6b5263394141414141414141414141414141414141414141414141414141442f4130487122"
//	logData, err := hex.DecodeString(data)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	body := &cell.Cell{}
//	if err := json.Unmarshal(logData, &body); err != nil {
//		t.Fatal(err)
//	}
//	slice := body.BeginParse()
//	orderID, err := slice.LoadUInt(64)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	from, err := slice.LoadRef()
//	if err != nil {
//		t.Fatal(err)
//	}
//	srcChain, err := from.LoadUInt(64)
//	if err != nil {
//		t.Fatal(err)
//	}
//	sender, err := from.LoadAddr()
//	if err != nil {
//		t.Fatal(err)
//	}
//	srcToken, err := from.LoadAddr()
//	if err != nil {
//		t.Fatal(err)
//	}
//	inAmount, err := from.LoadBigInt(64)
//	if err != nil {
//		t.Fatal(err)
//	}
//	t.Log(orderID, srcChain, sender, srcToken, inAmount)
//	slippage, err := from.LoadUInt(16)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	to, err := slice.LoadRef()
//	if err != nil {
//		t.Fatal(err)
//	}
//	dstChain, err := to.LoadUInt(64)
//	if err != nil {
//		t.Fatal(err)
//	}
//	receiverBigInt, err := to.LoadBigUInt(160)
//	if err != nil {
//		t.Fatal(err)
//	}
//	t.Log("======== receiver: ", common.BytesToAddress(receiverBigInt.Bytes()))
//
//	receiver := "0x" + hex.EncodeToString(receiverBigInt.Bytes())
//	tokenOut, err := to.LoadBigUInt(160)
//	if err != nil {
//		t.Fatal(err)
//	}
//	t.Log("======== tokenOut: ", common.BytesToAddress(tokenOut.Bytes()))
//
//	dstToken := "0x" + hex.EncodeToString(tokenOut.Bytes())
//	relayAmount, err := slice.LoadUInt(32)
//	if err != nil {
//		t.Fatal(err)
//	}
//	scrTokenStr := srcToken.String()
//	if scrTokenStr == params.NoneAddress {
//		scrTokenStr = params.NativeOfTON
//	}
//	_, afterAmount := deductFees(new(big.Float).SetUint64(relayAmount), FeeRate)
//	// convert token to float like 0.089
//	afterAmountFloat := new(big.Float).Quo(afterAmount, big.NewFloat(params.USDTDecimalOfTON))
//	inAmountFloat := new(big.Float).Quo(new(big.Float).SetInt(inAmount), big.NewFloat(params.InAmountDecimalOfTON))
//	//inAmountFloat := new(big.Float).Quo(new(big.Float).SetUint64(inAmount), big.NewFloat(params.InAmountDecimalOfTON))
//	order := &dao.Order{
//		OrderIDFromContract: orderID,
//		SrcChain:            strconv.FormatUint(srcChain, 10),
//		SrcToken:            scrTokenStr,
//		Sender:              sender.String(),
//		InAmount:            inAmountFloat.Text('f', -1),
//		RelayToken:          params.USDTOfTON,
//		RelayAmountInt:         afterAmountFloat.String(),
//		DstChain:            strconv.FormatUint(dstChain, 10),
//		DstToken:            dstToken,
//		Receiver:            receiver,
//		Action:              dao.OrderActionToEVM,
//		Stage:               dao.OrderStag1,
//		Status:              dao.OrderStatusTxConfirmed,
//		Slippage:            slippage,
//	}
//	t.Logf("order: %+v\n", order)
//}
