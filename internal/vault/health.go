package vault

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// HealthStatus represents the result of a Vault health check.
type HealthStatus struct {
	Initialized bool
	Sealed      bool
	Standby     bool
	Version     string
	ClusterName string
}

// Checker performs health checks against a Vault instance.
type Checker struct {
	client     *http.Client
	vaultAddr  string
}

// NewChecker creates a new Checker for the given Vault address.
func NewChecker(vaultAddr string, timeout time.Duration) *Checker {
	return &Checker{
		client:    &http.Client{Timeout: timeout},
		vaultAddr: vaultAddr,
	}
}

// Check performs a health check against the Vault /v1/sys/health endpoint.
// It returns a HealthStatus and an error if the request fails or Vault is
// sealed / uninitialized.
func (c *Checker) Check(ctx context.Context) (*HealthStatus, error) {
	url := fmt.Sprintf("%s/v1/sys/health", c.vaultAddr)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("building health request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing health request: %w", err)
	}
	defer resp.Body.Close()

	status := &HealthStatus{}

	switch resp.StatusCode {
	case http.StatusOK:
		status.Initialized = true
		status.Sealed = false
		status.Standby = false
	case http.StatusTooManyRequests:
		// 429 — active node, performance standby
		status.Initialized = true
		status.Sealed = false
		status.Standby = true
	case http.StatusNotImplemented:
		// 501 — not initialised
		status.Initialized = false
		return status, fmt.Errorf("vault is not initialized (HTTP %d)", resp.StatusCode)
	case http.StatusServiceUnavailable:
		// 503 — sealed
		status.Initialized = true
		status.Sealed = true
		return status, fmt.Errorf("vault is sealed (HTTP %d)", resp.StatusCode)
	default:
		return status, fmt.Errorf("unexpected health status code: %d", resp.StatusCode)
	}

	return status, nil
}
