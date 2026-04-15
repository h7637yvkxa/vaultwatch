package vault

import (
	"encoding/json"
	"fmt"
	"time"

	vaultapi "github.com/hashicorp/vault/api"
)

// Client wraps the Vault API client with lease-aware helpers.
type Client struct {
	api *vaultapi.Client
}

// LeaseInfo holds metadata about a single Vault lease.
type LeaseInfo struct {
	LeaseID  string
	TTL      time.Duration
	ExpireAt time.Time
}

// NewClient creates an authenticated Vault client using the provided address
// and token.
func NewClient(address, token string) (*Client, error) {
	cfg := vaultapi.DefaultConfig()
	cfg.Address = address

	raw, err := vaultapi.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("creating vault api client: %w", err)
	}
	raw.SetToken(token)

	return &Client{api: raw}, nil
}

// ListLeases returns all renewable leases visible under the given prefix.
func (c *Client) ListLeases(prefix string) ([]LeaseInfo, error) {
	secret, err := c.api.Logical().List("sys/leases/lookup/" + prefix)
	if err != nil {
		return nil, fmt.Errorf("listing leases under %q: %w", prefix, err)
	}
	if secret == nil || secret.Data == nil {
		return nil, nil
	}

	keys, ok := secret.Data["keys"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected keys format in lease list response")
	}

	var leases []LeaseInfo
	for _, k := range keys {
		key, _ := k.(string)
		if len(key) > 0 && key[len(key)-1] == '/' {
			sub, err := c.ListLeases(prefix + key)
			if err != nil {
				return nil, err
			}
			leases = append(leases, sub...)
			continue
		}
		info, err := c.LookupLease(prefix + key)
		if err != nil {
			return nil, err
		}
		leases = append(leases, *info)
	}
	return leases, nil
}

// LookupLease fetches TTL and expiry details for a specific lease ID.
func (c *Client) LookupLease(leaseID string) (*LeaseInfo, error) {
	secret, err := c.api.Logical().Write("sys/leases/lookup", map[string]interface{}{
		"lease_id": leaseID,
	})
	if err != nil {
		return nil, fmt.Errorf("looking up lease %q: %w", leaseID, err)
	}
	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("empty response for lease %q", leaseID)
	}

	ttlRaw, _ := json.Number(fmt.Sprintf("%v", secret.Data["ttl"])).Int64()
	ttl := time.Duration(ttlRaw) * time.Second

	return &LeaseInfo{
		LeaseID:  leaseID,
		TTL:      ttl,
		ExpireAt: time.Now().Add(ttl),
	}, nil
}
