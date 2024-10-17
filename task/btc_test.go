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
	"github.com/mapprotocol/fe-backend/params"
	"github.com/shopspring/decimal"
	"math/big"
	"reflect"
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

func TestConvert(t *testing.T) {
	amount := 0.12345678
	amountSat := amount / 1e8
	t.Log("amountSat: ", amountSat)
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

func Test_deductFees(t *testing.T) {
	type args struct {
		amount  *big.Int
		feeRate *big.Int
	}
	tests := []struct {
		name            string
		args            args
		wantFeeAmount   *big.Int
		wantAfterAmount *big.Int
	}{
		{
			name: "case1",
			args: args{
				amount:  big.NewInt(12345678),
				feeRate: big.NewInt(70),
			},
			wantFeeAmount:   big.NewInt(86419),
			wantAfterAmount: big.NewInt(12259259),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFeeAmount, gotAfterAmount := deductFees(tt.args.amount, tt.args.feeRate)
			if !reflect.DeepEqual(gotFeeAmount, tt.wantFeeAmount) {
				t.Errorf("deductFees() gotFeeAmount = %v, want %v", gotFeeAmount, tt.wantFeeAmount)
			}
			if !reflect.DeepEqual(gotAfterAmount, tt.wantAfterAmount) {
				t.Errorf("deductFees() gotAfterAmount = %v, want %v", gotAfterAmount, tt.wantAfterAmount)
			}
		})
	}
}

func Test_deductToTONBridgeFees(t *testing.T) {
	type args struct {
		amount        *big.Int
		bridgeFeeRate *big.Int
		swap          bool
	}
	tests := []struct {
		name            string
		args            args
		wantBridgeFees  *big.Int
		wantAfterAmount *big.Int
	}{
		{
			name: "t-1",
			args: args{
				amount:        convertDecimal(new(big.Int).Mul(big.NewInt(100), big.NewInt(1e18)), params.USDTDecimalNumberOfChainPool, params.FixedDecimalNumber), // 100 USDT
				bridgeFeeRate: big.NewInt(30),
			},
			wantBridgeFees:  big.NewInt(180000000),  // 0.3 + 1.5 = 1.8 USDT
			wantAfterAmount: big.NewInt(9820000000), // 98.2 USDT
		},
		// 32000000 * 30 / 10000 = 96000
		{
			name: "t-2",
			args: args{
				amount:        convertDecimal(new(big.Int).Mul(big.NewInt(32), big.NewInt(1e18)), params.USDTDecimalNumberOfChainPool, params.FixedDecimalNumber), // 32 USDT
				bridgeFeeRate: big.NewInt(30),
			},
			wantBridgeFees:  big.NewInt(159600000),  // 0.096 + 1.5 = 1.596 USDT
			wantAfterAmount: big.NewInt(3040400000), // 30.404 USDT
		},
		// 69000000 * 30 / 10000 = 207000
		{
			name: "t-3",
			args: args{
				amount:        convertDecimal(new(big.Int).Mul(big.NewInt(69), big.NewInt(1e18)), params.USDTDecimalNumberOfChainPool, params.FixedDecimalNumber), // 69 USDT
				bridgeFeeRate: big.NewInt(30),
			},
			wantBridgeFees:  big.NewInt(170700000),  // 20700000 + 150000000 = 1.707000 USDT
			wantAfterAmount: big.NewInt(6729300000), // 67.293 USDT
		},
		// 238457000000 * 30 / 10000 = 715371000
		// 715371000 + 1500000 = 716871000
		// 238457000000 - 716871000 = 237741629000
		{
			name: "t-3",
			args: args{
				amount:        convertDecimal(new(big.Int).Mul(big.NewInt(238457), big.NewInt(1e18)), params.USDTDecimalNumberOfChainPool, params.FixedDecimalNumber), // 238457 USDT
				bridgeFeeRate: big.NewInt(30),
			},
			wantBridgeFees:  big.NewInt(71687100000),    // 71537100000 + 150000000 = 716.871 USDT
			wantAfterAmount: big.NewInt(23774012900000), // 237740.129 USDT
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBridgeFees, gotAfterAmount := deductToTONBridgeFees(tt.args.amount, tt.args.bridgeFeeRate, tt.args.swap)
			if !reflect.DeepEqual(gotBridgeFees, tt.wantBridgeFees) {
				t.Errorf("deductToTonBridgeFees() gotBridgeFees = %v, want %v", gotBridgeFees, tt.wantBridgeFees)
			}
			if !reflect.DeepEqual(gotAfterAmount, tt.wantAfterAmount) {
				t.Errorf("deductToTonBridgeFees() gotAfterAmount = %v, want %v", gotAfterAmount, tt.wantAfterAmount)
			}
		})
	}
}

func Test_deductToBitcoinBridgeFees(t *testing.T) {
	type args struct {
		amount        *big.Int
		bridgeFeeRate *big.Int
	}
	tests := []struct {
		name            string
		args            args
		wantBridgeFees  *big.Int
		wantAfterAmount *big.Int
	}{
		{
			name: "t-1",
			args: args{
				amount:        convertDecimal(new(big.Int).Mul(big.NewInt(10), big.NewInt(1e18)), params.WBTCDecimalNumberOfChainPool, params.FixedDecimalNumber),
				bridgeFeeRate: big.NewInt(30),
			},
			wantBridgeFees:  big.NewInt(3001400),
			wantAfterAmount: big.NewInt(996998600),
		},
		{
			name: "t-2",
			args: args{
				amount:        convertDecimal(new(big.Int).Mul(big.NewInt(9636), big.NewInt(1e15)), params.WBTCDecimalNumberOfChainPool, params.FixedDecimalNumber), // 963600000
				bridgeFeeRate: big.NewInt(30),
			},
			wantBridgeFees:  big.NewInt(2892200),
			wantAfterAmount: big.NewInt(960707800),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			globalFeeRate = 5
			gotBridgeFees, gotAfterAmount := deductToBitcoinBridgeFees(tt.args.amount, tt.args.bridgeFeeRate)
			if !reflect.DeepEqual(gotBridgeFees, tt.wantBridgeFees) {
				t.Errorf("deductToBitcoinBridgeFees() gotBridgeFees = %v, want %v", gotBridgeFees, tt.wantBridgeFees)
			}
			if !reflect.DeepEqual(gotAfterAmount, tt.wantAfterAmount) {
				t.Errorf("deductToBitcoinBridgeFees() gotAfterAmount = %v, want %v", gotAfterAmount, tt.wantAfterAmount)
			}
		})
	}
}

func Test_deductTONToEVMBridgeFees(t *testing.T) {
	type args struct {
		amount        *big.Int
		bridgeFeeRate *big.Int
	}
	tests := []struct {
		name            string
		args            args
		wantBridgeFees  *big.Int
		wantAfterAmount *big.Int
	}{

		// amount * 30 /10000 + 1e8
		{
			name: "t-1",
			args: args{
				amount:        convertDecimal(new(big.Int).SetUint64(4950000), params.USDTDecimalNumberOfTON, params.FixedDecimalNumber), // 4.95
				bridgeFeeRate: big.NewInt(30),
			},
			wantBridgeFees:  big.NewInt(101485000),
			wantAfterAmount: big.NewInt(393515000),
		},
		{
			name: "t-2",
			args: args{
				amount:        convertDecimal(new(big.Int).SetUint64(10000000), params.USDTDecimalNumberOfTON, params.FixedDecimalNumber), // 10
				bridgeFeeRate: big.NewInt(30),
			},
			wantBridgeFees:  big.NewInt(103000000),
			wantAfterAmount: big.NewInt(897000000),
		},
		{
			name: "t-3",
			args: args{
				amount:        convertDecimal(new(big.Int).SetUint64(1286000000), params.USDTDecimalNumberOfTON, params.FixedDecimalNumber), // 1286
				bridgeFeeRate: big.NewInt(30),
			},
			wantBridgeFees:  big.NewInt(485800000),    // 4.858
			wantAfterAmount: big.NewInt(128114200000), // 1281.142
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBridgeFees, gotAfterAmount := deductTONToEVMBridgeFees(tt.args.amount, tt.args.bridgeFeeRate)
			if !reflect.DeepEqual(gotBridgeFees, tt.wantBridgeFees) {
				t.Errorf("deductTONToEVMBridgeFees() gotBridgeFees = %v, want %v", gotBridgeFees, tt.wantBridgeFees)
			}
			if !reflect.DeepEqual(gotAfterAmount, tt.wantAfterAmount) {
				t.Errorf("deductTONToEVMBridgeFees() gotAfterAmount = %v, want %v", gotAfterAmount, tt.wantAfterAmount)
			}
		})
	}
}

func Test_deductBitcoinToEVMBridgeFees(t *testing.T) {
	type args struct {
		amount        *big.Int
		bridgeFeeRate *big.Int
	}
	tests := []struct {
		name            string
		args            args
		wantBridgeFees  *big.Int
		wantAfterAmount *big.Int
	}{
		// amount * 70 /10000 + 1050
		{
			name: "t-1",
			args: args{
				amount:        convertDecimal(new(big.Int).SetUint64(50000), params.BTCDecimalNumber, params.FixedDecimalNumber),
				bridgeFeeRate: big.NewInt(70),
			},
			wantBridgeFees:  big.NewInt(1400),
			wantAfterAmount: big.NewInt(48600),
		},
		{
			name: "t-2",
			args: args{
				amount:        convertDecimal(new(big.Int).SetUint64(1280000), params.BTCDecimalNumber, params.FixedDecimalNumber),
				bridgeFeeRate: big.NewInt(70),
			},
			wantBridgeFees:  big.NewInt(10010),
			wantAfterAmount: big.NewInt(1269990),
		},
		{
			name: "t-3",
			args: args{
				amount:        convertDecimal(new(big.Int).SetUint64(432560000), params.BTCDecimalNumber, params.FixedDecimalNumber),
				bridgeFeeRate: big.NewInt(70),
			},
			wantBridgeFees:  big.NewInt(3028970),
			wantAfterAmount: big.NewInt(429531030),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBridgeFees, gotAfterAmount := deductBitcoinToEVMBridgeFees(tt.args.amount, tt.args.bridgeFeeRate)
			if !reflect.DeepEqual(gotBridgeFees, tt.wantBridgeFees) {
				t.Errorf("deductBitcoinToEVMBridgeFees() gotBridgeFees = %v, want %v", gotBridgeFees, tt.wantBridgeFees)
			}
			if !reflect.DeepEqual(gotAfterAmount, tt.wantAfterAmount) {
				t.Errorf("deductBitcoinToEVMBridgeFees() gotAfterAmount = %v, want %v", gotAfterAmount, tt.wantAfterAmount)
			}
		})
	}
}

func Test_getFeeRate(t *testing.T) {
	amount := big.NewInt(100000000) //100 USDT
	bridgeFeeRate := big.NewInt(30)
	bridgeFees := new(big.Int).Mul(amount, bridgeFeeRate)
	bridgeFees = new(big.Int).Div(bridgeFees, big.NewInt(BridgeFeeRateDenominator))
	bridgeFees = new(big.Int).Add(bridgeFees, ToTONBaseTxFee)

	afterAmount := new(big.Int).Sub(amount, bridgeFees)

	t.Log(bridgeFees, afterAmount)
}

func Test_convertDecimal(t *testing.T) {
	type args struct {
		amount     *big.Int
		srcDecimal uint64
		dstDecimal uint64
	}
	tests := []struct {
		name string
		args args
		want *big.Int
	}{
		{
			name: "t-1",
			args: args{
				amount:     big.NewInt(100),
				srcDecimal: 2,
				dstDecimal: 4,
			},
			want: big.NewInt(10000),
		},
		{
			name: "t-2",
			args: args{
				amount:     big.NewInt(10000),
				srcDecimal: 4,
				dstDecimal: 2,
			},
			want: big.NewInt(100),
		},
		{
			name: "t-3",
			args: args{
				amount:     big.NewInt(100),
				srcDecimal: 2,
				dstDecimal: 2,
			},
			want: big.NewInt(100),
		},
		{
			name: "t-4",
			args: args{
				amount:     big.NewInt(50000),
				srcDecimal: 8,
				dstDecimal: 8,
			},
			want: big.NewInt(50000),
		},
		{
			name: "t-5",
			args: args{
				amount:     big.NewInt(50000),
				srcDecimal: 8,
				dstDecimal: 6,
			},
			want: big.NewInt(500),
		},
		{
			name: "t-6",
			args: args{
				amount:     big.NewInt(50000),
				srcDecimal: 6,
				dstDecimal: 8,
			},
			want: big.NewInt(5000000),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := convertDecimal(tt.args.amount, tt.args.srcDecimal, tt.args.dstDecimal); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("convertDecimal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecimalPrecision(t *testing.T) {
	// DivisionPrecision = 16(default)
	d1 := decimal.NewFromFloat(2).Div(decimal.NewFromFloat(3))
	t.Log(d1.String()) // output: "0.6666666666666667"
	d2 := decimal.NewFromFloat(2).Div(decimal.NewFromFloat(30000))
	t.Log(d2.String()) // output: "0.0000666666666667"
	d3 := decimal.NewFromFloat(20000).Div(decimal.NewFromFloat(3))
	t.Log(d3.String()) // output: "6666.6666666666666667"
	decimal.DivisionPrecision = 3
	d4 := decimal.NewFromFloat(2).Div(decimal.NewFromFloat(3))
	t.Log(d4.String()) // output: "0.667"
	decimal.DivisionPrecision = 24
	d5 := decimal.NewFromFloat(2).Div(decimal.NewFromFloat(3))
	t.Log(d5.String()) // output: "0.666666666666666666666667"
}

//func TestCompare(t *testing.T) {
//	//amount := uint64(912_12345678)
//	amount := uint64(912_12340000)
//
//	fee := calcProtocolFees(new(big.Int).SetUint64(amount), 70, 1e8) //big int
//	t.Log("fee:", fee)
//	relayAmount := decimal.NewFromUint64(amount).Sub(decimal.NewFromUint64(fee.Uint64())) // decimal
//	t.Log("relayAmount:", relayAmount)
//	relayAmountStr := unwrapFixedDecimal(relayAmount).StringFixedBank(params.WBTCDecimalNumberOfChainPool)
//	t.Log("relayAmountStr:", relayAmountStr)
//
//	t.Log(912_12340000 - 638486380)
//}

// 73859259

func Test_calcProtocolFees(t *testing.T) {
	type args struct {
		inAmountSat *big.Int
		feeRatio    uint64
	}
	tests := []struct {
		name string
		args args
		want *big.Int
	}{
		{
			name: "t-1",
			args: args{
				inAmountSat: big.NewInt(0),
				feeRatio:    30,
			},
			want: big.NewInt(0),
		},
		{
			name: "t-2",
			args: args{
				inAmountSat: big.NewInt(50000),
				feeRatio:    0,
			},
			want: big.NewInt(0),
		},
		{
			name: "t-3",
			args: args{
				inAmountSat: big.NewInt(50000),
				feeRatio:    30,
			},
			want: big.NewInt(150),
		},
		{
			name: "t-5",
			args: args{
				inAmountSat: big.NewInt(3815200),
				feeRatio:    70,
			},
			want: big.NewInt(26706), // 3815200 * 70 / 10000 = 26706.4
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calcProtocolFees(tt.args.inAmountSat, tt.args.feeRatio); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("calcProtocolFees2() = %v, want %v", got, tt.want)
			}
		})
	}
}
