package task

import (
	"encoding/hex"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
	"testing"
)

func TestDecodeData1(t *testing.T) {
	type args struct {
		data string
	}
	tests := []struct {
		name    string
		args    args
		want    *SwapAndBridgeFunctionParams
		wantErr bool
	}{
		{
			name: "test",
			args: args{
				data: "6e1537da0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000f855a761f9182c4b22a04753681a1f6324ed3449000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000003635c9adc5dea00000000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000007c00000000000000000000000000000000000000000000000000000000000000ec00000000000000000000000000000000000000000000000000000000000000ee000000000000000000000000000000000000000000000000000000000000006a0000000000000000000000000000000000000000000000000000000000000002000000000000000000000000033daba9618a75a7aff103e53afe530fbacf4a3dd000000000000000000000000f855a761f9182c4b22a04753681a1f6324ed3449000000000000000000000000f855a761f9182c4b22a04753681a1f6324ed344900000000000000000000000000000000000000000000000071338559aa0b798b00000000000000000000000000000000000000000000000000000000000000a0000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000004000000000000000000000000002162b2aee2dd657fb131b28cc34dee6797b66f000000000000000000000000002162b2aee2dd657fb131b28cc34dee6797b66f00000000000000000000000000000000000000000000003635c9adc5dea0000000000000000000000000000000000000000000000000000000000000000000a000000000000000000000000000000000000000000000000000000000000004e00000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000006000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000444efa064650000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000000000000000000000000000000033daba9618a75a7aff103e53afe530fbacf4a3dd0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000f855a761f9182c4b22a04753681a1f6324ed344900000000000000000000000000000000000000000000000071338559aa0b798b00000000000000000000000000000000000000000000000000000000000000c00000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000000000000000000000000000003ef68d3f7664b2805d4e88381b64868a56f88bc40000000000000000000000003ef68d3f7664b2805d4e88381b64868a56f88bc400000000000000000000000000000000000000000000003635c9adc5dea0000000000000000000000000000000000000000000000000000000000000000000a00000000000000000000000000000000000000000000000000000000000000244ac9650d800000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000001c0000000000000000000000000000000000000000000000000000000000000014475ceafe6000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000a0000000000000000000000000002162b2aee2dd657fb131b28cc34dee6797b66f00000000000000000000000000000000000000000000003635c9adc5dea00000000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000ffffffff000000000000000000000000000000000000000000000000000000000000004213cb04d4a5dfb6398fc5ab005a6c84337256ee2300271005ab928d446d8ce6761e368c8e7be03c3168a9ec00271033daba9618a75a7aff103e53afe530fbacf4a3dd00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000412210e8a00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000006e0000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000890000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000000c0000000000000000000000000000000000000000000000000000000000000001415100fe7e94d03f00261b7c5d2de4708df33912d00000000000000000000000000000000000000000000000000000000000000000000000000000000000005e0000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000005c0000000000000000000000000000000000000000000000000000000000000056000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000000000000000000000000000000f855a761f9182c4b22a04753681a1f6324ed3449000000000000000000000000f855a761f9182c4b22a04753681a1f6324ed3449000000000000000000000000000000000000000000000000bb72e5d9b1862e3100000000000000000000000000000000000000000000000000000000000000a0000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000004000000000000000000000000002162b2aee2dd657fb131b28cc34dee6797b66f000000000000000000000000002162b2aee2dd657fb131b28cc34dee6797b66f00000000000000000000000000000000000000000000000000000000006dd52f00000000000000000000000000000000000000000000000000000000000000a000000000000000000000000000000000000000000000000000000000000003a00000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000006000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000304efa064650000000000000000000000000000000000000000000000000000000000000020000000000000000000000000c2132d05d31c914a87c6611c10748aeb04b58e8f00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000f855a761f9182c4b22a04753681a1f6324ed3449000000000000000000000000000000000000000000000000bb72e5d9b1862e3100000000000000000000000000000000000000000000000000000000000000c00000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000040000000000000000000000001111111254eeb25477b68fb85ed929f73a9605820000000000000000000000001111111254eeb25477b68fb85ed929f73a96058200000000000000000000000000000000000000000000000000000000006dd52f00000000000000000000000000000000000000000000000000000000000000a000000000000000000000000000000000000000000000000000000000000001200000000000000000000000000000000000000000000000000000000000000024000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000a8e449022e00000000000000000000000000000000000000000000000000000000006dd52f000000000000000000000000000000000000000000000000b7957a2eefd7de1000000000000000000000000000000000000000000000000000000000000000600000000000000000000000000000000000000000000000000000000000000001a000000000000000000000009b08288c3be4f62bbf8d1c20ac9c5e6f9467d8b752fd304d00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
			},
			want:    &SwapAndBridgeFunctionParams{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DecodeData(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("DecodeData() got = %v, want %v", got, tt.want)
			//}
			t.Logf("got: %+v\n", got)
			t.Logf("SwapData: %s\n", common.Bytes2Hex(got.SwapData))
			t.Logf("BridgeData: %s\n", common.Bytes2Hex(got.BridgeData))
			t.Logf("PermitData: %s\n", common.Bytes2Hex(got.PermitData))
			t.Logf("FeeData: %s\n", common.Bytes2Hex(got.FeeData))
		})
	}
}

func TestGenerateAccount(t *testing.T) {
	privateKey, err := generateKey()
	if err != nil {
		t.Fatal(err)
	}
	address, err := makeTaprootAddress(privateKey, &chaincfg.TestNet3Params)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("address: ", address.String())
	t.Log("privateKey: ", hex.EncodeToString(privateKey.Serialize()))
	t.Log("privateKey: ", privateKey.Key.String())
}

func TestParseHexNumberToBigInt(t *testing.T) {
	value, ok := new(big.Int).SetString("0f4240", 16)
	if !ok {
		t.Fatal("parse hex number to big int failed")
	}
	t.Log("value: ", value)
}

func TestEventHash(t *testing.T) {
	event := "OnReceived(bytes32,address,address,bytes,uint256,address)"
	eventHash := crypto.Keccak256Hash([]byte(event))
	t.Log("event hash: ", eventHash)
}

func generateKey() (*btcec.PrivateKey, error) {
	privateKey, err := btcec.NewPrivateKey()
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}

func makeTaprootAddress(privKey *btcec.PrivateKey, netParams *chaincfg.Params) (btcutil.Address, error) {
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
