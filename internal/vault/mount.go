package vault

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// MountInfo holds basic information about a Vault secret engine mount.
type MountInfo struct {
	Path        string
	Type        string
	Description string
	Accessor    string
}

// MountChecker lists secret engine mounts from Vault.
type MountChecker struct {
	address string
	token   string
	client  *http.Client
}

// NewMountChecker creates a new MountChecker.
func NewMountChecker(address, token string, client *http.Client) *MountChecker {
	if client == nil {
		client = http.DefaultClient
	}
	return &MountChecker{address: address, token: token, client: client}
}

// ListMounts returns all secret engine mounts from Vault.
func (m *MountChecker) ListMounts(ctx context.Context) ([]MountInfo, error) {
	url := fmt.Sprintf("%s/v1/sys/mounts", m.address)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("X-Vault-Token", m.token)

	resp, err := m.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
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
		return nil, fmt.Errorf("decode response: %w", err)
	}

	mounts := make([]MountInfo, 0, len(raw))
	for path, info := range raw {
		mounts = append(mounts, MountInfo{
			Path:        path,
			Type:        info.Type,
			Description: info.Description,
			Accessor:    info.Accessor,
		})
	}
	return mounts, nil
}
