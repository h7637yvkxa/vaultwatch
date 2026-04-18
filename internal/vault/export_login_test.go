package vault

import (
	"net/http"
	"testing"
)

// NewLoginCheckerFromParts constructs a LoginChecker directly for testing.
func NewLoginCheckerFromParts(t *testing.T, baseURL, token string) *LoginChecker {
	t.Helper()
	return &LoginChecker{
		baseURL:    baseURL,
		token:      token,
		httpClient: &http.Client{},
	}
}
