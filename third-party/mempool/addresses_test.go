package mempool

import (
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"testing"
)

func TestListUnspent(t *testing.T) {
	// https://mempool.space/signet/api/address/tb1p8lh4np5824u48ppawq3numsm7rss0de4kkxry0z70dcfwwwn2fcspyyhc7/utxo
	netParams := &chaincfg.SigNetParams
	client := NewClient(netParams)
	address, _ := btcutil.DecodeAddress("tb1p8lh4np5824u48ppawq3numsm7rss0de4kkxry0z70dcfwwwn2fcspyyhc7", netParams)
	unspentList, err := client.ListUnspent(address)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(len(unspentList))
		for _, output := range unspentList {
			t.Log(output.Outpoint.Hash.String(), "    ", output.Outpoint.Index)
		}
	}
}

func TestMempoolClient_ListUnspent(t *testing.T) {
	type fields struct {
		network *chaincfg.Params
	}
	type args struct {
		address string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		wantErr bool
	}{
		{
			name: "t-1",
			fields: fields{
				network: &chaincfg.MainNetParams,
			},
			args: args{
				address: "bc1qgn3j5wr8kg6lmv87jk2dlu4eaz8rveeg0ylcj3",
			},
			want:    0,
			wantErr: false,
		},
		{
			name: "t-2",
			fields: fields{
				network: &chaincfg.MainNetParams,
			},
			args: args{
				address: "bc1qq7c4dk5v3wfmccthmkmc8jd5nld74e0scwtgxd",
			},
			want:    4239287,
			wantErr: false,
		},
		{
			name: "t-3",
			fields: fields{
				network: &chaincfg.MainNetParams,
			},
			args: args{
				address: "bc1q32sxnq5hecdurfzgzp5x0zh8du86v9x84wdqdx",
			},
			want:    3037655080,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewClient(tt.fields.network)
			t.Log("mempool url: ", c.baseURL)

			address, err := btcutil.DecodeAddress(tt.args.address, tt.fields.network)
			if err != nil {
				t.Fatal(err)
			}
			unspentList, err := c.ListUnspent(address)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListUnspent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("ListUnspent() got = %v, want %v", got, tt.want)
			//}
			t.Logf("ListUnspent() got = %v", unspentList)
			if len(unspentList) == 0 {
				t.Fatal("no unspent outputs")
			}

			total := int64(0)
			for _, u := range unspentList {
				total += u.Output.Value
			}
			if total != tt.want {
				t.Errorf("ListUnspent() got = %v, want %v", total, tt.want)
			}
		})
	}
}

func TestMempoolClient_Balance(t *testing.T) {
	type fields struct {
		network *chaincfg.Params
	}
	type args struct {
		address string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		wantErr bool
	}{
		{
			name: "t-1",
			fields: fields{
				network: &chaincfg.TestNet3Params,
			},
			args: args{
				address: "tb1q9l9pph9e8ds5calkf76y40gw9d37zmt6q4mnu3",
			},
			want:    7040802,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewClient(tt.fields.network)
			t.Log("mempool url: ", c.baseURL)
			address, err := btcutil.DecodeAddress(tt.args.address, tt.fields.network)
			if err != nil {
				t.Fatal(err)
			}
			got, err := c.Balance(address)
			if (err != nil) != tt.wantErr {
				t.Errorf("Balance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Balance() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSatoshiToBTC(t *testing.T) {

	type fields struct {
		network *chaincfg.Params
	}
	type args struct {
		satoshi int64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "t-1",
			args: args{
				satoshi: 7040802,
			},
			want: 0.07040802,
		},
		{
			name: "t-2",
			args: args{
				satoshi: 50000,
			},
			want: 0.0005,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := btcutil.Amount(tt.args.satoshi).ToBTC()
			if got != tt.want {
				t.Errorf("SatoshiToBTC() got = %v, want %v", got, tt.want)
			}
		})
	}
}
