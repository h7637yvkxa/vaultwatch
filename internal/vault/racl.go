package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// RACLEntry represents a single response-wrapping ACL rule.
type RACLEntry struct {
	Path        string   `json:"path"`
	Capabilities []string `json:"capabilities"`
}

// RACLChecker fetches ACL rules from a Vault path.
type RACLChecker struct {
	address string
	token   string
	client  *http.Client
}

// NewRACLChecker creates a new RACLChecker using the provided Vault client.
func NewRACLChecker(c *Client) *RACLChecker {
	return &RACLChecker{
		address: c.Address,
		token:   c.Token,
		client:  c.HTTP,
	}
}

// ListACLPaths returns ACL entries for the given policy name.
func (r *RACLChecker) ListACLPaths(policy string) ([]RACLEntry, error) {
	if policy == "" {
		return nil, fmt.Errorf("policy name must not be empty")
	}

	url := fmt.Sprintf("%s/v1/sys/policy/%s", r.address, policy)
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

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var result struct {
		Rules []RACLEntry `json:"rules"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return result.Rules, nil
}
