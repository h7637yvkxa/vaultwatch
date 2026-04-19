package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// TOTPKey represents a TOTP key entry in Vault.
type TOTPKey struct {
	Name      string
	AccountID string
	Issuer    string
	Period    int
}

// TOTPChecker lists TOTP keys from the totp secrets engine.
type TOTPChecker struct {
	client *http.Client
	base   string
	token  string
}

// NewTOTPChecker creates a new TOTPChecker.
func NewTOTPChecker(client *http.Client, base, token string) *TOTPChecker {
	return &TOTPChecker{client: client, base: base, token: token}
}

// ListKeys lists all TOTP keys at the given mount path.
func (c *TOTPChecker) ListKeys(mount string) ([]TOTPKey, error) {
	url := fmt.Sprintf("%s/v1/%s/keys?list=true", c.base, mount)
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
		return nil, fmt.Errorf("totp list: unexpected status %d", resp.StatusCode)
	}

	var result struct {
		Data struct {
			Keys []string `json:"keys"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	keys := make([]TOTPKey, 0, len(result.Data.Keys))
	for _, k := range result.Data.Keys {
		keys = append(keys, TOTPKey{Name: k})
	}
	return keys, nil
}
