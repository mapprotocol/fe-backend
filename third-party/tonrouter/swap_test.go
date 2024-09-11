package tonrouter

import (
	"github.com/mapprotocol/fe-backend/utils"
	"log"
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
