package vault

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// MFAMethod represents a single MFA method configured in Vault.
type MFAMethod struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

// MFAChecker lists MFA methods from Vault.
type MFAChecker struct {
	address string
	token   string
	client  *http.Client
}

// NewMFAChecker creates a new MFAChecker using the provided Vault client.
func NewMFAChecker(c *Client) *MFAChecker {
	return &MFAChecker{
		address: c.Address,
		token:   c.Token,
		client:  c.HTTP,
	}
}

// ListMFAMethods returns all MFA methods configured in Vault.
func (m *MFAChecker) ListMFAMethods(ctx context.Context) ([]MFAMethod, error) {
	url := fmt.Sprintf("%s/v1/identity/mfa/method", m.address)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("mfa: build request: %w", err)
	}
	req.Header.Set("X-Vault-Token", m.token)

	resp, err := m.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("mfa: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("mfa: unexpected status %d", resp.StatusCode)
	}

	var result struct {
		Data struct {
			Keys []string `json:"keys"`
			KeyInfo map[string]MFAMethod `json:"key_info"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("mfa: decode response: %w", err)
	}

	methods := make([]MFAMethod, 0, len(result.Data.Keys))
	for _, k := range result.Data.Keys {
		if m, ok := result.Data.KeyInfo[k]; ok {
			methods = append(methods, m)
		}
	}
	return methods, nil
}
