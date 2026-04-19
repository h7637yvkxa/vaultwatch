package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// TransitKey represents a Vault transit encryption key.
type TransitKey struct {
	Name            string `json:"name"`
	Type            string `json:"type"`
	DeletionAllowed bool   `json:"deletion_allowed"`
	Exportable      bool   `json:"exportable"`
	MinDecryptVersion int  `json:"min_decryption_version"`
	LatestVersion   int    `json:"latest_version"`
}

// TransitChecker lists transit keys from Vault.
type TransitChecker struct {
	addr   string
	token  string
	client *http.Client
}

// NewTransitChecker creates a new TransitChecker.
func NewTransitChecker(addr, token string, client *http.Client) *TransitChecker {
	if client == nil {
		client = http.DefaultClient
	}
	return &TransitChecker{addr: addr, token: token, client: client}
}

// ListTransitKeys returns all transit keys found under the given mount path.
func (tc *TransitChecker) ListTransitKeys(mount string) ([]TransitKey, error) {
	url := fmt.Sprintf("%s/v1/%s/keys?list=true", tc.addr, mount)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("transit: build request: %w", err)
	}
	req.Header.Set("X-Vault-Token", tc.token)

	resp, err := tc.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("transit: request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("transit: unexpected status %d", resp.StatusCode)
	}

	var body struct {
		Data struct {
			Keys []string `json:"keys"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, fmt.Errorf("transit: decode: %w", err)
	}

	keys := make([]TransitKey, 0, len(body.Data.Keys))
	for _, name := range body.Data.Keys {
		keys = append(keys, TransitKey{Name: name})
	}
	return keys, nil
}
