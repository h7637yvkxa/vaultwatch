package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// SSHRole represents a Vault SSH secret engine role.
type SSHRole struct {
	Name     string
	KeyType  string `json:"key_type"`
	TTL      string `json:"ttl"`
	MaxTTL   string `json:"max_ttl"`
	AllowedUsers string `json:"allowed_users"`
}

// SSHChecker lists SSH roles from a given mount path.
type SSHChecker struct {
	addr   string
	token  string
	client *http.Client
}

// NewSSHChecker creates a new SSHChecker.
func NewSSHChecker(addr, token string) *SSHChecker {
	return &SSHChecker{addr: addr, token: token, client: &http.Client{}}
}

// ListRoles returns all SSH roles under the given mount (e.g. "ssh").
func (s *SSHChecker) ListRoles(mount string) ([]SSHRole, error) {
	url := fmt.Sprintf("%s/v1/%s/roles?list=true", s.addr, mount)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Vault-Token", s.token)

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ssh list roles: unexpected status %d", resp.StatusCode)
	}

	var body struct {
		Data struct {
			Keys []string `json:"keys"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, err
	}

	roles := make([]SSHRole, 0, len(body.Data.Keys))
	for _, k := range body.Data.Keys {
		roles = append(roles, SSHRole{Name: k})
	}
	return roles, nil
}
