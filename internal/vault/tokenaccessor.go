package vault

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// TokenAccessorEntry represents a token accessor returned by Vault.
type TokenAccessorEntry struct {
	Accessor    string `json:"accessor"`
	CreationTime int64  `json:"creation_time"`
	DisplayName string `json:"display_name"`
	Policies    []string `json:"policies"`
	TTL         int    `json:"ttl"`
}

// TokenAccessorChecker lists token accessors from Vault.
type TokenAccessorChecker struct {
	address string
	token   string
	client  *http.Client
}

// NewTokenAccessorChecker creates a new TokenAccessorChecker.
func NewTokenAccessorChecker(address, token string) *TokenAccessorChecker {
	return &TokenAccessorChecker{
		address: address,
		token:   token,
		client:  &http.Client{},
	}
}

// ListTokenAccessors returns all token accessors from the auth/token/accessors path.
func (c *TokenAccessorChecker) ListTokenAccessors(ctx context.Context) ([]TokenAccessorEntry, error) {
	url := fmt.Sprintf("%s/v1/auth/token/accessors", c.address)
	req, err := http.NewRequestWithContext(ctx, "LIST", url, nil)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("X-Vault-Token", c.token)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("list token accessors: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var body struct {
		Data struct {
			Keys []string `json:"keys"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	entries := make([]TokenAccessorEntry, 0, len(body.Data.Keys))
	for _, k := range body.Data.Keys {
		entries = append(entries, TokenAccessorEntry{Accessor: k})
	}
	return entries, nil
}
