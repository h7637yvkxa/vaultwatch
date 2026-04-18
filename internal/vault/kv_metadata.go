package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// KVMetadata holds metadata for a KV v2 secret.
type KVMetadata struct {
	Path            string
	CurrentVersion  int
	OldestVersion   int
	CreatedTime     time.Time
	UpdatedTime     time.Time
	MaxVersions     int
	DeleteVersionAfter string
}

// KVMetadataChecker reads KV v2 secret metadata from Vault.
type KVMetadataChecker struct {
	addr   string
	token  string
	client *http.Client
}

// NewKVMetadataChecker creates a new KVMetadataChecker.
func NewKVMetadataChecker(addr, token string) *KVMetadataChecker {
	return &KVMetadataChecker{addr: addr, token: token, client: &http.Client{Timeout: 10 * time.Second}}
}

// ReadMetadata fetches metadata for the given KV v2 mount and secret path.
func (c *KVMetadataChecker) ReadMetadata(mount, path string) (*KVMetadata, error) {
	url := fmt.Sprintf("%s/v1/%s/metadata/%s", c.addr, mount, path)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("X-Vault-Token", c.token)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var body struct {
		Data struct {
			CurrentVersion     int       `json:"current_version"`
			OldestVersion      int       `json:"oldest_version"`
			CreatedTime        time.Time `json:"created_time"`
			UpdatedTime        time.Time `json:"updated_time"`
			MaxVersions        int       `json:"max_versions"`
			DeleteVersionAfter string    `json:"delete_version_after"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return &KVMetadata{
		Path:               fmt.Sprintf("%s/%s", mount, path),
		CurrentVersion:     body.Data.CurrentVersion,
		OldestVersion:      body.Data.OldestVersion,
		CreatedTime:        body.Data.CreatedTime,
		UpdatedTime:        body.Data.UpdatedTime,
		MaxVersions:        body.Data.MaxVersions,
		DeleteVersionAfter: body.Data.DeleteVersionAfter,
	}, nil
}
