package task

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/mapprotocol/fe-backend/utils"
	"math/big"
	"reflect"
	"testing"
)

func TestUnpackDeliverAndSwap(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *DeliverAndSwapEventParams
		wantErr bool
	}{
		{
			name: "test-1",
			args: args{
				data: common.Hex2Bytes("0x460000174876e802000000000000000000000000000000000000000000000000214a9568fd625718bff77cbdc51f05acb117c018e6d744c5ee2d69e5c0b62d94000000000000000000000000c2132d05d31c914a87c6611c10748aeb04b58e8f000000000000000000000000000000000000000000000000000000000009ff5b"),
			},
			want: &DeliverAndSwapEventParams{
				OrderId:  utils.Uint64ToByte32(5044031682654955522),
				BridgeId: [32]byte{0x21, 0x4a, 0x95, 0x68, 0xfd, 0x62, 0x57, 0x18, 0xbf, 0xf7, 0x7c, 0xbd, 0xc5, 0x1f, 0x5, 0xac, 0xb1, 0x17, 0xc0, 0x18, 0xe6, 0xd7, 0x44, 0xc5, 0xee, 0x2d, 0x69, 0xe5, 0xc0, 0xb6, 0x2d, 0x94},
				Token:    common.HexToAddress("xc2132D05D31c914a87C6611C10748AEb04B58e8F"),
				Amount:   big.NewInt(655195),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UnpackDeliverAndSwap(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnpackDeliverAndSwap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UnpackDeliverAndSwap() got = %v, want %v", got, tt.want)
			}
		})
	}
}
