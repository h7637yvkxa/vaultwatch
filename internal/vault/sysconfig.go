package vault

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// SysConfig holds key Vault system configuration values.
type SysConfig struct {
	DefaultLeaseTTL string `json:"default_lease_ttl"`
	MaxLeaseTTL     string `json:"max_lease_ttl"`
	ForceNoCache    bool   `json:"force_no_cache"`
}

// SysConfigChecker fetches Vault system configuration.
type SysConfigChecker struct {
	address string
	token   string
	client  *http.Client
}

// NewSysConfigChecker creates a new SysConfigChecker using the provided Vault client.
func NewSysConfigChecker(c *Client) *SysConfigChecker {
	return &SysConfigChecker{
		address: c.address,
		token:   c.token,
		client:  c.http,
	}
}

// ReadSysConfig retrieves the current Vault system configuration from /v1/sys/config/state/sanitized.
func (s *SysConfigChecker) ReadSysConfig(ctx context.Context) (*SysConfig, error) {
	url := fmt.Sprintf("%s/v1/sys/config/state/sanitized", s.address)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("sysconfig: build request: %w", err)
	}
	req.Header.Set("X-Vault-Token", s.token)

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("sysconfig: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("sysconfig: unexpected status %d", resp.StatusCode)
	}

	var envelope struct {
		Data SysConfig `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
		return nil, fmt.Errorf("sysconfig: decode response: %w", err)
	}
	return &envelope.Data, nil
}
