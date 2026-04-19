package vault

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// CertStoreEntry represents a certificate in the Vault cert auth method.
type CertStoreEntry struct {
	Name       string
	DisplayName string
	TTL        time.Duration
	MaxTTL     time.Duration
}

// CertStoreChecker lists certificates from the Vault cert auth backend.
type CertStoreChecker struct {
	address string
	token   string
	client  *http.Client
	mount   string
}

// NewCertStoreChecker creates a new CertStoreChecker.
func NewCertStoreChecker(address, token, mount string) *CertStoreChecker {
	return &CertStoreChecker{
		address: address,
		token:   token,
		client:  &http.Client{Timeout: 10 * time.Second},
		mount:   mount,
	}
}

// ListCerts returns all certificate entries from the cert auth mount.
func (c *CertStoreChecker) ListCerts() ([]CertStoreEntry, error) {
	url := fmt.Sprintf("%s/v1/auth/%s/certs?list=true", c.address, c.mount)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Vault-Token", c.token)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("certstore list failed: %s", string(body))
	}

	var result struct {
		Data struct {
			Keys []string `json:"keys"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	var entries []CertStoreEntry
	for _, key := range result.Data.Keys {
		entries = append(entries, CertStoreEntry{Name: key})
	}
	return entries, nil
}
