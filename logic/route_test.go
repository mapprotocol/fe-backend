package logic

import (
	"github.com/shopspring/decimal"
	"math/big"
	"reflect"
	"strconv"
	"testing"
)

func Test_calcBridgeAndProtocolFees(t *testing.T) {
	type args struct {
		amount          func() decimal.Decimal
		bridgeFeeRate   decimal.Decimal
		protocolFeeRate decimal.Decimal
	}
	tests := []struct {
		name                string
		args                args
		wantBridgeFeesStr   string
		wantProtocolFeesStr string
		wantAfterAmountStr  string
	}{
		{
			name: "t-1",
			args: args{
				amount:          func() decimal.Decimal { return decimal.NewFromFloat(100) },
				bridgeFeeRate:   decimal.NewFromFloat(70.0 / 10000.0),
				protocolFeeRate: decimal.NewFromFloat(70.0 / 10000.0),
			},
			wantBridgeFeesStr:   "0.7",
			wantProtocolFeesStr: "0.7",
			wantAfterAmountStr:  "98.6",
		},
		//{
		//	name: "t-2",
		//	args: args{
		//		amount:          func() decimal.Decimal { return decimal.NewFromFloat(100.123456789123456789123456789123456789) }, // replace to decimal.NewFromString()
		//		bridgeFeeRate:   decimal.NewFromFloat(70.0 / 10000.0),
		//		protocolFeeRate: decimal.NewFromFloat(70.0 / 10000.0),
		//	},
		//	wantBridgeFeesStr:   "0.700864197523864197523864197523864197523",
		//	wantProtocolFeesStr: "0.700864197523864197523864197523864197523",
		//	wantAfterAmountStr:  "98.721728394075728394075728394075728393954",
		//},

		{

			name: "t-2",
			args: args{
				amount: func() decimal.Decimal {
					amount, _ := decimal.NewFromString("100.123456789123456789123456789123456789")
					return amount
				},
				bridgeFeeRate:   decimal.NewFromFloat(70.0 / 10000.0),
				protocolFeeRate: decimal.NewFromFloat(70.0 / 10000.0),
			},
			wantBridgeFeesStr:   "0.700864197523864197523864197523864197523",
			wantProtocolFeesStr: "0.700864197523864197523864197523864197523",
			wantAfterAmountStr:  "98.721728394075728394075728394075728393954",
		},
		{

			name: "t-3",
			args: args{
				amount: func() decimal.Decimal {
					amount, _ := decimal.NewFromString("123456.123456789123456789")
					return amount
				},
				bridgeFeeRate:   decimal.NewFromFloat(70.0 / 10000.0),
				protocolFeeRate: decimal.NewFromFloat(70.0 / 10000.0),
			},
			wantBridgeFeesStr:   "864.192864197523864197523",
			wantProtocolFeesStr: "864.192864197523864197523",
			wantAfterAmountStr:  "121727.737728394075728393954",
		},
		{

			name: "t-4",
			args: args{
				amount: func() decimal.Decimal {
					amount, _ := decimal.NewFromString("10.123456789")
					return amount
				},
				bridgeFeeRate:   decimal.NewFromFloat(70.0 / 10000.0),
				protocolFeeRate: decimal.NewFromFloat(70.0 / 10000.0),
			},
			wantBridgeFeesStr:   "0.070864197523",
			wantProtocolFeesStr: "0.070864197523",
			wantAfterAmountStr:  "9.981728393954",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBridgeFeesStr, gotProtocolFeesStr, gotAfterAmountStr := calcBridgeAndProtocolFees(tt.args.amount(), tt.args.bridgeFeeRate, tt.args.protocolFeeRate)
			if gotBridgeFeesStr != tt.wantBridgeFeesStr {
				t.Errorf("calcBridgeAndProtocolFees() gotBridgeFeesStr = %v, want %v", gotBridgeFeesStr, tt.wantBridgeFeesStr)
			}
			if gotProtocolFeesStr != tt.wantProtocolFeesStr {
				t.Errorf("calcBridgeAndProtocolFees() gotProtocolFeesStr = %v, want %v", gotProtocolFeesStr, tt.wantProtocolFeesStr)
			}
			if gotAfterAmountStr != tt.wantAfterAmountStr {
				t.Errorf("calcBridgeAndProtocolFees() gotAfterAmountStr = %v, want %v", gotAfterAmountStr, tt.wantAfterAmountStr)
			}
		})
	}
}

func Test_calcBridgeAndProtocolFees1(t *testing.T) {
	type args struct {
		amount          *big.Float
		bridgeFeeRate   *big.Float
		protocolFeeRate *big.Float
	}
	tests := []struct {
		name             string
		args             args
		wantBridgeFees   *big.Float
		wantProtocolFees *big.Float
		wantAfterAmount  *big.Float
	}{
		{
			name: "t-1",
			args: args{
				amount:          big.NewFloat(100),
				bridgeFeeRate:   big.NewFloat(70),
				protocolFeeRate: big.NewFloat(70),
			},
			wantBridgeFees:   big.NewFloat(0.7),
			wantProtocolFees: big.NewFloat(0.7),
			wantAfterAmount:  big.NewFloat(98.6),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBridgeFees, gotProtocolFees, gotAfterAmount := calcBridgeAndProtocolFees1(tt.args.amount, tt.args.bridgeFeeRate, tt.args.protocolFeeRate)
			if !reflect.DeepEqual(gotBridgeFees, tt.wantBridgeFees) {
				t.Errorf("calcBridgeAndProtocolFees() gotBridgeFees = %v, want %v", gotBridgeFees, tt.wantBridgeFees)
			}
			if !reflect.DeepEqual(gotProtocolFees, tt.wantProtocolFees) {
				t.Errorf("calcBridgeAndProtocolFees() gotProtocolFees = %v, want %v", gotProtocolFees, tt.wantProtocolFees)
			}
			_ = gotAfterAmount
			//	if !reflect.DeepEqual(gotAfterAmount, tt.wantAfterAmount) {
			//		t.Errorf("calcBridgeAndProtocolFees() gotAfterAmount = %v, want %v", gotAfterAmount, tt.wantAfterAmount)
			//	}
		})
	}
}

func Test_calcBridgeAndProtocolFees2(t *testing.T) {
	type args struct {
		amount          *big.Rat
		bridgeFeeRate   *big.Rat
		protocolFeeRate *big.Rat
	}
	tests := []struct {
		name                string
		args                args
		wantBridgeFeesStr   string
		wantProtocolFeesStr string
		wantAfterAmountStr  string
	}{
		{
			name: "t-1",
			args: args{
				amount:          big.NewRat(100, 1),
				bridgeFeeRate:   big.NewRat(70, 1),
				protocolFeeRate: big.NewRat(70, 1),
			},
			wantBridgeFeesStr:   "0.7",
			wantProtocolFeesStr: "0.7",
			wantAfterAmountStr:  "98.6",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBridgeFeesStr, gotProtocolFeesStr, gotAfterAmountStr := calcBridgeAndProtocolFees2(tt.args.amount, tt.args.bridgeFeeRate, tt.args.protocolFeeRate)
			if gotBridgeFeesStr != tt.wantBridgeFeesStr {
				t.Errorf("calcBridgeAndProtocolFees2() gotBridgeFeesStr = %v, want %v", gotBridgeFeesStr, tt.wantBridgeFeesStr)
			}
			if gotProtocolFeesStr != tt.wantProtocolFeesStr {
				t.Errorf("calcBridgeAndProtocolFees2() gotProtocolFeesStr = %v, want %v", gotProtocolFeesStr, tt.wantProtocolFeesStr)
			}
			if gotAfterAmountStr != tt.wantAfterAmountStr {
				t.Errorf("calcBridgeAndProtocolFees2() gotAfterAmountStr = %v, want %v", gotAfterAmountStr, tt.wantAfterAmountStr)
			}
		})
	}
}

func TestDecimal(t *testing.T) {
	s := "3012345678.141592653589793238462643383279502884197169399375105820974944592307816406286208998628034825342117067982148086513282306647093844609550582231725359408128481117450284102701938521105559644622948954930381964428810975665933446128475648233786783165271201909145648566923460348610454326648213393607260249141273724587006606315588174881520920962829254091715364367892590360011330530548820466521384146951"

	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("==================== f ====================")
	t.Log(f)

	f1, err := decimal.NewFromString(s)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("==================== f1 ====================")
	t.Log(f1.String())
	t.Log(f1.Float64())
	t.Log(f1.StringFixed(390))
	t.Log(f1.StringFixed(0))
	t.Log(f1.StringFixed(10))
	t.Log(f1.StringFixedBank(390))
	t.Log(f1.StringFixedBank(0))
	t.Log(f1.StringFixedBank(10))
	t.Log(f1.BigFloat().Text('f', -1))

	f2 := decimal.NewFromFloat(12345678.123456789123456789123456789)
	t.Log("==================== f2 ====================")
	t.Log(f2.String())
	t.Log(f2.Float64())
	t.Log(f2.StringFixed(18))
	t.Log(f2.StringFixedBank(18))
	t.Log(f2.BigFloat().Text('f', -1)) //

	f3 := decimal.NewFromFloat(1.79769313486231570814527423731704356798070e+308)
	t.Log("==================== f3 ====================")
	t.Log(f3.String())
}

func TestCalcFees(t *testing.T) {
	amount, err := decimal.NewFromString("123456.123456789123456789") // 700864197523864197523864197523864197523
	if err != nil {
		t.Fatal(err)
	}
	//bridgeFeeRate := decimal.NewFromFloat(0.007)   // 0.700864197523864197523864197523864197523
	//protocolFeeRate := decimal.NewFromFloat(0.007) // 0.700864197523864197523864197523864197523

	bridgeFeeRate := decimal.NewFromFloat(0.007)     // 0.700864197523864197523864197523864197523
	protocolFeeRate := decimal.NewFromFloat(0.00007) // 0.700864197523864197523864197523864197523

	bridgeFees := amount.Mul(bridgeFeeRate)
	protocolFees := amount.Mul(protocolFeeRate)
	afterAmount := amount.Sub(bridgeFees).Sub(protocolFees)

	t.Log(bridgeFees.String())
	t.Log(protocolFees.String())
	t.Log(afterAmount.String())
}
