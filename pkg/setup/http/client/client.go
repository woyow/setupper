package client

import (
	"net/http"
	"time"
)

type HTTPClient struct {
	*http.Client
}

// NewHTTPClient - .
func NewHTTPClient(cfg *Config) *HTTPClient {
	httpClient := getHTTPClient(cfg)

	return &HTTPClient{
		httpClient,
	}
}

// disallowFollowRedirectFunc - redirect wrapper.
func disallowFollowRedirectFunc(_ *http.Request, _ []*http.Request) error {
	return ErrRedirectNotAllowed
}

// getHTTPClient - Returns *http.Client.
func getHTTPClient(cfg *Config) *http.Client {
	var transport *http.Transport

	if _, ok := http.DefaultTransport.(*http.Transport); ok {
		transport = http.DefaultTransport.(*http.Transport).Clone()
	} else {
		panic(ok)
	}

	transport.MaxIdleConns = cfg.MaxIdleConnections
	transport.MaxConnsPerHost = cfg.MaxConnectionsPerHost
	transport.MaxIdleConnsPerHost = cfg.MaxIdleConnections

	var redirectFunc func(req *http.Request, via []*http.Request) error

	if cfg.AllowFollowRedirect {
		redirectFunc = nil
	} else {
		redirectFunc = disallowFollowRedirectFunc
	}

	return &http.Client{
		Timeout:       time.Duration(cfg.Timeout) * time.Second,
		Transport:     transport,
		CheckRedirect: redirectFunc,
	}
}
