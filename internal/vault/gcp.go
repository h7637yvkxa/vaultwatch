package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// GCPRole represents a GCP secrets engine role.
type GCPRole struct {
	Name        string
	SecretType  string `json:"secret_type"`
	Project     string `json:"project"`
	Bindings    string `json:"bindings"`
}

// GCPChecker lists GCP roles from a Vault GCP secrets engine.
type GCPChecker struct {
	client *http.Client
	base   string
	token  string
	mount  string
}

// NewGCPChecker creates a GCPChecker using the provided Vault client.
func NewGCPChecker(c *Client, mount string) *GCPChecker {
	return &GCPChecker{
		client: c.http,
		base:   c.addr,
		token:  c.token,
		mount:  mount,
	}
}

// ListRoles returns all GCP roles under the configured mount.
func (g *GCPChecker) ListRoles() ([]GCPRole, error) {
	url := fmt.Sprintf("%s/v1/%s/roles?list=true", g.base, g.mount)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Vault-Token", g.token)

	resp, err := g.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("gcp list roles: unexpected status %d", resp.StatusCode)
	}

	var body struct {
		Data struct {
			Keys []string `json:"keys"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, err
	}

	roles := make([]GCPRole, 0, len(body.Data.Keys))
	for _, k := range body.Data.Keys {
		roles = append(roles, GCPRole{Name: k})
	}
	return roles, nil
}
