package butter

import (
	"github.com/mapprotocol/fe-backend/utils"
	"reflect"
	"testing"
)

func TestRouteAndSwap(t *testing.T) {
	type args struct {
		request *RouteAndSwapRequest
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
				request: &RouteAndSwapRequest{
					FromChainID:     "22776",
					ToChainID:       "137",
					Amount:          "1",
					TokenInAddress:  "0x0000000000000000000000000000000000000000",
					TokenOutAddress: "0x0000000000000000000000000000000000000000",
					Kind:            "exactIn",
					Slippage:        "100",
					Entrance:        "Butter+",
					From:            "0xf855a761f9182c4b22A04753681A1F6324Ed3449",
					Receiver:        "0xf855a761f9182c4b22A04753681A1F6324Ed3449",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "t-2",
			args: args{
				request: &RouteAndSwapRequest{
					FromChainID:     "22776",
					ToChainID:       "137",
					Amount:          "10",
					TokenInAddress:  "0x0000000000000000000000000000000000000000",
					TokenOutAddress: "0x0000000000000000000000000000000000000000",
					Kind:            "exactIn",
					Slippage:        "100",
					Entrance:        "Butter+",
					From:            "0xf855a761f9182c4b22A04753681A1F6324Ed3449",
					Receiver:        "0xf855a761f9182c4b22A04753681A1F6324Ed3449",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "t-3",
			args: args{
				request: &RouteAndSwapRequest{
					FromChainID:     "22776",
					ToChainID:       "137",
					Amount:          "100",
					TokenInAddress:  "0x0000000000000000000000000000000000000000",
					TokenOutAddress: "0x0000000000000000000000000000000000000000",
					Kind:            "exactIn",
					Slippage:        "100",
					Entrance:        "Butter+",
					From:            "0xf855a761f9182c4b22A04753681A1F6324Ed3449",
					Receiver:        "0xf855a761f9182c4b22A04753681A1F6324Ed3449",
				},
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RouteAndSwap(tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("RouteAndSwap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RouteAndSwap() got = %v, want %v", utils.JSON(got), tt.want)
			}
		})
	}
}
