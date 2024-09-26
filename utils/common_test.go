package utils

import (
	"fmt"
	"github.com/shopspring/decimal"
	"math"
	"math/big"
	"reflect"
	"testing"
)

func TestTrimHexPrefix(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "t-1",
			args: args{
				s: "0x1234567890abcdef",
			},
			want: "1234567890abcdef",
		},
		{
			name: "t-2",
			args: args{
				s: "0X1234567890abcdef",
			},
			want: "1234567890abcdef",
		},
		{
			name: "t-3",
			args: args{
				s: "1234567890abcdef",
			},
			want: "1234567890abcdef",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TrimHexPrefix(tt.args.s); got != tt.want {
				t.Errorf("TrimHexPrefix() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCompare(t *testing.T) {
	btcAmount := "0.00049650"
	btcDecimal := 1e8
	wbtcDecimal := 1e18

	// big float
	amountBigFloat, ok := new(big.Float).SetString(btcAmount)
	if !ok {
		t.Fatal("Float SetString failed")
	}

	amountBigFloat = new(big.Float).Mul(amountBigFloat, big.NewFloat(btcDecimal))
	fmt.Println("============================== amount float: ", amountBigFloat.Text('f', -1))

	amountInt, _ := amountBigFloat.Int(nil)
	fmt.Println("============================== amount: ", amountInt.String())

	// big rat
	amountBigRat, ok := new(big.Rat).SetString(btcAmount)
	if !ok {
		t.Fatal("Rat SetString failed")
	}
	amountBigRat = new(big.Rat).Mul(amountBigRat, new(big.Rat).SetUint64(uint64(btcDecimal)))
	amountInt = amountBigRat.Num()
	fmt.Println("============================== amount: ", amountInt.String())

	// decimal
	relayAmount, err := decimal.NewFromString(btcAmount)
	if err != nil {
		t.Fatal("Decimal SetString failed")
	}

	relayAmount = relayAmount.Mul(decimal.NewFromUint64(uint64(wbtcDecimal)))
	amountInt = relayAmount.BigInt()
	fmt.Println("============================== amount: ", amountInt.String())
}

func Test333(t *testing.T) {
	btcAmount := "0.00100000"
	btcDecimal := 8

	amountBigFloat, ok := new(big.Float).SetString(btcAmount)
	if !ok {
		t.Fatal("failed to set big float")
	}

	exp := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(btcDecimal)), nil)

	// big float
	amount := new(big.Float).Mul(amountBigFloat, new(big.Float).SetInt(exp))
	t.Log(amount.Text('f', -1))
	amountBigInt, _ := amount.Int(nil)
	t.Log(amountBigInt.String())

	// decimal
	amountDecimal, err := decimal.NewFromString(btcAmount)
	if err != nil {
		t.Fatal("Decimal SetString failed")
	}

	amountDecimal = amountDecimal.Mul(decimal.NewFromBigInt(exp, 0))
	t.Log(amountDecimal.String())
	amountInt := amountDecimal.BigInt()
	t.Log(amountInt.String())
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
			name: "t-1",
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

func deductFees(amount, feeRate *big.Int) (feeAmount, afterAmount *big.Int) {
	//feeRate = new(big.Int).Quo(feeRate, big.NewInt(10000))
	feeAmount = new(big.Int).Mul(amount, feeRate)
	feeAmount = new(big.Int).Div(feeAmount, big.NewInt(10000))
	afterAmount = new(big.Int).Sub(amount, feeAmount)
	return feeAmount, afterAmount
}

func TestDecimal(t *testing.T) {
	//value := big.NewInt(100)
	//dec := decimal.NewFromBigInt(value, 0)
	//t.Log("value: ", value)
	//t.Log("value: ", dec.BigInt())
	//
	//btc := decimal.NewFromBigInt(big.NewInt(50000), 0).Div(decimal.NewFromFloat(params.BTCDecimal))
	//t.Log("btc: ", btc.String())
	//
	//t.Log("sats: ", decimal.NewFromFloat(0.00100000).Mul(decimal.NewFromFloat(params.BTCDecimal)))
	//t.Log("sats: ", 0.00100000*params.BTCDecimal)
	//
	//exp := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(8)), nil)
	//amount := new(big.Float).Mul(big.NewFloat(0.00100000), new(big.Float).SetInt(exp))
	//amountBigInt, acc := amount.Int(nil)
	//t.Log("amount: ", amountBigInt.String(), acc)

	//t.Log("=====: ", decimal.NewFromFloat(12345.1234567891).StringFixedBank(-1))
	//t.Log("=====: ", decimal.NewFromFloat(12345.1234567891).StringFixedBank(-2))
	//t.Log("=====: ", decimal.NewFromFloat(12345.1234567891).StringFixedBank(0))

	f1 := 12345.1234567895
	t.Log("===== f1: ", decimal.NewFromFloat(f1).StringFixedBank(8))
	t.Log("===== f1: ", decimal.NewFromFloat(f1).StringFixed(8))
	t.Log("===== f1: ", decimal.NewFromFloat(f1).String())

	f2 := 12345.123456685
	t.Log("===== f2: ", decimal.NewFromFloat(f2).StringFixedBank(8))
	t.Log("===== f2: ", decimal.NewFromFloat(f2).StringFixed(8))
	t.Log("===== f2: ", decimal.NewFromFloat(f2).String())

	f3 := 12345.1234
	t.Log("===== f3: ", decimal.NewFromFloat(f3).StringFixedBank(8))
	t.Log("===== f3: ", decimal.NewFromFloat(f3).StringFixed(8))
	t.Log("===== f3: ", decimal.NewFromFloat(f3).String())

	f4 := 12345.1235
	t.Log("===== f4: ", decimal.NewFromFloat(f4).StringFixedBank(8))
	t.Log("===== f4: ", decimal.NewFromFloat(f4).StringFixed(8))
	t.Log("===== f4: ", decimal.NewFromFloat(f4).String())
}

func TestBase64ToHex(t *testing.T) {
	type args struct {
		base64Str string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "t-1",
			args: args{
				base64Str: "qD2s/oD4Qt2QS3JXsRMb9A95/y6KLXKRu9Yhg+RD504=",
			},
			want:    "a83dacfe80f842dd904b7257b1131bf40f79ff2e8a2d7291bbd62183e443e74e",
			wantErr: false,
		},
		{
			name: "t-1",
			args: args{
				base64Str: "yPVVUz6qd0nZZMzifAgNggXzWpRxBgvuSuGol4cp/tw=",
			},
			want:    "c8f555533eaa7749d964cce27c080d8205f35a9471060bee4ae1a8978729fedc",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Base64ToHex(tt.args.base64Str)
			if (err != nil) != tt.wantErr {
				t.Errorf("Base64ToHex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Base64ToHex() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUint64ToByte32(t *testing.T) {
	type args struct {
		num uint64
	}
	tests := []struct {
		name string
		args args
		want [32]byte
	}{
		{
			name: "t-1",
			args: args{
				num: 0,
			},
			want: [32]byte{0},
		},
		{
			name: "t-2",
			args: args{
				num: 1,
			},
			want: [32]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		},
		{
			name: "t-3",
			args: args{
				num: math.MaxUint64 - 1,
			},
			want: [32]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 255, 255, 255, 255, 255, 254},
		},
		{
			name: "t-4",
			args: args{
				num: math.MaxUint32,
			},
			want: [32]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 255, 255},
		},
		{
			name: "t-5",
			args: args{
				num: math.MaxUint32 - 10,
			},
			want: [32]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 255, 245},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Uint64ToByte32(tt.args.num); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Uint64ToByte32() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBytesToUint64(t *testing.T) {
	type args struct {
		bs []byte
	}
	tests := []struct {
		name string
		args args
		want uint64
	}{
		{
			name: "t-1",
			args: args{
				bs: []byte{0},
			},
			want: 0,
		},
		{
			name: "t-2",
			args: args{
				bs: []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			},
			want: 1,
		},
		{
			name: "t-3",
			args: args{
				bs: []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 255, 255, 255, 255, 255, 254},
			},
			want: math.MaxUint64 - 1,
		},
		{
			name: "t-4",
			args: args{
				bs: []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 255, 255},
			},
			want: math.MaxUint32,
		},
		{
			name: "t-5",
			args: args{
				bs: []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 255, 245},
			},
			want: math.MaxUint32 - 10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BytesToUint64(tt.args.bs); got != tt.want {
				t.Errorf("BytesToUint64() = %v, want %v", got, tt.want)
			}
		})
	}
}
