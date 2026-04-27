package vault

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// StoredTokenEntry represents a stored token accessor entry in Vault.
type StoredTokenEntry struct {
	Accessor    string `json:"accessor"`
	DisplayName string `json:"display_name"`
	Policies    []string `json:"policies"`
	TTL         int    `json:"ttl"`
	Renewable   bool   `json:"renewable"`
}

// StoredTokenResult holds the list of stored token entries.
type StoredTokenResult struct {
	Entries []StoredTokenEntry
}

// StoredTokenChecker retrieves stored token entries from Vault.
type StoredTokenChecker struct {
	client *http.Client
	base   string
	token  string
}

// NewStoredTokenChecker creates a new StoredTokenChecker.
func NewStoredTokenChecker(client *http.Client, base, token string) *StoredTokenChecker {
	return &StoredTokenChecker{client: client, base: base, token: token}
}

// List retrieves all stored token entries from Vault.
func (c *StoredTokenChecker) List(ctx context.Context) (*StoredTokenResult, error) {
	url := fmt.Sprintf("%s/v1/auth/token/accessors", c.base)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("storedtoken: build request: %w", err)
	}
	req.Header.Set("X-Vault-Token", c.token)
	q := req.URL.Query()
	q.Set("list", "true")
	req.URL.RawQuery = q.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("storedtoken: request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return &StoredTokenResult{}, nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("storedtoken: unexpected status %d", resp.StatusCode)
	}

	var body struct {
		Data struct {
			Keys []string `json:"keys"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, fmt.Errorf("storedtoken: decode: %w", err)
	}

	entries := make([]StoredTokenEntry, 0, len(body.Data.Keys))
	for _, k := range body.Data.Keys {
		entries = append(entries, StoredTokenEntry{Accessor: k})
	}
	return &StoredTokenResult{Entries: entries}, nil
}
