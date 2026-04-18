package vault

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// WrappingInfo holds metadata about a wrapped token.
type WrappingInfo struct {
	Token          string        `json:"token"`
	Accessor       string        `json:"accessor"`
	TTL            time.Duration `json:"-"`
	TTLSeconds     int           `json:"ttl"`
	CreationTime   string        `json:"creation_time"`
	CreationPath   string        `json:"creation_path"`
}

// WrappingChecker looks up wrapping token metadata.
type WrappingChecker struct {
	client *http.Client
	base   string
	token  string
}

// NewWrappingChecker creates a WrappingChecker from a Client.
func NewWrappingChecker(c *Client) *WrappingChecker {
	return &WrappingChecker{client: c.http, base: c.addr, token: c.token}
}

// Lookup calls sys/wrapping/lookup for the given wrapping token.
func (w *WrappingChecker) Lookup(ctx context.Context, wrappingToken string) (*WrappingInfo, error) {
	if wrappingToken == "" {
		return nil, fmt.Errorf("wrapping token must not be empty")
	}

	body := fmt.Sprintf(`{"token":%q}`, wrappingToken)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		w.base+"/v1/sys/wrapping/lookup",
		stringReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Vault-Token", w.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := w.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("wrapping lookup returned status %d", resp.Status struct {
		DataWrappingInfo `json:"		return nil, err
	}
	out.Data.TTL = time.Duration(out.Data.TTLSeconds) * time.Second
	return &out.Data, nil
}
