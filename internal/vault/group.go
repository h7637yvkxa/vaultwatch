package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// GroupEntry represents a Vault identity group.
type GroupEntry struct {
	ID       string            `json:"id"`
	Name     string            `json:"name"`
	Type     string            `json:"type"`
	Policies []string          `json:"policies"`
	Metadata map[string]string `json:"metadata"`
}

// GroupChecker lists identity groups from Vault.
type GroupChecker struct {
	addr   string
	token  string
	client *http.Client
}

// NewGroupChecker returns a GroupChecker using the provided Vault client.
func NewGroupChecker(c *Client) *GroupChecker {
	return &GroupChecker{
		addr:   c.addr,
		token:  c.token,
		client: c.http,
	}
}

// ListGroups returns all identity groups from Vault.
func (g *GroupChecker) ListGroups() ([]GroupEntry, error) {
	url := fmt.Sprintf("%s/v1/identity/group/id?list=true", g.addr)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("group: build request: %w", err)
	}
	req.Header.Set("X-Vault-Token", g.token)

	resp, err := g.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("group: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return []GroupEntry{}, nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("group: unexpected status %d", resp.StatusCode)
	}

	var body struct {
		Data struct {
			KeyInfo map[string]GroupEntry `json:"key_info"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, fmt.Errorf("group: decode: %w", err)
	}

	groups := make([]GroupEntry, 0, len(body.Data.KeyInfo))
	for id, entry := range body.Data.KeyInfo {
		entry.ID = id
		groups = append(groups, entry)
	}
	return groups, nil
}
