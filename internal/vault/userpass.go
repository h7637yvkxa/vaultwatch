package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// UserpassUser represents a user in the userpass auth method.
type UserpassUser struct {
	Username string
	Policies []string
}

// UserpassChecker lists users from the userpass auth method.
type UserpassChecker struct {
	client *http.Client
	base   string
	token  string
	mount  string
}

// NewUserpassChecker creates a UserpassChecker using the given Vault client.
func NewUserpassChecker(c *Client, mount string) *UserpassChecker {
	if mount == "" {
		mount = "userpass"
	}
	return &UserpassChecker{
		client: c.http,
		base:   c.address,
		token:  c.token,
		mount:  mount,
	}
}

// ListUsers returns all users configured in the userpass mount.
func (u *UserpassChecker) ListUsers() ([]UserpassUser, error) {
	url := fmt.Sprintf("%s/v1/auth/%s/users?list=true", u.base, u.mount)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Vault-Token", u.token)

	resp, err := u.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("userpass list: unexpected status %d", resp.StatusCode)
	}

	var body struct {
		Data struct {
			Keys []string `json:"keys"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, err
	}

	users := make([]UserpassUser, 0, len(body.Data.Keys))
	for _, k := range body.Data.Keys {
		users = append(users, UserpassUser{Username: k})
	}
	return users, nil
}
