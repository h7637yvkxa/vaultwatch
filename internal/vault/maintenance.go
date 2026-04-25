package vault

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// MaintenanceStatus holds the current maintenance mode status of Vault.
type MaintenanceStatus struct {
	Enabled    bool   `json:"enabled"`
	Message    string `json:"message"`
	RequestID  string `json:"request_id"`
}

// MaintenanceChecker checks the maintenance mode status of a Vault instance.
type MaintenanceChecker struct {
	address string
	token   string
	client  *http.Client
}

// NewMaintenanceChecker creates a new MaintenanceChecker using the provided Vault client.
func NewMaintenanceChecker(c *Client) *MaintenanceChecker {
	return &MaintenanceChecker{
		address: c.address,
		token:   c.token,
		client:  c.http,
	}
}

// Check queries the Vault sys/maintenance endpoint and returns the current status.
func (m *MaintenanceChecker) Check(ctx context.Context) (*MaintenanceStatus, error) {
	url := fmt.Sprintf("%s/v1/sys/maintenance", m.address)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("building maintenance request: %w", err)
	}
	req.Header.Set("X-Vault-Token", m.token)

	resp, err := m.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("maintenance request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status from maintenance endpoint: %d", resp.StatusCode)
	}

	var result MaintenanceStatus
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decoding maintenance response: %w", err)
	}

	return &result, nil
}
