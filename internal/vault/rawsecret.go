package vault

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// RawSecretEntry represents a raw secret value read from a KV v1 path.
type RawSecretEntry struct {
	Path string
	Data map[string]interface{}
}

// RawSecretChecker reads raw secrets from a KV v1 mount.
type RawSecretChecker struct {
	address string
	token   string
	client  *http.Client
}

// NewRawSecretChecker creates a new RawSecretChecker using the provided Vault client.
func NewRawSecretChecker(c *Client) *RawSecretChecker {
	return &RawSecretChecker{
		address: c.Address,
		token:   c.Token,
		client:  c.HTTP,
	}
}

// ReadRawSecret reads a KV v1 secret at the given path.
func (r *RawSecretChecker) ReadRawSecret(ctx context.Context, path string) (*RawSecretEntry, error) {
	if path == "" {
		return nil, fmt.Errorf("rawsecret: path must not be empty")
	}

	url := fmt.Sprintf("%s/v1/%s", r.address, path)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("rawsecret: build request: %w", err)
	}
	req.Header.Set("X-Vault-Token", r.token)

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("rawsecret: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("rawsecret: path %q not found", path)
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("rawsecret: unexpected status %d: %s", resp.StatusCode, string(body))
	}

	var envelope struct {
		Data map[string]interface{} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
		return nil, fmt.Errorf("rawsecret: decode response: %w", err)
	}

	return &RawSecretEntry{
		Path: path,
		Data: envelope.Data,
	}, nil
}
