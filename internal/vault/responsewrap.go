package vault

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// ResponseWrapInfo holds metadata about a response-wrapped secret.
type ResponseWrapInfo struct {
	Token          string    `json:"token"`
	Accessor       string    `json:"accessor"`
	TTL            int       `json:"ttl"`
	CreationTime   time.Time `json:"creation_time"`
	CreationPath   string    `json:"creation_path"`
	WrappedAccessor string   `json:"wrapped_accessor"`
}

// ResponseWrapChecker inspects a wrapping token via the sys/wrapping/lookup endpoint.
type ResponseWrapChecker struct {
	address string
	token   string
	client  *http.Client
}

// NewResponseWrapChecker creates a ResponseWrapChecker using the provided Vault client.
func NewResponseWrapChecker(c *Client) *ResponseWrapChecker {
	return &ResponseWrapChecker{
		address: c.Address,
		token:   c.Token,
		client:  c.HTTP,
	}
}

// Lookup returns the wrapping info for the given wrapping token.
func (r *ResponseWrapChecker) Lookup(ctx context.Context, wrappingToken string) (*ResponseWrapInfo, error) {
	if wrappingToken == "" {
		return nil, fmt.Errorf("wrapping token must not be empty")
	}

	url := fmt.Sprintf("%s/v1/sys/wrapping/lookup", r.address)
	body := fmt.Sprintf(`{"token":%q}`, wrappingToken)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url,
		io.NopCloser(newStringReader(body)))
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("X-Vault-Token", r.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var result struct {
		Data ResponseWrapInfo `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return &result.Data, nil
}

// newStringReader is a helper to avoid importing strings in this file.
func newStringReader(s string) io.Reader {
	return stringReader(s)
}

type stringReader string

func (s stringReader) Read(p []byte) (int, error) {
	copy(p, s)
	n := len(s)
	if n > len(p) {
		n = len(p)
	}
	return n, io.EOF
}
