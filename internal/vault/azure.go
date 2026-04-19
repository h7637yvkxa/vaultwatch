package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// AzureRole represents a Vault Azure secrets engine role.
type AzureRole struct {
	Name            string   `json:"name"`
	ApplicationObjectID string `json:"application_object_id"`
	AzureRoles      []string `json:"azure_roles"`
	TTL             string   `json:"ttl"`
	MaxTTL          string   `json:"max_ttl"`
}

// AzureChecker lists Azure roles from Vault.
type AzureChecker struct {
	base string
	token string
	client *http.Client
}

// NewAzureChecker creates an AzureChecker from a Vault client.
func NewAzureChecker(c *Client) *AzureChecker {
	return &AzureChecker{base: c.Address, token: c.Token, client: c.HTTP}
}

// ListAzureRoles returns all Azure roles from the given mount path.
func (a *AzureChecker) ListAzureRoles(mount string) ([]AzureRole, error) {
	url := fmt.Sprintf("%s/v1/%s/roles?list=true", a.base, mount)
	req, err := http.NewRequest(http.MethodGet, url, nil)
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
		return nil, fmt.Errorf("azure list roles: unexpected status %d", resp.StatusCode)
	}

	var body struct {
		Data struct {
			Keys []string `json:"keys"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, err
	}

	roles := make([]AzureRole, 0, len(body.Data.Keys))
	for _, k := range body.Data.Keys {
		roles = append(roles, AzureRole{Name: k})
	}
	return roles, nil
}
