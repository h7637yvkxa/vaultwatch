package vault

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// MountEntry represents a single secret engine mount.
type MountEntry struct {
	Path        string
	Type        string
	Description string
	Accessor    string
}

// MountChecker lists secret engine mounts from Vault.
type MountChecker struct {
	client *Client
}

// NewMountChecker creates a MountChecker using the provided Client.
func NewMountChecker(c *Client) *MountChecker {
	return &MountChecker{client: c}
}

// ListMounts returns all mounted secret engines.
func (m *MountChecker) ListMounts(ctx context.Context) ([]MountEntry, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		m.client.address+"/v1/sys/mounts", nil)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("X-Vault-Token", m.client.token)

	resp, err := m.client.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var raw map[string]struct {
		Type        string `json:"type"`
		Description string `json:"description"`
		Accessor    string `json:"accessor"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}

	var entries []MountEntry
	for path, info := range raw {
		entries = append(entries, MountEntry{
			Path:        path,
			Type:        info.Type,
			Description: info.Description,
			Accessor:    info.Accessor,
		})
	}
	return entries, nil
}
