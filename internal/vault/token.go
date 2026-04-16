package vault

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	vaultapi "github.com/hashicorp/vault/api"
)

// TokenInfo holds metadata about a Vault token.
type TokenInfo struct {
	ID          string
	DisplayName string
	Policies    []string
	TTL         time.Duration
	Renewable   bool
}

// TokenChecker inspects the current Vault token.
type TokenChecker struct {
	client *vaultapi.Client
}

// NewTokenChecker returns a TokenChecker using the provided client.
func NewTokenChecker(client *vaultapi.Client) *TokenChecker {
	return &TokenChecker{client: client}
}

// LookupSelf returns metadata about the token currently configured on the client.
func (tc *TokenChecker) LookupSelf(ctx context.Context) (*TokenInfo, error) {
	secret, err := tc.client.Auth().Token().LookupSelfWithContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("token lookup-self: %w", err)
	}
	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("token lookup-self: empty response")
	}

	info := &TokenInfo{}

	if v, ok := secret.Data["id"].(string); ok {
		info.ID = v
	}
	if v, ok := secret.Data["display_name"].(string); ok {
		info.DisplayName = v
	}
	if v, ok := secret.Data["renewable"].(bool); ok {
		info.Renewable = v
	}
	if policies, ok := secret.Data["policies"].([]interface{}); ok {
		for _, p := range policies {
			if s, ok := p.(string); ok {
				info.Policies = append(info.Policies, s)
			}
		}
	}
	if raw, ok := secret.Data["ttl"]; ok {
		var n json.Number
		switch v := raw.(type) {
		case json.Number:
			n = v
		case float64:
			n = json.Number(fmt.Sprintf("%.0f", v))
		}
		if i, err := n.Int64(); err == nil {
			info.TTL = time.Duration(i) * time.Second
		}
	}

	return info, nil
}
