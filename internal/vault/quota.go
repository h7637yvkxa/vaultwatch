package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// QuotaEntry represents a rate-limit or lease-count quota in Vault.
type QuotaEntry struct {
	Name       string  `json:"name"`
	Type       string  `json:"type"`
	Path       string  `json:"path"`
	MaxLeases  int     `json:"max_leases"`
	Rate       float64 `json:"rate"`
	Interval   int     `json:"interval"`
	BlockInterval int  `json:"block_interval"`
}

type quotaListResponse struct {
	Data struct {
		Keys []string `json:"keys"`
	} `json:"data"`
}

type quotaGetResponse struct {
	Data QuotaEntry `json:"data"`
}

// QuotaChecker fetches quota rules from Vault.
type QuotaChecker struct {
	address string
	token   string
	client  *http.Client
}

// NewQuotaChecker creates a QuotaChecker using the provided Vault client.
func NewQuotaChecker(address, token string, client *http.Client) *QuotaChecker {
	if client == nil {
		client = http.DefaultClient
	}
	return &QuotaChecker{address: address, token: token, client: client}
}

// ListQuotas returns all configured quota rules.
func (q *QuotaChecker) ListQuotas() ([]QuotaEntry, error) {
	req, err := http.NewRequest(http.MethodGet, q.address+"/v1/sys/quotas/rate-limit?list=true", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Vault-Token", q.token)

	resp, err := q.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return []QuotaEntry{}, nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("quota list: unexpected status %d", resp.StatusCode)
	}

	var list quotaListResponse
	if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
		return nil, err
	}

	var entries []QuotaEntry
	for _, name := range list.Data.Keys {
		entry, err := q.getQuota(name)
		if err != nil {
			continue
		}
		entries = append(entries, *entry)
	}
	return entries, nil
}

func (q *QuotaChecker) getQuota(name string) (*QuotaEntry, error) {
	req, err := http.NewRequest(http.MethodGet, q.address+"/v1/sys/quotas/rate-limit/"+name, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Vault-Token", q.token)

	resp, err := q.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("quota get %s: status %d", name, resp.StatusCode)
	}

	var gr quotaGetResponse
	if err := json.NewDecoder(resp.Body).Decode(&gr); err != nil {
		return nil, err
	}
	return &gr.Data, nil
}
