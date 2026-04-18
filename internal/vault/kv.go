package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// KVSecret represents a secret read from a KV v2 mount.
type KVSecret struct {
	Path      string
	Version   int
	CreatedAt time.Time
	Data      map[string]interface{}
}

// KVChecker reads KV v2 secrets from Vault.
type KVChecker struct {
	client *http.Client
	baseURL string
	token   string
}

// NewKVChecker returns a KVChecker using the provided Vault client.
func NewKVChecker(baseURL, token string, client *http.Client) *KVChecker {
	if client == nil {
		client = http.DefaultClient
	}
	return &KVChecker{client: client, baseURL: strings.TrimRight(baseURL, "/"), token: token}
}

// ReadSecret reads a KV v2 secret at the given mount and path.
func (k *KVChecker) ReadSecret(mount, path string) (*KVSecret, error) {
	url := fmt.Sprintf("%s/v1/%s/data/%s", k.baseURL, mount, strings.TrimLeft(path, "/"))
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("kv: build request: %w", err)
	}
	req.Header.Set("X-Vault-Token", k.token)

	resp, err := k.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("kv: request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("kv: unexpected status %d for %s/%s", resp.StatusCode, mount, path)
	}

	var body struct {
		Data struct {
			Data     map[string]interface{} `json:"data"`
			Metadata struct {
				Version   int       `json:"version"`
				CreatedAt time.Time `json:"created_time"`
			} `json:"metadata"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, fmt.Errorf("kv: decode: %w", err)
	}

	return &KVSecret{
		Path:      path,
		Version:   body.Data.Metadata.Version,
		CreatedAt: body.Data.Metadata.CreatedAt,
		Data:      body.Data.Data,
	}, nil
}
