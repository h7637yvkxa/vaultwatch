package vault

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// EntityEntry represents a Vault identity entity.
type EntityEntry struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Policies []string `json:"policies"`
	Disabled bool     `json:"disabled"`
}

// EntityChecker lists identity entities from Vault.
type EntityChecker struct {
	client  *http.Client
	baseURL string
	token   string
}

// NewEntityChecker constructs an EntityChecker from a Client.
func NewEntityChecker(c *Client) *EntityChecker {
	return &EntityChecker{
		client:  c.http,
		baseURL: c.baseURL,
		token:   c.token,
	}
}

// ListEntities returns all identity entities.
func (e *EntityChecker) ListEntities(ctx context.Context) ([]EntityEntry, error) {
	url := fmt.Sprintf("%s/v1/identity/entity/id?list=true", e.baseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Vault-Token", e.token)

	resp, err := e.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("entity list: unexpected status %d", resp.StatusCode)
	}

	var body struct {
		Data struct {
			KeyInfo map[string]EntityEntry `json:"key_info"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, err
	}

	entries := make([]EntityEntry, 0, len(body.Data.KeyInfo))
	for id, entry := range body.Data.KeyInfo {
		entry.ID = id
		entries = append(entries, entry)
	}
	return entries, nil
}
