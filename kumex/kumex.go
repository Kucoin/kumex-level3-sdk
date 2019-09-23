package kumex

import (
	"time"

	"github.com/Kucoin/kumex-level3-sdk/kumex/http_client"
)

type KuMEX struct {
	httpClient *http_client.Client
}

func NewKuMEX(baseUrl string, skipVerifyTls bool, timeout time.Duration) *KuMEX {
	client := http_client.NewClient(baseUrl, skipVerifyTls, timeout)
	kumex := &KuMEX{
		httpClient: client,
	}
	return kumex
}
