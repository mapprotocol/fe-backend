package mempool

import (
	"fmt"
	"io"
	"log"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/wire"
	uhttp "github.com/mapprotocol/fe-backend/utils/http"
)

type MempoolClient struct {
	baseURL string
}

func NewClient(netParams *chaincfg.Params) *MempoolClient {
	baseURL := ""
	if netParams.Net == wire.MainNet {
		baseURL = "https://mempool.space/api"
	} else if netParams.Net == wire.TestNet3 {
		baseURL = "https://mempool.space/testnet/api"
	} else if netParams.Net == chaincfg.SigNetParams.Net {
		baseURL = "https://mempool.space/signet/api"
	} else {
		log.Fatal("mempool don't support other netParams")
	}
	return &MempoolClient{
		baseURL: baseURL,
	}
}

func (c *MempoolClient) request(method, subPath string, body io.Reader) ([]byte, error) {
	url := fmt.Sprintf("%s%s", c.baseURL, subPath)
	return uhttp.Request(url, method, nil, body)
}
