package vault

import "net/http"

// NewAuthCheckerFromParts exposes the internal constructor for testing.
func NewAuthCheckerFromParts(client *http.Client, base, token string) *AuthChecker {
	return NewAuthChecker(client, base, token)
}
