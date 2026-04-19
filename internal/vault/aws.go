package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// AWSRole represents an AWS auth role in Vault.
type AWSRole struct {
	Name     string
	AuthType string `json:"auth_type"`
	BoundAMI string `json:"bound_ami_id"`
	Policies []string `json:"policies"`
}

// AWSChecker lists AWS auth roles from Vault.
type AWSChecker struct {
	address string
	token   string
	client  *http.Client
	mount   string
}

// NewAWSChecker creates a new AWSChecker.
func NewAWSChecker(address, token, mount string) *AWSChecker {
	if mount == "" {
		mount = "aws"
	}
	return &AWSChecker{
		address: address,
		token:   token,
		client:  &http.Client{},
		mount:   mount,
	}
}

// ListAWSRoles returns all AWS auth roles configured in Vault.
func (c *AWSChecker) ListAWSRoles() ([]AWSRole, error) {
	url := fmt.Sprintf("%s/v1/auth/%s/role?list=true", c.address, c.mount)
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
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var result struct {
		Data struct {
			Keys []string `json:"keys"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	roles := make([]AWSRole, 0, len(result.Data.Keys))
	for _, k := range result.Data.Keys {
		roles = append(roles, AWSRole{Name: k})
	}
	return roles, nil
}
