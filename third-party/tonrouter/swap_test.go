package tonrouter

import (
	"github.com/mapprotocol/fe-backend/utils"
	"testing"
)

func TestRouteAndSwap(t *testing.T) {
	type args struct {
		request *BridgeSwapRequest
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
				request: &BridgeSwapRequest{
					Sender:   "0xf855a761f9182c4b22A04753681A1F6324Ed3449",
					Receiver: "0xf855a761f9182c4b22A04753681A1F6324Ed3449",
					Hash:     "", // invalid hash
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "t-2",
			args: args{
				request: &BridgeSwapRequest{
					Sender:   "0xf855a761f9182c4b22A04753681A1F6324Ed3449", // "invalid sender 0xf855a761f9182c4b22A04753681A1F6324Ed3449"
					Receiver: "0xf855a761f9182c4b22A04753681A1F6324Ed3449",
					Hash:     "0x0ba6fdc2da0f791be0952adf02bec836fb6ff7528d3051ebb4d715581b3c09d8",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "t-3",
			args: args{
				request: &BridgeSwapRequest{
					Sender:       "UQAiqP-qy6O3Fbe8aGNTMz_aubZDMBqTqE0Nwf7qfyj1nB9c",
					Receiver:     "0xf855a761f9182c4b22A04753681A1F6324Ed3449",
					FeeCollector: "UQAiqP-qy6O3Fbe8aGNTMz_aubZDMBqTqE0Nwf7qfyj1nB9c",
					FeeRatio:     "200",
					Hash:         "0x9bafbe79da0214dc9464ba62546ae86d68dd483d342128219dc1d201d61f812d",
				},
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BridgeSwap(tt.args.request)
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
