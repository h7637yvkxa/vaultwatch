package vault

import "net/http"

// NewSecretCheckerFromParts constructs a SecretChecker with a custom HTTP client for testing.
func NewSecretCheckerFromParts(client *http.Client, address, token string) *SecretChecker {
	return &SecretChecker{
		client:  client,
		address: address,
		token:   token,
	}
}
