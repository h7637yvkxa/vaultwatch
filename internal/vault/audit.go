package vault

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// AuditEntry represents a summary of a Vault audit log device.
type AuditEntry struct {
	Path        string    `json:"path"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	Enabled     bool      `json:"enabled"`
	CheckedAt   time.Time `json:"checked_at"`
}

// AuditChecker fetches enabled audit devices from Vault.
type AuditChecker struct {
	client *Client
}

// NewAuditChecker creates a new AuditChecker.
func NewAuditChecker(c *Client) *AuditChecker {
	return &AuditChecker{client: c}
}

type auditDeviceResponse struct {
	Data map[string]struct {
		Type        string `json:"type"`
		Description string `json:"description"`
	} `json:"data"`
}

// ListAuditDevices returns all enabled audit devices.
func (a *AuditChecker) ListAuditDevices(ctx context.Context) ([]AuditEntry, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		a.client.address+"/v1/sys/audit", nil)
	if err != nil {
		return nil, fmt.Errorf("audit: build request: %w", err)
	}
	req.Header.Set("X-Vault-Token", a.client.token)

	resp, err := a.client.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("audit: request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("audit: unexpected status %d", resp.StatusCode)
	}

	var parsed auditDeviceResponse
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return nil, fmt.Errorf("audit: decode: %w", err)
	}

	now := time.Now().UTC()
	var entries []AuditEntry
	for path, dev := range parsed.Data {
		entries = append(entries, AuditEntry{
			Path:        path,
			Type:        dev.Type,
			Description: dev.Description,
			Enabled:     true,
			CheckedAt:   now,
		})
	}
	return entries, nil
}
