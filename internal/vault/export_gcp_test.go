package vault

import "net/http"

// NewGCPCheckerFromParts constructs a GCPChecker directly for white-box testing.
func NewGCPCheckerFromParts(httpClient *http.Client, base, token, mount string) *GCPChecker {
	return &GCPChecker{
		client: httpClient,
		base:   base,
		token:  token,
		mount:  mount,
	}
}
