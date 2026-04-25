package vault

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// UnsealKeyStatus holds information about the unseal key configuration.
type UnsealKeyStatus struct {
	SecretShares    int  `json:"secret_shares"`
	SecretThreshold int  `json:"secret_threshold"`
	PGPFingerprints []string `json:"pgp_fingerprints"`
	Nonce           string   `json:"nonce"`
	StoredShares    int  `json:"stored_shares"`
}

// UnsealKeyChecker retrieves unseal key configuration from Vault.
type UnsealKeyChecker struct {
	address string
	token   string
	client  *http.Client
}

// NewUnsealKeyChecker creates a new UnsealKeyChecker.
func NewUnsealKeyChecker(address, token string) *UnsealKeyChecker {
	return &UnsealKeyChecker{
		address: address,
		token:   token,
		client:  &http.Client{},
	}
}

// GetUnsealKeyStatus retrieves the current unseal key configuration.
func (c *UnsealKeyChecker) GetUnsealKeyStatus(ctx context.Context) (*UnsealKeyStatus, error) {
	url := fmt.Sprintf("%s/v1/sys/rekey/init", c.address)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}
	req.Header.Set("X-Vault-Token", c.token)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("requesting unseal key status: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d from unseal key endpoint", resp.StatusCode)
	}

	var status UnsealKeyStatus
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		return nil, fmt.Errorf("decoding unseal key status: %w", err)
	}
	return &status, nil
}
