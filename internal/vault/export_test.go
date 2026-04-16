package vault

import (
	"net/http"
	"testing"
)

// NewClientFromParts constructs a Client for testing without going through
// full config loading.
func NewClientFromParts(address, token string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &Client{address: address, token: token, http: httpClient}
}

// newVaultClient is a test helper shared across vault package tests.
func newVaultClient(t *testing.T, address string) *Client {
	t.Helper()
	return &Client{address: address, token: "test-token", http: http.DefaultClient}
}
