package vault

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// PluginEntry represents a registered Vault plugin.
type PluginEntry struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Version string `json:"version"`
	Builtin bool   `json:"builtin"`
}

// PluginChecker lists registered plugins from Vault.
type PluginChecker struct {
	address string
	token   string
	client  *http.Client
}

// NewPluginChecker creates a PluginChecker using the provided Vault client.
func NewPluginChecker(address, token string, client *http.Client) *PluginChecker {
	if client == nil {
		client = http.DefaultClient
	}
	return &PluginChecker{address: address, token: token, client: client}
}

// ListPlugins returns all registered plugins across all types.
func (p *PluginChecker) ListPlugins(ctx context.Context) ([]PluginEntry, error) {
	url := fmt.Sprintf("%s/v1/sys/plugins/catalog", p.address)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("X-Vault-Token", p.token)

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var body struct {
		Data struct {
			Detailed []PluginEntry `json:"detailed"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return body.Data.Detailed, nil
}
