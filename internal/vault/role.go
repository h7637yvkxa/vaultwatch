package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// RoleEntry represents a single role returned from Vault.
type RoleEntry struct {
	Name     string
	Path     string
	TokenTTL int `json:"token_ttl"`
	MaxTTL   int `json:"token_max_ttl"`
}

// RoleChecker lists roles for a given auth mount path.
type RoleChecker struct {
	client *http.Client
	base   string
	token  string
}

// NewRoleChecker creates a RoleChecker using the provided Vault client.
func NewRoleChecker(c *Client) *RoleChecker {
	return &RoleChecker{
		client: c.HTTP,
		base:   c.Address,
		token:  c.Token,
	}
}

// ListRoles returns all roles under the given auth mount (e.g. "approle").
func (r *RoleChecker) ListRoles(mount string) ([]RoleEntry, error) {
	url := fmt.Sprintf("%s/v1/auth/%s/role?list=true", r.base, mount)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Vault-Token", r.token)

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return []RoleEntry{}, nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("vault returned status %d", resp.StatusCode)
	}

	var body struct {
		Data struct {
			Keys []string `json:"keys"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, err
	}

	roles := make([]RoleEntry, 0, len(body.Data.Keys))
	for _, k := range body.Data.Keys {
		roles = append(roles, RoleEntry{Name: k, Path: fmt.Sprintf("auth/%s/role/%s", mount, k)})
	}
	return roles, nil
}
