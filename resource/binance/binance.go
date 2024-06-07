package binance

import (
	connector "github.com/binance/binance-connector-go"
	"net/http"
	"time"
)

const defaultHTTPTimeout = 20 * time.Second

var client *connector.Client

func Init(apiKey, secretKey, baseURL string) {
	// Initialise the client
	client = connector.NewClient(apiKey, secretKey, baseURL)
	client.HTTPClient = &http.Client{
		Timeout: defaultHTTPTimeout,
	}
}

func GetClient() *connector.Client {
	return client
}
