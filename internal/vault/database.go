package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// DatabaseRole represents a Vault database role.
type DatabaseRole struct {
	Name           string `json:"name"`
	DBName         string `json:"db_name"`
	DefaultTTL     int    `json:"default_ttl"`
	MaxTTL         int    `json:"max_ttl"`
	CreationStmts  string `json:"creation_statements"`
}

// DatabaseChecker lists database roles from Vault.
type DatabaseChecker struct {
	address string
	token   string
	client  *http.Client
	mount   string
}

// NewDatabaseChecker creates a new DatabaseChecker.
func NewDatabaseChecker(address, token, mount string) *DatabaseChecker {
	if mount == "" {
		mount = "database"
	}
	return &DatabaseChecker{
		address: address,
		token:   token,
		client:  &http.Client{},
		mount:   mount,
	}
}

// ListDatabaseRoles returns all database roles from the given mount.
func (c *DatabaseChecker) ListDatabaseRoles() ([]DatabaseRole, error) {
	url := fmt.Sprintf("%s/v1/%s/roles?list=true", c.address, c.mount)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("X-Vault-Token", c.token)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return []DatabaseRole{}, nil
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

	roles := make([]DatabaseRole, 0, len(result.Data.Keys))
	for _, k := range result.Data.Keys {
		roles = append(roles, DatabaseRole{Name: k})
	}
	return roles, nil
}
