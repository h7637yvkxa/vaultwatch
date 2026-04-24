package vault

import "net/http"

// NewLeaseCountCheckerFromParts constructs a LeaseCountChecker with a custom HTTP client.
// Used only in tests.
func NewLeaseCountCheckerFromParts(address, token string, client *http.Client) *LeaseCountChecker {
	return &LeaseCountChecker{
		address: address,
		token:   token,
		client:  client,
	}
}
