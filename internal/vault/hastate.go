package vault

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// HAState represents the high-availability state of a Vault node.
type HAState struct {
	HAEnabled     bool   `json:"ha_enabled"`
	IsSelf        bool   `json:"is_self"`
	ActiveTime    string `json:"active_time"`
	LeaderAddress string `json:"leader_address"`
	LeaderCluster string `json:"leader_cluster_address"`
	PerfStandby   bool   `json:"performance_standby"`
}

// HAStateChecker fetches HA leader information from Vault.
type HAStateChecker struct {
	address string
	token   string
	client  *http.Client
}

// NewHAStateChecker creates a new HAStateChecker.
func NewHAStateChecker(address, token string) *HAStateChecker {
	return &HAStateChecker{
		address: address,
		token:   token,
		client:  &http.Client{},
	}
}

// Check retrieves the current HA state from Vault.
func (h *HAStateChecker) Check(ctx context.Context) (*HAState, error) {
	url := fmt.Sprintf("%s/v1/sys/leader", h.address)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("hastate: build request: %w", err)
	}
	req.Header.Set("X-Vault-Token", h.token)

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("hastate: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("hastate: unexpected status %d", resp.StatusCode)
	}

	var state HAState
	if err := json.NewDecoder(resp.Body).Decode(&state); err != nil {
		return nil, fmt.Errorf("hastate: decode response: %w", err)
	}
	return &state, nil
}
