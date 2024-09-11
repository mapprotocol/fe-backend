package tx

import (
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/mapprotocol/fe-backend/utils"
	"math/big"
	"reflect"
	"testing"
)

func TestTransactor_Deliver(t1 *testing.T) {
	endpoint := "https://rpc.maplabs.io"
	privateKey, err := crypto.ToECDSA(common.FromHex("0xe5cb2b60ddaa2c087543a3f730340fdb7b6491071f4351a6e860075d2fd8b570")) // keeper private key
	if err != nil {
		t1.Fatal(err)
	}

	client, err := ethclient.Dial(endpoint)
	if err != nil {
		t1.Fatal(err)
	}

	//amount, ok := new(big.Int).SetString("1621889354640251400000", 10)
	amount, ok := new(big.Int).SetString("1621000000000000000000", 10)
	if !ok {
		t1.Fatal(err)
	}

	type fields struct {
		endpoint           string
		client             *ethclient.Client
		address            common.Address
		privateKey         *ecdsa.PrivateKey
		gasLimitMultiplier float64
		chainPoolContract  common.Address
	}
	type args struct {
		orderID     [32]byte
		token       common.Address
		amount      *big.Int
		receiver    common.Address
		fee         *big.Int
		feeReceiver common.Address
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    common.Hash
		wantErr bool
	}{
		{
			name: "test",
			fields: fields{
				endpoint:           endpoint,
				client:             client,
				address:            common.HexToAddress("0x848bAc8cc5a68c381fdb95FB4D8C0C4bb52DE2AB"), // keeper address
				privateKey:         privateKey,
				gasLimitMultiplier: 1.2,
				chainPoolContract:  common.HexToAddress("0x63d13711AFcD6Ddd07154c19F020b490C25bDCD0"), // fe router contract
			},
			args: args{
				orderID: utils.Uint64ToByte32(1000000000002),
				//token:       common.HexToAddress("0xc2132d05d31c914a87c6611c10748aeb04b58e8f"), // USDT of polygon
				//token:       common.HexToAddress("0x1BFD67037B42Cf73acF2047067bd4F2C47D9BfD6"), // WBTC of polygon
				token:       common.HexToAddress("0x33daba9618a75a7aff103e53afe530fbacf4a3dd"), // USDT of MAP
				amount:      amount,
				receiver:    common.HexToAddress("0x2E9B4be739453cdDbB3641FB61052BA46873D41f"),
				fee:         big.NewInt(0),
				feeReceiver: common.Address{},
			},
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &Transactor{
				endpoint:           tt.fields.endpoint,
				client:             tt.fields.client,
				address:            tt.fields.address,
				privateKey:         tt.fields.privateKey,
				gasLimitMultiplier: tt.fields.gasLimitMultiplier,
				chainPoolContract:  tt.fields.chainPoolContract,
			}
			got, err := t.Deliver(tt.args.orderID, tt.args.token, tt.args.amount, tt.args.receiver, tt.args.fee, tt.args.feeReceiver)
			if (err != nil) != tt.wantErr {
				t1.Errorf("Deliver() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("Deliver() got = %v, want %v", got, tt.want)
			}
		})
	}
}
