package task

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"testing"
)

func TestEncodeButterData(t *testing.T) {
	type args struct {
		initiator  common.Address
		dstToken   common.Address
		swapData   []byte
		bridgeData []byte
		feeData    []byte
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "t-1",
			args: args{
				initiator:  common.HexToAddress("0x22776"),
				dstToken:   common.HexToAddress("0xc2132d05d31c914a87c6611c10748aeb04b58e8f"),
				swapData:   nil,
				bridgeData: nil,
				feeData:    nil,
			},
			want:    "0x0000000000000000000000000000000000000000000000000000000000022776000000000000000000000000c2132d05d31c914a87c6611c10748aeb04b58e8f00000000000000000000000000000000000000000000000000000000000000a000000000000000000000000000000000000000000000000000000000000000c000000000000000000000000000000000000000000000000000000000000000e0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EncodeButterData(tt.args.initiator, tt.args.dstToken, tt.args.swapData, tt.args.bridgeData, tt.args.feeData)
			if (err != nil) != tt.wantErr {
				t.Errorf("EncodeButterData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("EncodeButterData() got = %v, want %v", got, tt.want)
			//}
			if hexutil.Encode(got) != tt.want {
				t.Errorf("EncodeButterData() got = %v, want %v", hexutil.Encode(got), tt.want)
			}
		})
	}
}

func TestName2(t *testing.T) {
	t.Log(Initiator)
}
