package http_client

import (
	"net/http"
	"time"
)

type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

type defaultHTTPClient struct {
	client *http.Client
}

func NewDefaultHTTPClient(timeout time.Duration) HTTPClient {
	return &defaultHTTPClient{
		client: &http.Client{
			Timeout: timeout,
			Transport: &http.Transport{
				MaxIdleConnsPerHost: 100,
				IdleConnTimeout:     90 * time.Second,
				TLSHandshakeTimeout: 10 * time.Second,
			},
		},
	}
}

func (c *defaultHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return c.client.Do(req)
}
