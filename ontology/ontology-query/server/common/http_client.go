package common

import (
	"sync"

	"github.com/kweaver-ai/kweaver-go-lib/rest"
)

var (
	httpClientOnce sync.Once
	httpClient     rest.HTTPClient
)

func NewHTTPClient() rest.HTTPClient {
	httpClientOnce.Do(func() {
		httpClient = rest.NewHTTPClient()
	})

	return httpClient
}
