package vault

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// EGPPolicy represents an Endpoint Governing Policy in Vault.
type EGPPolicy struct {
	Name            string   `json:"name"`
	Paths           []string `json:"paths"`
	EnforcementLevel string  `json:"enforcement_level"`
}

// EGPChecker lists EGP policies from Vault.
type EGPChecker struct {
	client *http.Client
	base   string
	token  string
}

// NewEGPChecker creates a new EGPChecker.
func NewEGPChecker(client *http.Client, base, token string) *EGPChecker {
	return &EGPChecker{client: client, base: base, token: token}
}

// ListEGPPolicies returns all EGP policies from Vault.
func (e *EGPChecker) ListEGPPolicies(ctx context.Context) ([]EGPPolicy, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, e.base+"/v1/sys/policies/egp?list=true", nil)
	if err != nil {
		return nil, fmt.Errorf("egp: build request: %w", err)
	}
	req.Header.Set("X-Vault-Token", e.token)

	resp, err := e.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("egp: request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("egp: unexpected status %d", resp.StatusCode)
	}

	var body struct {
		Data struct {
			Keys []string `json:"keys"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, fmt.Errorf("egp: decode: %w", err)
	}

	policies := make([]EGPPolicy, 0, len(body.Data.Keys))
	for _, k := range body.Data.Keys {
		policies = append(policies, EGPPolicy{Name: k})
	}
	return policies, nil
}
