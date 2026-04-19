package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// LDAPGroup represents a Vault LDAP group mapping.
type LDAPGroup struct {
	Name     string   `json:"name"`
	Policies []string `json:"policies"`
}

// LDAPChecker lists LDAP group mappings from Vault.
type LDAPChecker struct {
	client *http.Client
	base   string
	token  string
	mount  string
}

// NewLDAPChecker creates a new LDAPChecker.
func NewLDAPChecker(client *http.Client, base, token, mount string) *LDAPChecker {
	if mount == "" {
		mount = "ldap"
	}
	return &LDAPChecker{client: client, base: base, token: token, mount: mount}
}

// ListGroups returns all LDAP group mappings.
func (c *LDAPChecker) ListGroups() ([]LDAPGroup, error) {
	url := fmt.Sprintf("%s/v1/auth/%s/groups?list=true", c.base, c.mount)
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
		return nil, fmt.Errorf("ldap list groups: unexpected status %d", resp.StatusCode)
	}

	var result struct {
		Data struct {
			Keys []string `json:"keys"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	groups := make([]LDAPGroup, 0, len(result.Data.Keys))
	for _, k := range result.Data.Keys {
		groups = append(groups, LDAPGroup{Name: k})
	}
	return groups, nil
}
