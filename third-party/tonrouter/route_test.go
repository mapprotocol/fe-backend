package tonrouter

import (
	"fmt"
	"github.com/mapprotocol/fe-backend/utils"
	"math/big"
	"reflect"
	"testing"
)

func TestMain(m *testing.M) {
	//endpoint = os.Getenv("BUTTER_ENDPOINT")
	endpoint = "https://ton-router-test.chainservice.io"
	fmt.Println("============================== ", endpoint)
	m.Run()
}

func TestRoute(t *testing.T) {
	type args struct {
		request *RouteRequest
	}
	tests := []struct {
		name    string
		args    args
		want    []*RouteData
		wantErr bool
	}{
		{
			name: "t-1",
			args: args{
				request: &RouteRequest{
					TokenInAddress:  "0x0000000000000000000000000000000000000000",
					TokenOutAddress: "EQBlqsm144Dq6SjbPI4jjZvA1hqTIP3CvHovbIfW_t-SCALE",
					Amount:          "10",
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
					Amount:          "10",
					TokenInAddress:  "0x000000000000000000000000000000", // invalid tokenInAddress
					TokenOutAddress: "EQBlqsm144Dq6SjbPI4jjZvA1hqTIP3CvHovbIfW_t-SCALE",
					Slippage:        100,
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "t-3",
			args: args{
				request: &RouteRequest{
					TokenInAddress:  "0x0000000000000000000000000000000000000000",
					TokenOutAddress: "0x0000000000000000000000000000000000000000", // tokenInAddress and tokenOutAddress is same
					Amount:          "10",
					Slippage:        100,
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "t-4",
			args: args{
				request: &RouteRequest{
					TokenInAddress:  "0x0000000000000000000000000000000000000000",
					TokenOutAddress: "0x0000000000000000000000000000000000000001", // invalid tokenOutAddress
					Amount:          "10",
					Slippage:        100,
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "t-5",
			args: args{
				request: &RouteRequest{
					TokenInAddress:  "0x0000000000000000000000000000000000000000",
					TokenOutAddress: "EQCxE6mUtQJKFnGfaROTKOt1lZbDiiX1kCixRv7Nw2Id_sDs",
					Amount:          "10",
					Slippage:        1000,
				},
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "t-6",
			args: args{
				request: &RouteRequest{
					Amount:          "10",
					TokenInAddress:  "0x0000000000000000000000000000000000000000",
					TokenOutAddress: "EQCxE6mUtQJKFnGfaROTKOt1lZbDiiX1kCixRv7Nw2Id_sDs",
					Slippage:        100,
				},
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "t-7",
			args: args{
				request: &RouteRequest{
					TokenInAddress:  "0x0000000000000000000000000000000000000000",
					TokenOutAddress: "EQCxE6mUtQJKFnGfaROTKOt1lZbDiiX1kCixRv7Nw2Id_sDs",
					Amount:          "10",
					Slippage:        5000, //
				},
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "t-8",
			args: args{
				request: &RouteRequest{
					TokenInAddress:  "0x0000000000000000000000000000000000000000",
					TokenOutAddress: "EQCxE6mUtQJKFnGfaROTKOt1lZbDiiX1kCixRv7Nw2Id_sDs",
					Amount:          "9.87",
					Slippage:        100000, // invalid slippage(<5000) 100000
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

func TestBridgeRoute(t *testing.T) {
	type args struct {
		request *BridgeRouteRequest
	}
	tests := []struct {
		name    string
		args    args
		want    []*RouteData
		wantErr bool
	}{
		{
			name: "t-1",
			args: args{
				request: &BridgeRouteRequest{
					ToChainID:       "56",
					TokenInAddress:  "EQBlqsm144Dq6SjbPI4jjZvA1hqTIP3CvHovbIfW_t-SCALE",
					TokenOutAddress: "0x0000000000000000000000000000000000000000",
					Amount:          "10",
					TonSlippage:     100,
					Slippage:        300,
				},
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "t-2",
			args: args{
				request: &BridgeRouteRequest{
					ToChainID:       "137",
					TokenInAddress:  "0x000000000000000000000000000000", // invalid tokenInAddress (length: 32)
					TokenOutAddress: "EQBlqsm144Dq6SjbPI4jjZvA1hqTIP3CvHovbIfW_t-SCALE",
					Amount:          "10",
					TonSlippage:     100,
					Slippage:        300,
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "t-3",
			args: args{
				request: &BridgeRouteRequest{
					ToChainID:       "1",
					TokenInAddress:  "0x0000000000000000000000000000000000000000",
					TokenOutAddress: "0x0000000000000000000000000000000000000000",
					Amount:          "100",
					TonSlippage:     100,
					Slippage:        300,
				},
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "t-4",
			args: args{
				request: &BridgeRouteRequest{
					ToChainID:       "22776",
					TokenInAddress:  "0x0000000000000000000000000000000000000000",
					TokenOutAddress: "0x0000000000000000000000000000000000000001",
					Amount:          "1.56",
					TonSlippage:     100,
					Slippage:        300,
				},
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "t-5",
			args: args{
				request: &BridgeRouteRequest{
					ToChainID:       "22776",
					TokenInAddress:  "0x0000000000000000000000000000000000000000",
					TokenOutAddress: "EQCxE6mUtQJKFnGfaROTKOt1lZbDiiX1kCixRv7Nw2Id_sDs", // invalid tokenOutAddress
					Amount:          "0.0089",
					TonSlippage:     100,
					Slippage:        300,
				},
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "t-6",
			args: args{
				request: &BridgeRouteRequest{
					ToChainID:       "1",
					TokenInAddress:  "EQCxE6mUtQJKFnGfaROTKOt1lZbDiiX1kCixRv7Nw2Id_sDs",
					TokenOutAddress: "0x0000000000000000000000000000000000000000",
					Amount:          "1.2",
					TonSlippage:     100,
					Slippage:        300,
				},
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "t-7",
			args: args{
				request: &BridgeRouteRequest{
					ToChainID:       "1",
					TokenInAddress:  "0x0000000000000000000000000000000000000000",
					TokenOutAddress: "0x0000000000000000000000000000000000000000",
					Amount:          "1",
					TonSlippage:     100,
					Slippage:        300,
				},
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "t-8",
			args: args{
				request: &BridgeRouteRequest{
					ToChainID:       "1",
					TokenInAddress:  "EQCxE6mUtQJKFnGfaROTKOt1lZbDiiX1kCixRv7Nw2Id_sDs",
					TokenOutAddress: "0x0000000000000000000000000000000000000000",
					Amount:          "0.0089",
					TonSlippage:     3000,
					Slippage:        9000, // invalid slippage(<5000) 9000
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BridgeRoute(tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("BridgeRoute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("BridgeRoute() got = %v", utils.JSON(got))
		})
	}
}

func TestGetRouteAmountOut(t *testing.T) {
	type args struct {
		hash string
	}
	tests := []struct {
		name    string
		args    args
		want    *big.Float
		wantErr bool
	}{
		{
			name: "t-1",

			args: args{
				hash: "0x9ba8a9525ef4658c24541cf889089a9cbb9251cdd446bfdd4178ace32fe42bc7",
			},
			want:    big.NewFloat(0.128110533),
			wantErr: false,
		},
		{
			name: "t-2",

			args: args{
				hash: "0xb5e62d3729298b58def20c7474dfd0bf756481ad1b6f25efcb17cd2b93623fc3",
			},
			want:    big.NewFloat(0.128110164),
			wantErr: false,
		},
		{
			name: "t-3",

			args: args{
				hash: "0x0",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetRouteAmountOut(tt.args.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRouteAmountOut() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetRouteAmountOut() got = %v, want %v", got, tt.want)
			}
		})
	}
}
