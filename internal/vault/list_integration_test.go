//go:build integration
// +build integration

package vault_test

import (
	"context"
	"os"
	"testing"

	vaultapi "github.com/hashicorp/vault/api"

	"github.com/yourusername/vaultwatch/internal/vault"
)

// TestListLeases_Integration requires a running Vault instance.
// Set VAULT_ADDR and VAULT_TOKEN env vars before running.
//
// Run with: go test -tags integration ./internal/vault/...
func TestListLeases_Integration(t *testing.T) {
	addr := os.Getenv("VAULT_ADDR")
	token := os.Getenv("VAULT_TOKEN")
	prefix := os.Getenv("VAULT_LEASE_PREFIX")

	if addr == "" || token == "" {
		t.Skip("VAULT_ADDR and VAULT_TOKEN must be set for integration tests")
	}
	if prefix == "" {
		prefix = "aws/creds"
	}

	cfg := vaultapi.DefaultConfig()
	cfg.Address = addr
	client, err := vaultapi.NewClient(cfg)
	if err != nil {
		t.Fatalf("failed to create vault client: %v", err)
	}
	client.SetToken(token)

	lister := vault.NewLister(client)
	ids, err := lister.ListLeases(context.Background(), prefix)
	if err != nil {
		t.Fatalf("ListLeases error: %v", err)
	}

	t.Logf("found %d leases under %q", len(ids), prefix)
	for _, id := range ids {
		t.Logf("  lease: %s", id)
	}
}
