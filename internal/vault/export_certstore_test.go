package vault

import "net/http"

// NewCertStoreCheckerFromParts constructs a CertStoreChecker with a custom HTTP client for testing.
func NewCertStoreCheckerFromParts(address, token, mount string, client *http.Client) *CertStoreChecker {
	return &CertStoreChecker{
		address: address,
		token:   token,
		mount:   mount,
		client:  client,
	}
}
