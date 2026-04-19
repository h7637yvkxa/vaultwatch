package vault

import "net/http"

// NewAWSCheckerFromParts constructs an AWSChecker with a custom HTTP client for testing.
func NewAWSCheckerFromParts(address, token, mount string, client *http.Client) *AWSChecker {
	c := NewAWSChecker(address, token, mount)
	c.client = client
	return c
}
