package vault

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// ReplicationStatus holds the DR and performance replication state.
type ReplicationStatus struct {
	DRMode          string `json:"dr_mode"`
	PerformanceMode string `json:"performance_mode"`
	DRPrimary       bool   `json:"dr_primary"`
	PerfPrimary     bool   `json:"perf_primary"`
}

// ReplicationChecker fetches replication status from Vault.
type ReplicationChecker struct {
	address    string
	token      string
	httpClient *http.Client
}

// NewReplicationChecker creates a ReplicationChecker using the given Vault client.
func NewReplicationChecker(address, token string, hc *http.Client) *ReplicationChecker {
	if hc == nil {
		hc = http.DefaultClient
	}
	return &ReplicationChecker{address: address, token: token, httpClient: hc}
}

// Check retrieves the current replication status.
func (r *ReplicationChecker) Check(ctx context.Context) (*ReplicationStatus, error) {
	url := fmt.Sprintf("%s/v1/sys/replication/status", r.address)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("X-Vault-Token", r.token)

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var payload struct {
		Data struct {
			DR struct {
				Mode    string `json:"mode"`
				Primary bool   `json:"primary"`
			} `json:"dr"`
			Performance struct {
				Mode    string `json:"mode"`
				Primary bool   `json:"primary"`
			} `json:"performance"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return &ReplicationStatus{
		DRMode:          payload.Data.DR.Mode,
		DRPrimary:       payload.Data.DR.Primary,
		PerformanceMode: payload.Data.Performance.Mode,
		PerfPrimary:     payload.Data.Performance.Primary,
	}, nil
}
