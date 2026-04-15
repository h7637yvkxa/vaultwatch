package vault

import (
	"context"
	"fmt"
	"strings"

	vaultapi "github.com/hashicorp/vault/api"
)

// LeaseEntry represents a single lease returned from Vault's sys/leases/lookup.
type LeaseEntry struct {
	LeaseID    string
	TTL        int
	Renewable  bool
}

// Lister lists leases under a given prefix from Vault.
type Lister struct {
	client *vaultapi.Client
}

// NewLister creates a new Lister using the provided Vault client.
func NewLister(client *vaultapi.Client) *Lister {
	return &Lister{client: client}
}

// ListLeases returns all lease IDs found under the given prefix path.
func (l *Lister) ListLeases(ctx context.Context, prefix string) ([]string, error) {
	prefix = strings.TrimPrefix(prefix, "/")
	path := fmt.Sprintf("sys/leases/lookup/%s", prefix)

	secret, err := l.client.Logical().ListWithContext(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("listing leases at %q: %w", path, err)
	}
	if secret == nil || secret.Data == nil {
		return []string{}, nil
	}

	keys, ok := secret.Data["keys"]
	if !ok {
		return []string{}, nil
	}

	raw, ok := keys.([]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected type for lease keys at %q", path)
	}

	leaseIDs := make([]string, 0, len(raw))
	for _, k := range raw {
		s, ok := k.(string)
		if !ok {
			continue
		}
		// Skip sub-directories (trailing slash), only collect leaf IDs.
		if !strings.HasSuffix(s, "/") {
			leaseIDs = append(leaseIDs, prefix+"/"+s)
		}
	}
	return leaseIDs, nil
}
