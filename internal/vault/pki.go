package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// PKICert holds metadata about a PKI certificate.
type PKICert struct {
	SerialNumber string
	Expiry       time.Time
	IssuingCA    string
}

// PKIChecker checks PKI certificate expiry via Vault's PKI secrets engine.
type PKIChecker struct {
	client *http.Client
	baseURL string
	token   string
}

// NewPKIChecker creates a new PKIChecker.
func NewPKIChecker(baseURL, token string) *PKIChecker {
	return &PKIChecker{
		client:  &http.Client{Timeout: 10 * time.Second},
		baseURL: baseURL,
		token:   token,
	}
}

// ListCerts lists PKI certificate serial numbers under the given mount path.
func (p *PKIChecker) ListCerts(mount string) ([]PKICert, error) {
	url := fmt.Sprintf("%s/v1/%s/certs?list=true", p.baseURL, mount)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("pki: build request: %w", err)
	}
	req.Header.Set("X-Vault-Token", p.token)

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("pki: request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("pki: unexpected status %d", resp.StatusCode)
	}

	var body struct {
		Data struct {
			Keys []string `json:"keys"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, fmt.Errorf("pki: decode: %w", err)
	}

	certs := make([]PKICert, 0, len(body.Data.Keys))
	for _, serial := range body.Data.Keys {
		certs = append(certs, PKICert{SerialNumber: serial})
	}
	return certs, nil
}
