package vault

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// RequestCountResult holds Vault request count metrics.
type RequestCountResult struct {
	StartTime  string         `json:"start_time"`
	EndTime    string         `json:"end_time"`
	ByNamespace []NSReqCount  `json:"by_namespace"`
	Total       int           `json:"total"`
}

// NSReqCount holds per-namespace request counts.
type NSReqCount struct {
	NamespaceID   string `json:"namespace_id"`
	NamespacePath string `json:"namespace_path"`
	Counts        int    `json:"counts"`
}

// RequestCountChecker fetches request count activity from Vault.
type RequestCountChecker struct {
	address string
	token   string
	client  *http.Client
}

// NewRequestCountChecker creates a RequestCountChecker using the provided Vault client.
func NewRequestCountChecker(address, token string, client *http.Client) *RequestCountChecker {
	if client == nil {
		client = http.DefaultClient
	}
	return &RequestCountChecker{address: address, token: token, client: client}
}

// GetRequestCounts queries the Vault activity log endpoint for request counts.
func (r *RequestCountChecker) GetRequestCounts(ctx context.Context) (*RequestCountResult, error) {
	url := fmt.Sprintf("%s/v1/sys/internal/counters/requests", r.address)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("X-Vault-Token", r.token)

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request counts request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d from request counts endpoint", resp.StatusCode)
	}

	var wrapper struct {
		Data RequestCountResult `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
		return nil, fmt.Errorf("decode request counts: %w", err)
	}
	return &wrapper.Data, nil
}
