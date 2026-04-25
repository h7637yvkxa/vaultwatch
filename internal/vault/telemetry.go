package vault

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// TelemetryMetric represents a single metric from Vault's telemetry endpoint.
type TelemetryMetric struct {
	Name   string             `json:"Name"`
	Labels map[string]string  `json:"Labels"`
	Value  float64            `json:"Value"`
}

// TelemetryResult holds a subset of Vault telemetry data.
type TelemetryResult struct {
	Gauges   []TelemetryMetric `json:"Gauges"`
	Counters []TelemetryMetric `json:"Counters"`
	Samples  []TelemetryMetric `json:"Samples"`
}

// TelemetryChecker fetches Vault telemetry metrics.
type TelemetryChecker struct {
	client *http.Client
	base   string
	token  string
}

// NewTelemetryChecker creates a new TelemetryChecker using the provided Vault client.
func NewTelemetryChecker(c *Client) *TelemetryChecker {
	return &TelemetryChecker{
		client: c.HTTP,
		base:   c.BaseURL,
		token:  c.Token,
	}
}

// ReadTelemetry fetches the telemetry metrics from Vault.
func (tc *TelemetryChecker) ReadTelemetry(ctx context.Context) (*TelemetryResult, error) {
	url := fmt.Sprintf("%s/v1/sys/metrics?format=json", tc.base)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("X-Vault-Token", tc.token)

	resp, err := tc.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("telemetry request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var result TelemetryResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode telemetry: %w", err)
	}
	return &result, nil
}
