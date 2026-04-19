package vault

import (
	"encoding/json"
	"fmt"
	"net/http"

	hashivault "github.com/hashicorp/vault/api"
)

// OIDCRole represents a Vault OIDC role.
type OIDCRole struct {
	Name        string
	BoundAudiences []string `json:"bound_audiences"`
	UserClaim   string   `json:"user_claim"`
	TTL         string   `json:"ttl"`
}

// OIDCChecker lists OIDC roles from Vault.
type OIDCChecker struct {
	client *hashivault.Client
}

// NewOIDCChecker returns a new OIDCChecker.
func NewOIDCChecker(client *hashivault.Client) *OIDCChecker {
	return &OIDCChecker{client: client}
}

// ListRoles returns all OIDC roles configured in Vault.
func (o *OIDCChecker) ListRoles() ([]OIDCRole, error) {
	req := o.client.NewRequest(http.MethodList, "/v1/auth/oidc/role")
	resp, err := o.client.RawRequest(req)
	if err != nil {
		return nil, fmt.Errorf("oidc list roles: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("oidc list roles: unexpected status %d", resp.StatusCode)
	}

	var body struct {
		Data struct {
			Keys []string `json:"keys"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, fmt.Errorf("oidc list roles decode: %w", err)
	}

	roles := make([]OIDCRole, 0, len(body.Data.Keys))
	for _, k := range body.Data.Keys {
		roles = append(roles, OIDCRole{Name: k})
	}
	return roles, nil
}
