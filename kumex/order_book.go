package kumex

import (
	"net/http"

	"github.com/Kucoin/kumex-level3-sdk/kumex/http_client"
)

func (kumex *KuMEX) AtomicFullOrderBook(symbol string) (*http_client.Response, error) {
	return kumex.httpClient.Request(http.MethodGet, "/api/v1/level3/snapshot", map[string]string{"symbol": symbol})
}
