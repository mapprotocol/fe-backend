package tonrouter

import (
	"github.com/mapprotocol/fe-backend/utils"
	"testing"
)

func TestRouteAndSwap(t *testing.T) {
	type args struct {
		request *RouteAndSwapRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *RouteData
		wantErr bool
	}{
		//fromChainId=1360104473493505&
		//toChainId=1360104473493505&
		//amount=0.1&
		//slippage=500&
		//tokenInAddress=0x0000000000000000000000000000000000000000&
		//tokenOutAddress=EQBlqsm144Dq6SjbPI4jjZvA1hqTIP3CvHovbIfW_t-SCALE&
		//receiver=UQCcgIOWXxRpCWmQ8n2QLC2crtysNNjIAXzfhuqmVJEBH7Dl
		{
			name: "t-1",
			args: args{
				request: &RouteAndSwapRequest{
					FromChainID:     "1360104473493505",
					ToChainID:       "1",
					Amount:          "1",
					TokenInAddress:  "0x0000000000000000000000000000000000000000",
					TokenOutAddress: "0x0000000000000000000000000000000000000000",
					Slippage:        100,
					Receiver:        "0xf855a761f9182c4b22A04753681A1F6324Ed3449",
				},
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "t-2",
			args: args{
				request: &RouteAndSwapRequest{
					FromChainID:     "1360104473493505",
					ToChainID:       "1360104473493505",
					Amount:          "10",
					TokenInAddress:  "0x0000000000000000000000000000000000000000",
					TokenOutAddress: "EQBlqsm144Dq6SjbPI4jjZvA1hqTIP3CvHovbIfW_t-SCALE",
					Slippage:        100,
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
					FromChainID:     "1360104473493505",
					ToChainID:       "1360104473493505",
					Amount:          "10",
					TokenInAddress:  "0x0000000000000000000000000000000000000000",
					TokenOutAddress: "EQBlqsm144Dq6SjbPI4jjZvA1hqTIP3CvHovbIfW_t-SCALE",
					Slippage:        100,
					Receiver:        "UQCcgIOWXxRpCWmQ8n2QLC2crtysNNjIAXzfhuqmVJEBH7Dl",
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
				t.Errorf("RouteAndSwap() got error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("RouteAndSwap() got = %v, want %v", utils.JSON(got), tt.want)
			//}
			t.Logf("RouteAndSwap() got = %v", utils.JSON(got))
		})
	}
}
