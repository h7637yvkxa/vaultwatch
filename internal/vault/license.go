package vault

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// LicenseInfo holds Vault Enterprise license details.
type LicenseInfo struct {
	LicenseID      string    `json:"license_id"`
	CustomerName   string    `json:"customer_name"`
	InstallationID string    `json:"installation_id"`
	IssueTime      time.Time `json:"issue_time"`
	StartTime      time.Time `json:"start_time"`
	ExpirationTime time.Time `json:"expiration_time"`
	Terminated     bool      `json:"terminated"`
	Features       []string  `json:"features"`
}

// LicenseChecker fetches Vault Enterprise license information.
type LicenseChecker struct {
	address string
	token   string
	client  *http.Client
}

// NewLicenseChecker constructs a LicenseChecker from a Vault client.
func NewLicenseChecker(c *Client) *LicenseChecker {
	return &LicenseChecker{
		address: c.Address,
		token:   c.Token,
		client:  c.HTTP,
	}
}

// GetLicense retrieves the current Vault Enterprise license.
func (lc *LicenseChecker) GetLicense(ctx context.Context) (*LicenseInfo, error) {
	url := fmt.Sprintf("%s/v1/sys/license/status", lc.address)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("license: build request: %w", err)
	}
	req.Header.Set("X-Vault-Token", lc.token)

	resp, err := lc.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("license: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("license: unexpected status %d", resp.StatusCode)
	}

	var wrapper struct {
		Data LicenseInfo `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
		return nil, fmt.Errorf("license: decode response: %w", err)
	}
	return &wrapper.Data, nil
}
