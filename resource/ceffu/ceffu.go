package ceffu

import (
	"github.com/mapprotocol/ceffu-go/client"
	"net/http"
	"time"
)

var ceffuClient *client.Client

func Init(domain, apiKey, apiKeySecret string) {
	opts := client.Options{
		Domain: domain,
		HttpClient: &http.Client{
			Timeout: 20 * time.Second,
		},
	}
	cli, err := client.New(apiKey, apiKeySecret, opts)
	if err != nil {
		panic(err)
	}
	ceffuClient = &cli
}

func GetClient() *client.Client {
	return ceffuClient
}
