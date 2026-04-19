package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// RADIUSUser represents a user configured in the RADIUS auth method.
type RADIUSUser struct {
	Username string
	Policies []string
}

// RADIUSChecker fetches RADIUS users from Vault.
type RADIUSChecker struct {
	client *http.Client
	base   string
	token  string
	mount  string
}

// NewRADIUSChecker creates a new RADIUSChecker.
func NewRADIUSChecker(client *http.Client, base, token, mount string) *RADIUSChecker {
	if mount == "" {
		mount = "radius"
	}
	return &RADIUSChecker{client: client, base: base, token: token, mount: mount}
}

// ListUsers returns all RADIUS users configured under the given mount.
func (r *RADIUSChecker) ListUsers() ([]RADIUSUser, error) {
	url := fmt.Sprintf("%s/v1/auth/%s/users?list=true", r.base, r.mount)
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
		return nil, nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("radius list users: unexpected status %d", resp.StatusCode)
	}

	var body struct {
		Data struct {
			Keys []string `json:"keys"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, err
	}

	users := make([]RADIUSUser, 0, len(body.Data.Keys))
	for _, k := range body.Data.Keys {
		users = append(users, RADIUSUser{Username: k})
	}
	return users, nil
}
