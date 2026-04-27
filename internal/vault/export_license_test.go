package vault

import (
	"net/http"
	"testing"
)

// NewLicenseCheckerFromParts constructs a LicenseChecker directly from parts
// for use in external test packages.
func NewLicenseCheckerFromParts(t *testing.T, address, token string) *LicenseChecker {
	t.Helper()
	return &LicenseChecker{
		address: address,
		token:   token,
		client:  &http.Client{},
	}
}
