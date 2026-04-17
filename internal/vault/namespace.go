package vault

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// NamespaceEntry holds metadata for a single Vault namespace.
type NamespaceEntry struct {
	Path string
	ID   string
}

// NamespaceChecker lists namespaces from Vault.
type NamespaceChecker struct {
	address string
	token   string
	client  *http.Client
}

// NewNamespaceChecker creates a new NamespaceChecker.
func NewNamespaceChecker(address, token string, client *http.Client) *NamespaceChecker {
	if client == nil {
		client = http.DefaultClient
	}
	return &NamespaceChecker{address: address, token: token, client: client}
}

// ListNamespaces returns all namespaces under the given prefix.
func (nc *NamespaceChecker) ListNamespaces(ctx context.Context, prefix string) ([]NamespaceEntry, error) {
	url := fmt.Sprintf("%s/v1/sys/namespaces/%s?list=true", nc.address, prefix)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Vault-Token", nc.token)

	resp, err := nc.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return []NamespaceEntry{}, nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("namespace list: unexpected status %d", resp.StatusCode)
	}

	var payload struct {
		Data struct {
			Keys []string `json:"keys"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}

	entries := make([]NamespaceEntry, 0, len(payload.Data.Keys))
	for _, k := range payload.Data.Keys {
		entries = append(entries, NamespaceEntry{Path: prefix + k, ID: k})
	}
	return entries, nil
}
