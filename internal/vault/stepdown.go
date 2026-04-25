package vault

import (
	"context"
	"fmt"
	"net/http"
)

// StepDownResult holds the result of a Vault leader step-down request.
type StepDownResult struct {
	Success bool
	Message string
}

// StepDownChecker can request the active Vault node to step down.
type StepDownChecker struct {
	client *http.Client
	base   string
	token  string
}

// NewStepDownChecker creates a StepDownChecker using the provided Vault client.
func NewStepDownChecker(c *Client) *StepDownChecker {
	return &StepDownChecker{
		client: c.http,
		base:   c.address,
		token:  c.token,
	}
}

// StepDown sends a PUT request to /v1/sys/step-down, asking the active node
// to relinquish leadership.
func (s *StepDownChecker) StepDown(ctx context.Context) (*StepDownResult, error) {
	url := fmt.Sprintf("%s/v1/sys/step-down", s.base)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, nil)
	if err != nil {
		return nil, fmt.Errorf("stepdown: build request: %w", err)
	}
	req.Header.Set("X-Vault-Token", s.token)

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("stepdown: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent || resp.StatusCode == http.StatusOK {
		return &StepDownResult{Success: true, Message: "step-down accepted"}, nil
	}

	return &StepDownResult{
		Success: false,
		Message: fmt.Sprintf("unexpected status: %d", resp.StatusCode),
	}, nil
}
