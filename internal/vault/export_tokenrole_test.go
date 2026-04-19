package vault

import "net/http"

// NewTokenRoleCheckerFromParts constructs a TokenRoleChecker with a custom HTTP client for testing.
func NewTokenRoleCheckerFromParts(address, token string, client *http.Client) *TokenRoleChecker {
	return &TokenRoleChecker{address: address, token: token, client: client}
}
