package mempool

import (
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"reflect"
	"testing"
)

func TestGetRawTransaction(t *testing.T) {
	//https://mempool.space/signet/api/tx/b752d80e97196582fd02303f76b4b886c222070323fb7ccd425f6c89f5445f6c/hex
	client := NewClient(&chaincfg.TestNet3Params)
	txId, _ := chainhash.NewHashFromStr("04cb8b3e453cc8eded98b50e24cdd498254223b0e97332100f06999c504b3bab")
	transaction, err := client.GetRawTransaction(txId)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(transaction.TxHash().String())
	}
}

func TestMempoolClient_GetRawTransaction(t *testing.T) {
	type fields struct {
		network *chaincfg.Params
	}
	type args struct {
		txHash string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "t-1",
			fields: fields{
				network: &chaincfg.MainNetParams,
			},
			args: args{
				txHash: "04cb8b3e453cc8eded98b50e24cdd498254223b0e97332100f06999c504b3bab",
			},
			want:    "04cb8b3e453cc8eded98b50e24cdd498254223b0e97332100f06999c504b3bab",
			wantErr: false,
		},
		{
			name: "t-2",
			fields: fields{
				network: &chaincfg.MainNetParams,
			},
			args: args{
				txHash: "6eac170d38bb9bc82359c6a94d4b0b8fa4208328501b9c251b526d9dc8cff5b9",
			},
			want:    "6eac170d38bb9bc82359c6a94d4b0b8fa4208328501b9c251b526d9dc8cff5b9",
			wantErr: false,
		},
		{
			name: "t-3",
			fields: fields{
				network: &chaincfg.TestNet3Params,
			},
			args: args{
				txHash: "c4e6f34a7c012916cbcc5b91bdab978a71f9f9e882e5ccc559cd2867abd57c74",
			},
			want:    "c4e6f34a7c012916cbcc5b91bdab978a71f9f9e882e5ccc559cd2867abd57c74",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewClient(tt.fields.network)
			t.Log("mempool url: ", c.baseURL)

			txHash, err := chainhash.NewHashFromStr(tt.args.txHash)
			if err != nil {
				t.Fatal(err)
			}
			got, err := c.GetRawTransaction(txHash)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRawTransaction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.TxHash().String() != tt.want {
				t.Errorf("GetRawTransaction() got = %v, want %v", got.TxHash().String(), tt.want)
			}
		})
	}
}

func TestMempoolClient_TransactionStatus(t *testing.T) {
	type fields struct {
		net *chaincfg.Params
	}
	type args struct {
		txHash string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *TransactionStatusResponse
		wantErr bool
	}{
		{
			name: "unconfirmed",
			fields: fields{
				net: &chaincfg.TestNet3Params,
			},
			args: args{
				txHash: "15e10745f15593a899cef391191bdd3d7c12412cc4696b7bcb669d0feadc8521",
			},
			want: &TransactionStatusResponse{
				Confirmed:   false,
				BlockHeight: 0,
				BlockHash:   "",
				BlockTime:   0,
			},
			wantErr: false,
		},
		{
			name: "confirmed",
			fields: fields{
				net: &chaincfg.TestNet3Params,
			},
			args: args{
				txHash: "01edb3c42acae9765ff9f7824831a7fe9438c4817c056cf9637a126400c559cd",
			},
			want: &TransactionStatusResponse{
				Confirmed:   true,
				BlockHeight: 2539424,
				BlockHash:   "0000000000000024bf4c584d37889f6be903257dfabf7568942cbec5da169621",
				BlockTime:   1700471227,
			},
			wantErr: false,
		},
		{
			name: "confirmed-mainnet-1",
			fields: fields{
				net: &chaincfg.MainNetParams,
			},
			args: args{
				txHash: "04fb363f62f0745ba71cb0f067212decc65c0973c61462690bcc5e88fcdd46a7",
			},
			want: &TransactionStatusResponse{
				Confirmed:   true,
				BlockHeight: 821080,
				BlockHash:   "00000000000000000001ab0ba9e1294d373b85cfe270eeac87210294b751eaa8",
				BlockTime:   1702515068,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewClient(tt.fields.net)
			t.Log("mempool url: ", c.baseURL)
			txId, err := chainhash.NewHashFromStr(tt.args.txHash)
			if err != nil {
				t.Fatal(err)
			}
			got, err := c.TransactionStatus(txId)
			if (err != nil) != tt.wantErr {
				t.Errorf("TransactionStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TransactionStatus() got = %v, want %v", got, tt.want)
			}
		})
	}
}
