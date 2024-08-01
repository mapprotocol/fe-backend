package butter

import (
	"github.com/mapprotocol/fe-backend/utils"
	"testing"
)

func TestRoute(t *testing.T) {
	type args struct {
		request *RouteRequest
	}
	tests := []struct {
		name    string
		args    args
		want    []*RouteResponseData
		wantErr bool
	}{
		{
			name: "t-1",
			args: args{
				request: &RouteRequest{
					FromChainID:     "1",
					ToChainID:       "137",
					Amount:          "23",
					TokenInAddress:  "0x0000000000000000000000000000000000000000",
					TokenOutAddress: "0x0000000000000000000000000000000000000000",
					Type:            "exactIn",
					Slippage:        100,
				},
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "t-2",
			args: args{
				request: &RouteRequest{
					FromChainID:     "1",
					ToChainID:       "22776",
					Amount:          "18",
					TokenInAddress:  "0x0000000000000000000000000000000000000000",
					TokenOutAddress: "0x0000000000000000000000000000000000000000",
					Type:            "exactIn",
					Slippage:        200,
				},
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "t-3",
			args: args{
				request: &RouteRequest{
					FromChainID:     "22776",
					ToChainID:       "137",
					Amount:          "1",
					TokenInAddress:  "0x0000000000000000000000000000000000000000",
					TokenOutAddress: "0x0000000000000000000000000000000000000000",
					Type:            "exactIn",
					Slippage:        200,
				},
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "t-4",
			args: args{
				request: &RouteRequest{
					FromChainID:     "22776",
					ToChainID:       "137",
					Amount:          "1",
					TokenInAddress:  "0x0000000000000000000000000000000000000000",
					TokenOutAddress: "0x0000000000000000000000000000000000000000",
					Type:            "exactIn",
					Slippage:        200,
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Route(tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("Route() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("Route() got = %v, want %v", got, tt.want)
			//}
			t.Logf("Route() got = %v", utils.JSON(got))
		})
	}
}
