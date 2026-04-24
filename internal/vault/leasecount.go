package vault

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// LeaseCountResult holds the result of a lease count check.
type LeaseCountResult struct {
	Total   int
	ByMount map[string]int
}

// LeaseCountChecker queries Vault for current lease counts.
type LeaseCountChecker struct {
	address string
	token   string
	client  *http.Client
}

// NewLeaseCountChecker creates a new LeaseCountChecker.
func NewLeaseCountChecker(address, token string) *LeaseCountChecker {
	return &LeaseCountChecker{
		address: address,
		token:   token,
		client:  &http.Client{},
	}
}

type leaseCountResponse struct {
	Data struct {
		LeaseCount     int            `json:"lease_count"`
		CountPerMount  map[string]int `json:"count_per_mount"`
	} `json:"data"`
}

// Count retrieves the current lease count summary from Vault.
func (c *LeaseCountChecker) Count(ctx context.Context) (*LeaseCountResult, error) {
	url := fmt.Sprintf("%s/v1/sys/leases/count", c.address)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("leasecount: build request: %w", err)
	}
	req.Header.Set("X-Vault-Token", c.token)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("leasecount: request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("leasecount: unexpected status %d", resp.StatusCode)
	}

	var body leaseCountResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, fmt.Errorf("leasecount: decode: %w", err)
	}

	return &LeaseCountResult{
		Total:   body.Data.LeaseCount,
		ByMount: body.Data.CountPerMount,
	}, nil
}
