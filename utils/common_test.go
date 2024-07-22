package utils

import (
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"testing"
)

func TestIsValidEvmAddress(t *testing.T) {
	type args struct {
		address string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "t-1",
			args: args{
				address: "0x0000000000000000000000000000000000000000",
			},
			want: true,
		},
		{
			name: "t-2",
			args: args{
				address: "0xffffffffffffffffffffffffffffffffffffffff",
			},
			want: true,
		},
		{
			name: "t-3",
			args: args{
				address: "0x5553915067d79D4b0B90e48d9C94DFad2715Db49",
			},
			want: true,
		},
		{
			name: "t-4",
			args: args{
				address: "0x9eaD03F7136Fc6b4bDb0780B00a1c14aE5A8B6d0",
			},
			want: true,
		},
		{
			name: "t-5",
			args: args{
				address: "0x9ead03f7136fc6b4bdb0780b00a1c14ae5a8b6d0",
			},
			want: true,
		},
		{
			name: "t-6",
			args: args{
				address: "0X9EAD03F7136FC6B4BDB0780B00A1C14AE5A8B6D0",
			},
			want: true,
		},
		{
			name: "t-5",
			args: args{
				address: "0x9eaD03F7136Fc6b4bDb0780B00a1c14aE5A8B6d", // length 41
			},
			want: false,
		},
		{
			name: "t-6",
			args: args{
				address: "0x9e51897D4062cfcc0284B0e80Ced59706284438G",
			},
			want: false,
		},
		{
			name: "t-7",
			args: args{
				address: "0xbb623c5edc60d06ddfb61f500ecbae3c3a1ec1422", // length 43
			},
			want: false,
		},

		{
			name: "t-8",
			args: args{
				address: "0e7a64057f2524ab76f19e64e8886b05abd3b5d00",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidEvmAddress(tt.args.address); got != tt.want {
				t.Errorf("IsValidEvmAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsValidEvmHash(t *testing.T) {
	type args struct {
		hash string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "t-1",
			args: args{
				hash: "0x0000000000000000000000000000000000000000000000000000000000000000",
			},
			want: true,
		},
		{
			name: "t-2",
			args: args{
				hash: "0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
			},
			want: true,
		},
		{
			name: "t-3",
			args: args{
				hash: "0XFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF",
			},
			want: true,
		},
		{
			name: "t-4",
			args: args{
				hash: "0xaa6a63c301f3abe5fb3a72dd42e1d931a613f5cd3793bef8136fd0cb01a26453",
			},
			want: true,
		},
		{
			name: "t-5",
			args: args{
				hash: "aa6a63c301f3abe5fb3a72dd42e1d931a613f5cd3793bef8136fd0cb01a26453",
			},
			want: true,
		},
		{
			name: "t-6",
			args: args{
				hash: "0xfffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffg",
			},
			want: false,
		},
		{
			name: "t-7",
			args: args{
				hash: "0c0b566a1b8b674376e72b5174c00126b3a3409a2ed120bb7e5ee4e44306f0c74",
			},
			want: false,
		},
		{
			name: "t-8",
			args: args{
				hash: "0x1234567890abcdefabcdef1234567890",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidEvmHash(tt.args.hash); got != tt.want {
				t.Errorf("IsValidEvmHash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestName(t *testing.T) {
	//t.Log(common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"))
	//t.Log(common.HexToHash("0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"))

	//t.Log(strings.ToLower("0x9eaD03F7136Fc6b4bDb0780B00a1c14aE5A8B6d0"))
	//t.Log(strings.ToUpper("0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"))

	addr, err := btcutil.DecodeAddress("tb1p6hhxg5mtghk2jm5jd0c7rhs4sm2k7hmlpncmwlr0gdv09fzmadjsl3kkfm", &chaincfg.MainNetParams)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(addr.EncodeAddress())
}

func TestIsValidBitcoinAddress(t *testing.T) {
	type args struct {
		address string
		network *chaincfg.Params
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "t-1",
			args: args{
				address: "tb1ptgz6dcvxf9gdjgk3l3th023gnf24v7zk68ul2eh72rjk9afaftssjxh6lc",
				network: &chaincfg.TestNet3Params,
			},
			want: true,
		},
		{
			name: "t-2",
			args: args{
				address: "tb1p6hhxg5mtghk2jm5jd0c7rhs4sm2k7hmlpncmwlr0gdv09fzmadjsl3kkfm",
				network: &chaincfg.TestNet3Params,
			},
			want: true,
		},
		{
			name: "t-3",
			args: args{
				address: "1LKmxxmRSzWJpRmfKeVgbqEvaj8vQaVcDn",
				network: &chaincfg.MainNetParams,
			},
			want: true,
		},
		{
			name: "t-4",
			args: args{
				address: "1CpSD4NTJPFc9ezshCkeb8XRSpwihNUgwm",
				network: &chaincfg.MainNetParams,
			},
			want: true,
		},
		{
			name: "t-5",
			args: args{
				address: "bc1qvrcsam86q4taerkclxngfjh4y38ef3m9s0ncva",
				network: &chaincfg.MainNetParams,
			},
			want: true,
		},
		{
			name: "t-6",
			args: args{
				address: "3Gr8q9gXYJVgykcLmWveJzko7yuU4RtqeB",
				network: &chaincfg.MainNetParams,
			},
			want: true,
		},
		{
			name: "t-7",
			args: args{
				address: "tb1p6hhxg5mtghk2jm5jd0c7rhs4sm2k7hmlpncmwlr0gdv09fzmadjsl3kkfm",
				network: &chaincfg.MainNetParams,
			},
			want: true,
		},
		{
			name: "t-8",
			args: args{
				address: "1p6hhxg5mtghk2jm5jd0c7rhs4sm2k7hmlpncmwlr0gdv09fzmadjsl3kkfm",
				network: &chaincfg.MainNetParams,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidBitcoinAddress(tt.args.address, tt.args.network); got != tt.want {
				t.Errorf("IsValidBitcoinAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsValidBitcoinAddress1(t *testing.T) {
	type args struct {
		address string
		network *chaincfg.Params
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "t-1",
			args: args{
				address: "tb1q7kf0tanmqw9guuhdrrhpsfxc8qyp8gxe5xpwce",
				network: &chaincfg.TestNet3Params,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidBitcoinAddress(tt.args.address, tt.args.network); got != tt.want {
				t.Errorf("IsValidBitcoinAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}
