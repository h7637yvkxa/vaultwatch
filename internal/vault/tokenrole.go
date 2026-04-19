package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// TokenRole represents a Vault token role.
type TokenRole struct {
	Name            string   `json:"name"`
	AllowedPolicies []string `json:"allowed_policies"`
	Orphan          bool     `json:"orphan"`
	Renewable       bool     `json:"renewable"`
	ExplicitMaxTTL  int      `json:"explicit_max_ttl"`
	Period          int      `json:"period"`
}

// TokenRoleChecker lists token roles from Vault.
type TokenRoleChecker struct {
	address string
	token   string
	client  *http.Client
}

// NewTokenRoleChecker creates a new TokenRoleChecker.
func NewTokenRoleChecker(address, token string) *TokenRoleChecker {
	return &TokenRoleChecker{address: address, token: token, client: http.DefaultClient}
}

// ListTokenRoles returns all token roles defined in Vault.
func (c *TokenRoleChecker) ListTokenRoles() ([]TokenRole, error) {
	url := fmt.Sprintf("%s/v1/auth/token/roles?list=true", c.address)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Vault-Token", c.token)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
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
		return nil, err
	}

	roles := make([]TokenRole, 0, len(result.Data.Keys))
	for _, k := range result.Data.Keys {
		roles = append(roles, TokenRole{Name: k})
	}
	return roles, nil
}
