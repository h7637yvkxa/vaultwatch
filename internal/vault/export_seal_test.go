package vault

import "net/http"

// NewSealCheckerFromParts exposes the constructor for external test packages.
func NewSealCheckerFromParts(address string, httpClient *http.Client) *SealChecker {
	return NewSealChecker(address, httpClient)
}
