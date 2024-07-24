package tonclient

import (
	"context"
	"github.com/xssnick/tonutils-go/liteclient"
	"github.com/xssnick/tonutils-go/ton"
	tonwallet "github.com/xssnick/tonutils-go/ton/wallet"
	"strings"
)

var (
	client ton.APIClientWrapped
	wallet *tonwallet.Wallet
)

func Init(words, password string) {
	pool := liteclient.NewConnectionPool()

	cfg, err := liteclient.GetConfigFromUrl(context.Background(), "https://ton.org/global.config.json")
	//cfg, err := liteclient.GetConfigFromUrl(context.Background(), "https://ton.org/testnet-global.config.json")
	if err != nil {
		panic(err)
	}

	// connect to mainnet lite servers
	err = pool.AddConnectionsFromConfig(context.Background(), cfg)
	if err != nil {
		panic(err)
	}
	api := ton.NewAPIClient(pool, ton.ProofCheckPolicySecure).WithRetry()
	api.SetTrustedBlockFromConfig(cfg)
	client = api

	// seed words of account, you can generate them with any wallet or using wallet.NewSeed() method
	seed := strings.Split(words, " ")

	w, err := tonwallet.FromSeedWithPassword(api, seed, password, tonwallet.V3R2)
	if err != nil {
		panic(err)
	}
	wallet = w
}

func Client() ton.APIClientWrapped {
	return client
}

func Wallet() *tonwallet.Wallet {
	return wallet
}
