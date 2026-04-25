package vault

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// CapabilityResult holds the capabilities for a given path and token.
type CapabilityResult struct {
	Path         string
	Capabilities []string
}

// CapabilityChecker checks token capabilities against Vault paths.
type CapabilityChecker struct {
	address string
	token   string
	client  *http.Client
}

// NewCapabilityChecker creates a new CapabilityChecker.
func NewCapabilityChecker(address, token string, client *http.Client) *CapabilityChecker {
	if client == nil {
		client = &http.Client{}
	}
	return &CapabilityChecker{address: address, token: token, client: client}
}

// CheckCapabilities queries /v1/sys/capabilities-self for the given paths.
func (c *CapabilityChecker) CheckCapabilities(ctx context.Context, paths []string) ([]CapabilityResult, error) {
	if len(paths) == 0 {
		return nil, nil
	}

	body, err := json.Marshal(map[string]interface{}{"paths": paths})
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		c.address+"/v1/sys/capabilities-self",
		newStringReader(string(body)))
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("X-Vault-Token", c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("capabilities request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("capabilities returned status %d", resp.StatusCode)
	}

	var raw map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	var results []CapabilityResult
	for _, p := range paths {
		caps := []string{}
		if v, ok := raw[p]; ok {
			if arr, ok := v.([]interface{}); ok {
				for _, cap := range arr {
					if s, ok := cap.(string); ok {
						caps = append(caps, s)
					}
				}
			}
		}
		results = append(results, CapabilityResult{Path: p, Capabilities: caps})
	}
	return results, nil
}
