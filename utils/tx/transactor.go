package tx

import (
	"context"
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
)

var defaultAddress common.Address
var defaultPrivateKey *ecdsa.PrivateKey
var defaultGasLimitMultiplier = 1.2

type Transactor struct {
	endpoint           string
	client             *ethclient.Client
	address            common.Address
	privateKey         *ecdsa.PrivateKey
	gasLimitMultiplier float64
}

func InitTransactor(private string) {
	var err error
	defaultPrivateKey, err = crypto.ToECDSA(common.FromHex(private))
	if err != nil {
		panic(err)
	}
	defaultAddress = crypto.PubkeyToAddress(defaultPrivateKey.PublicKey)
}

func NewTransactor(endpoint string) (*Transactor, error) {
	client, err := ethclient.Dial(endpoint)
	if err != nil {
		return nil, err
	}

	return &Transactor{
		endpoint:   endpoint,
		client:     client,
		address:    defaultAddress,
		privateKey: defaultPrivateKey,
	}, nil
}

func (t *Transactor) NewTransactOpts() (*bind.TransactOpts, error) {
	nonce, err := t.client.PendingNonceAt(context.Background(), t.address)
	if err != nil {
		return nil, err
	}

	id, err := t.client.ChainID(context.Background())
	if err != nil {
		return nil, err
	}

	opts, err := bind.NewKeyedTransactorWithChainID(t.privateKey, id)
	if err != nil {
		return nil, err
	}

	opts.Nonce = big.NewInt(int64(nonce))
	// get gas price and gas limit from network
	opts.Context = context.Background()

	return opts, nil
}

func (t *Transactor) NewCallOpts() *bind.CallOpts {
	return &bind.CallOpts{From: t.address}
}
