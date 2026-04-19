package vault

import (
	"net/http"
	"testing"
)

// NewUserpassCheckerFromParts constructs a UserpassChecker for testing.
func NewUserpassCheckerFromParts(t *testing.T, baseURL, token, mount string) *UserpassChecker {
	t.Helper()
	return &UserpassChecker{
		client: &http.Client{},
		base:   baseURL,
		token:  token,
		mount:  mount,
	}
}
