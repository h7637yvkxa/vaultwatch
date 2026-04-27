package vault

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// LogicalBackend represents a mounted logical backend in Vault.
type LogicalBackend struct {
	Path        string `json:"path"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Local       bool   `json:"local"`
	SealWrap    bool   `json:"seal_wrap"`
}

// LogicalBackendChecker checks logical backends mounted in Vault.
type LogicalBackendChecker struct {
	address string
	token   string
	client  *http.Client
}

// NewLogicalBackendChecker creates a new LogicalBackendChecker.
func NewLogicalBackendChecker(address, token string) *LogicalBackendChecker {
	return &LogicalBackendChecker{
		address: address,
		token:   token,
		client:  &http.Client{},
	}
}

// ListLogicalBackends returns all logical backends mounted in Vault.
func (c *LogicalBackendChecker) ListLogicalBackends(ctx context.Context) ([]LogicalBackend, error) {
	url := fmt.Sprintf("%s/v1/sys/mounts", c.address)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("X-Vault-Token", c.token)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var result map[string]json.RawMessage
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	var backends []LogicalBackend
	for path, raw := range result {
		var entry struct {
			Type        string `json:"type"`
			Description string `json:"description"`
			Local       bool   `json:"local"`
			SealWrap    bool   `json:"seal_wrap"`
		}
		if err := json.Unmarshal(raw, &entry); err != nil {
			continue
		}
		if entry.Type == "" {
			continue
		}
		backends = append(backends, LogicalBackend{
			Path:        path,
			Type:        entry.Type,
			Description: entry.Description,
			Local:       entry.Local,
			SealWrap:    entry.SealWrap,
		})
	}
	return backends, nil
}
