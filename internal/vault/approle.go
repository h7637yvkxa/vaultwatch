package vault

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// AppRoleEntry holds metadata about a Vault AppRole.
type AppRoleEntry struct {
	RoleID   string `json:"role_id"`
	Name     string `json:"-"`
	BindSecretID bool `json:"bind_secret_id"`
	TokenTTL int    `json:"token_ttl"`
	TokenMaxTTL int `json:"token_max_ttl"`
}

// AppRoleChecker lists AppRoles from a Vault approle auth mount.
type AppRoleChecker struct {
	address   string
	token     string
	mountPath string
	client    *http.Client
}

// NewAppRoleChecker creates a new AppRoleChecker.
func NewAppRoleChecker(address, token, mountPath string) *AppRoleChecker {
	return &AppRoleChecker{
		address:   address,
		token:     token,
		mountPath: mountPath,
		client:    &http.Client{},
	}
}

// ListAppRoles returns all AppRole names and their metadata.
func (a *AppRoleChecker) ListAppRoles(ctx context.Context) ([]AppRoleEntry, error) {
	url := fmt.Sprintf("%s/v1/auth/%s/role?list=true", a.address, a.mountPath)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Vault-Token", a.token)

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
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

	var entries []AppRoleEntry
	for _, name := range body.Data.Keys {
		entries = append(entries, AppRoleEntry{Name: name})
	}
	return entries, nil
}
