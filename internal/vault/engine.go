package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// EngineInfo holds metadata about a secrets engine mount.
type EngineInfo struct {
	Path        string
	Type        string
	Description string
	Local       bool
	SealWrap    bool
}

// EngineChecker lists enabled secrets engines from Vault.
type EngineChecker struct {
	client *http.Client
	base   string
	token  string
}

// NewEngineChecker creates an EngineChecker using the provided Vault client.
func NewEngineChecker(c *Client) *EngineChecker {
	return &EngineChecker{
		client: c.http,
		base:   c.address,
		token:  c.token,
	}
}

// ListEngines returns all secrets engine mounts.
func (e *EngineChecker) ListEngines() ([]EngineInfo, error) {
	req, err := http.NewRequest(http.MethodGet, e.base+"/v1/sys/mounts", nil)
	if err != nil {
		return nil, fmt.Errorf("engine: build request: %w", err)
	}
	req.Header.Set("X-Vault-Token", e.token)

	resp, err := e.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("engine: request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("engine: unexpected status %d", resp.StatusCode)
	}

	var raw map[string]struct {
		Type        string `json:"type"`
		Description string `json:"description"`
		Local       bool   `json:"local"`
		SealWrap    bool   `json:"seal_wrap"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("engine: decode: %w", err)
	}

	var engines []EngineInfo
	for path, info := range raw {
		engines = append(engines, EngineInfo{
			Path:        path,
			Type:        info.Type,
			Description: info.Description,
			Local:       info.Local,
			SealWrap:    info.SealWrap,
		})
	}
	return engines, nil
}
