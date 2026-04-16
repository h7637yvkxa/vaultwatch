package vault

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// PolicyInfo holds metadata about a Vault policy.
type PolicyInfo struct {
	Name      string
	Rules     string
	FetchedAt time.Time
}

// PolicyChecker fetches policy information from Vault.
type PolicyChecker struct {
	client *Client
}

// NewPolicyChecker creates a new PolicyChecker.
func NewPolicyChecker(c *Client) *PolicyChecker {
	return &PolicyChecker{client: c}
}

// ListPolicies returns all ACL policy names from Vault.
func (p *PolicyChecker) ListPolicies(ctx context.Context) ([]string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		p.client.address+"/v1/sys/policies/acl?list=true", nil)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("X-Vault-Token", p.client.token)

	resp, err := p.client.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("list policies: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("list policies: status %d: %s", resp.StatusCode, body)
	}

	var result struct {
		Data struct {
			Keys []string `json:"keys"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return result.Data.Keys, nil
}

// GetPolicy fetches a single policy by name.
func (p *PolicyChecker) GetPolicy(ctx context.Context, name string) (*PolicyInfo, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		p.client.address+"/v1/sys/policies/acl/"+name, nil)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("X-Vault-Token", p.client.token)

	resp, err := p.client.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("get policy %s: %w", name, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("get policy %s: status %d: %s", name, resp.StatusCode, body)
	}

	var result struct {
		Data struct {
			Policy string `json:"policy"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return &PolicyInfo{Name: name, Rules: result.Data.Policy, FetchedAt: time.Now()}, nil
}
