package tonrouter

//import (
//	"testing"
//)
//
//func TestBridge(t *testing.T) {
//	type args struct {
//		req *BridgeRequest
//	}
//	tests := []struct {
//		name    string
//		args    args
//		want    *TxParams
//		wantErr bool
//	}{
//		{
//			name: "t-1",
//			args: args{
//				req: &BridgeRequest{
//					TokenInAddress:  "0x0000000000000000000000000000000000000000",
//					TokenOutAddress: "0xc2132D05D31c914a87C6611C10748AEb04B58e8F",
//					Sender:          "UQAiqP-qy6O3Fbe8aGNTMz_aubZDMBqTqE0Nwf7qfyj1nB9c",
//					Receiver:        "0x0Eb16A9cFDf8e3A4471EF190eE63de5A24f38787",
//					Amount:          "10",
//					ToChainId:       "56",
//					Slippage:        300,
//				},
//			},
//			want:    nil,
//			wantErr: false,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got, err := Bridge(tt.args.req)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("Bridge() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			t.Logf("Bridge() got = %v", got)
//		})
//	}
//}
