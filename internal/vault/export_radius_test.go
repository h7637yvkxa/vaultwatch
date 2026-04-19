package vault

import "net/http"

// NewRADIUSCheckerFromParts exposes constructor for white-box tests.
func NewRADIUSCheckerFromParts(client *http.Client, base, token, mount string) *RADIUSChecker {
	return NewRADIUSChecker(client, base, token, mount)
}
