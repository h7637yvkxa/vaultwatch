package vault

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// ConnectionStatus holds the result of a Vault connectivity check.
type ConnectionStatus struct {
	Reachable   bool
	StatusCode  int
	ClusterName string
	Version     string
	Error       string
}

// ConnectionChecker verifies basic connectivity to Vault.
type ConnectionChecker struct {
	address string
	token   string
	client  *http.Client
}

// NewConnectionChecker creates a ConnectionChecker using the given Vault client.
func NewConnectionChecker(address, token string, client *http.Client) *ConnectionChecker {
	if client == nil {
		client = &http.Client{}
	}
	return &ConnectionChecker{address: address, token: token, client: client}
}

// Check performs a connectivity test against the Vault sys/health endpoint.
func (c *ConnectionChecker) Check(ctx context.Context) (*ConnectionStatus, error) {
	url := fmt.Sprintf("%s/v1/sys/health?standbyok=true&perfstandbyok=true", c.address)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return &ConnectionStatus{Reachable: false, Error: err.Error()}, err
	}
	if c.token != "" {
		req.Header.Set("X-Vault-Token", c.token)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return &ConnectionStatus{Reachable: false, Error: err.Error()}, nil
	}
	defer resp.Body.Close()

	var body struct {
		ClusterName string `json:"cluster_name"`
		Version     string `json:"version"`
	}
	_ = json.NewDecoder(resp.Body).Decode(&body)

	return &ConnectionStatus{
		Reachable:   true,
		StatusCode:  resp.StatusCode,
		ClusterName: body.ClusterName,
		Version:     body.Version,
	}, nil
}
