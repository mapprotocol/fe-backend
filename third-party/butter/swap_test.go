package butter

import (
	"github.com/mapprotocol/fe-backend/utils"
	"testing"
)

func TestSwap(t *testing.T) {
	type args struct {
		request *SwapRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *TxData
		wantErr bool
	}{
		{
			name: "t-1",
			args: args{
				request: &SwapRequest{
					From:     "0xf855a761f9182c4b22A04753681A1F6324Ed3449",
					Receiver: "0xf855a761f9182c4b22A04753681A1F6324Ed3449",
					Hash:     "",
					Slippage: 100,
				},
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Swap(tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("Swap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("Swap() got = %v", utils.JSON(got))
		})
	}
}
