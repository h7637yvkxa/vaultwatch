package vault

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// SealStatus holds the seal/unseal state of the Vault cluster.
type SealStatus struct {
	Sealed      bool      `json:"sealed"`
	Initialized bool      `json:"initialized"`
	T           int       `json:"t"`
	N           int       `json:"n"`
	Progress    int       `json:"progress"`
	Version     string    `json:"version"`
	ClusterName string    `json:"cluster_name"`
	CheckedAt   time.Time `json:"-"`
}

// SealChecker queries the Vault seal status endpoint.
type SealChecker struct {
	baseURL    string
	httpClient *http.Client
}

// NewSealChecker creates a SealChecker using the provided Vault address.
func NewSealChecker(address string, httpClient *http.Client) *SealChecker {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 10 * time.Second}
	}
	return &SealChecker{baseURL: address, httpClient: httpClient}
}

// Check fetches the current seal status from Vault.
func (s *SealChecker) Check(ctx context.Context) (*SealStatus, error) {
	url := fmt.Sprintf("%s/v1/sys/seal-status", s.baseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("seal check: build request: %w", err)
	}
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("seal check: request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("seal check: unexpected status %d", resp.StatusCode)
	}
	var status SealStatus
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		return nil, fmt.Errorf("seal check: decode: %w", err)
	}
	status.CheckedAt = time.Now()
	return &status, nil
}
