package vault

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// SecretVersion holds metadata about a KV v2 secret version.
type SecretVersion struct {
	Path      string
	Version   int
	CreatedAt time.Time
	DeletedAt *time.Time
	Destroyed bool
}

// SecretChecker checks KV v2 secret version metadata.
type SecretChecker struct {
	client *http.Client
	address string
	token   string
}

// NewSecretChecker creates a new SecretChecker.
func NewSecretChecker(address, token string) *SecretChecker {
	return &SecretChecker{
		client:  &http.Client{Timeout: 10 * time.Second},
		address: address,
		token:   token,
	}
}

// ListSecretVersions returns version metadata for all versions of a KV v2 secret.
func (c *SecretChecker) ListSecretVersions(ctx context.Context, mount, path string) ([]SecretVersion, error) {
	url := fmt.Sprintf("%s/v1/%s/metadata/%s", c.address, mount, path)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("X-Vault-Token", c.token)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var body struct {
		Data struct {
			Versions map[string]struct {
				CreatedTime  time.Time  `json:"created_time"`
				DeletionTime *time.Time `json:"deletion_time"`
				Destroyed    bool       `json:"destroyed"`
			} `json:"versions"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}

	var versions []SecretVersion
	for k, v := range body.Data.Versions {
		var ver int
		fmt.Sscanf(k, "%d", &ver)
		versions = append(versions, SecretVersion{
			Path:      path,
			Version:   ver,
			CreatedAt: v.CreatedTime,
			DeletedAt: v.DeletionTime,
			Destroyed: v.Destroyed,
		})
	}
	return versions, nil
}
