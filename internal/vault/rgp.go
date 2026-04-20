package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// RGPPolicy represents a Role Governing Policy entry.
type RGPPolicy struct {
	Name            string   `json:"name"`
	EnforcementLevel string  `json:"enforcement_level"`
	Paths           []string `json:"paths"`
}

// RGPChecker lists RGP policies from Vault.
type RGPChecker struct {
	address string
	token   string
	client  *http.Client
}

// NewRGPChecker creates a new RGPChecker.
func NewRGPChecker(address, token string) *RGPChecker {
	return &RGPChecker{address: address, token: token, client: &http.Client{}}
}

// ListRGPPolicies returns all RGP policies from Vault.
func (r *RGPChecker) ListRGPPolicies() ([]RGPPolicy, error) {
	url := fmt.Sprintf("%s/v1/sys/policies/rgp?list=true", r.address)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("X-Vault-Token", r.token)

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return []RGPPolicy{}, nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var result struct {
		Data struct {
			Keys []string `json:"keys"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	policies := make([]RGPPolicy, 0, len(result.Data.Keys))
	for _, k := range result.Data.Keys {
		policies = append(policies, RGPPolicy{Name: k})
	}
	return policies, nil
}
