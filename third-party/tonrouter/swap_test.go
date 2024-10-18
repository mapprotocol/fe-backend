package tonrouter

import (
	"github.com/mapprotocol/fe-backend/utils"
	"log"
	"os"
	"reflect"
	"testing"
)

func TestMain(m *testing.M) {
	endpoint := os.Getenv("TON_ROUTER_ENDPOINT")
	if utils.IsEmpty(endpoint) {
		log.Fatal("FILTER_ROUTER_ENDPOINT environment variable is not set")
	}
	Domain = endpoint

	m.Run()
}

func TestBridgeSwap(t *testing.T) {
	type args struct {
		request *BridgeSwapRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *TxParams
		wantErr bool
	}{
		{
			name: "test",
			args: args{
				request: &BridgeSwapRequest{
					Amount:          "9.249750315",
					OrderID:         144115188075855873,
					Receiver:        "UQDtwSLVcwJyIYKKitlTgq4LR_MIfDlObsanmIAGJcUpxaFz",
					Slippage:        133,
					TokenOutAddress: "0x0000000000000000000000000000000000000000",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BridgeSwap(tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("BridgeSwap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Log("got: ", utils.JSON(got))
		})
	}
}

func TestBridgeStatus(t *testing.T) {
	type args struct {
		orderID uint64
	}
	tests := []struct {
		name    string
		args    args
		want    *Status
		wantErr bool
	}{
		{
			name: "t-1",
			args: args{
				orderID: 144115188075855873,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "t-2",
			args: args{
				orderID: 7349874591868649496,
			},
			want: &Status{
				Status:    "success",
				Hash:      "3629e92cfb143ff30331b42f582a7815e75a3143095a204efac625aa05b361f8",
				Symbol:    "Ton",
				AmountIn:  "9.467",
				AmountOut: "1.797299291",
			},
			wantErr: false,
		},
		{
			name: "t-3",
			args: args{
				orderID: 7349874591868649483,
			},
			want: &Status{
				Status:    "success",
				Hash:      "a17ab58ed5e6232adcf1bb57b79d99b8b52b30b6c5843d285dafdc241a90f904",
				Symbol:    "Ton",
				AmountIn:  "4.341124",
				AmountOut: "0.856524689",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BridgeStatus(tt.args.orderID)
			if (err != nil) != tt.wantErr {
				t.Errorf("BridgeStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BridgeStatus() got = %v, want %v", got, tt.want)
			}
		})
	}
}
