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

var senderAddress common.Address
var senderPrivateKey *ecdsa.PrivateKey

type Transactor struct {
	endpoint           string
	client             *ethclient.Client
	address            common.Address
	privateKey         *ecdsa.PrivateKey
	gasLimitMultiplier float64
	chainPoolContract  common.Address
}

func InitTransactor(privateKey string) {
	var err error
	senderPrivateKey, err = crypto.ToECDSA(common.FromHex(privateKey))
	if err != nil {
		panic(err)
	}
	senderAddress = crypto.PubkeyToAddress(senderPrivateKey.PublicKey)
}

func NewTransactor(endpoint, chainPoolContract string, gasLimitMultiplier float64) (*Transactor, error) {
	client, err := ethclient.Dial(endpoint)
	if err != nil {
		return nil, err
	}

	multiplier := 1.0
	if gasLimitMultiplier > 1 {
		multiplier = gasLimitMultiplier
	}

	chainPoolAddress := common.HexToAddress(chainPoolContract)

	return &Transactor{
		endpoint:           endpoint,
		client:             client,
		address:            senderAddress,
		privateKey:         senderPrivateKey,
		gasLimitMultiplier: multiplier,
		chainPoolContract:  chainPoolAddress,
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
