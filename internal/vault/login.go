package vault

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// LoginInfo holds details about the current token from a login response.
type LoginInfo struct {
	ClientToken   string    `json:"client_token"`
	Accessor      string    `json:"accessor"`
	Policies      []string  `json:"policies"`
	LeaseDuration int       `json:"lease_duration"`
	Renewable     bool      `json:"renewable"`
	IssuedAt      time.Time
}

// LoginChecker fetches token info from a Vault login endpoint.
type LoginChecker struct {
	baseURL    string
	token      string
	httpClient *http.Client
}

// NewLoginChecker constructs a LoginChecker from a Vault client.
func NewLoginChecker(c *Client) *LoginChecker {
	return &LoginChecker{
		baseURL:    c.baseURL,
		token:      c.token,
		httpClient: c.httpClient,
	}
}

// LookupToken calls /v1/auth/token/lookup-self and returns LoginInfo.
func (lc *LoginChecker) LookupToken(ctx context.Context) (*LoginInfo, error) {
	url := lc.baseURL + "/v1/auth/token/lookup-self"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("login: build request: %w", err)
	}
	req.Header.Set("X-Vault-Token", lc.token)

	resp, err := lc.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("login: request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("login: unexpected status %d", resp.StatusCode)
	}

	var envelope struct {
		Data struct {
			ID            string   `json:"id"`
			Accessor      string   `json:"accessor"`
			Policies      []string `json:"policies"`
			TTL           int      `json:"ttl"`
			Renewable     bool     `json:"renewable"`
			CreationTime  int64    `json:"creation_time"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
		return nil, fmt.Errorf("login: decode: %w", err)
	}

	d := envelope.Data
	return &LoginInfo{
		ClientToken:   d.ID,
		Accessor:      d.Accessor,
		Policies:      d.Policies,
		LeaseDuration: d.TTL,
		Renewable:     d.Renewable,
		IssuedAt:      time.Unix(d.CreationTime, 0),
	}, nil
}
