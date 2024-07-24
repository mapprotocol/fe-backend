package tonrouter

import (
	"testing"
)

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
					Amount:          "0.1",
					OrderID:         1,
					Receiver:        "UQAiqP-qy6O3Fbe8aGNTMz_aubZDMBqTqE0Nwf7qfyj1nB9c",
					Slippage:        100,
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
			t.Log("got: ", got)
		})
	}
}
