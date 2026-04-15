package vault_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	vaultapi "github.com/hashicorp/vault/api"

	"github.com/yourusername/vaultwatch/internal/vault"
)

func newListTestServer(t *testing.T, responseBody map[string]interface{}, statusCode int) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"data": responseBody,
		})
	}))
}

func newVaultClient(t *testing.T, addr string) *vaultapi.Client {
	t.Helper()
	cfg := vaultapi.DefaultConfig()
	cfg.Address = addr
	client, err := vaultapi.NewClient(cfg)
	if err != nil {
		t.Fatalf("failed to create vault client: %v", err)
	}
	client.SetToken("test-token")
	return client
}

func TestListLeases_Empty(t *testing.T) {
	srv := newListTestServer(t, map[string]interface{}{}, http.StatusOK)
	defer srv.Close()

	lister := vault.NewLister(newVaultClient(t, srv.URL))
	ids, err := lister.ListLeases(context.Background(), "aws/creds")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ids) != 0 {
		t.Errorf("expected 0 leases, got %d", len(ids))
	}
}

func TestListLeases_ReturnsLeafIDs(t *testing.T) {
	body := map[string]interface{}{
		"keys": []interface{}{
			"lease-abc",
			"lease-def",
			"subdir/", // should be skipped
		},
	}
	srv := newListTestServer(t, body, http.StatusOK)
	defer srv.Close()

	lister := vault.NewLister(newVaultClient(t, srv.URL))
	ids, err := lister.ListLeases(context.Background(), "aws/creds")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ids) != 2 {
		t.Fatalf("expected 2 lease IDs, got %d", len(ids))
	}
	for _, id := range ids {
		if id == "subdir/" {
			t.Error("sub-directory entry should have been skipped")
		}
	}
}

func TestListLeases_StripLeadingSlash(t *testing.T) {
	body := map[string]interface{}{
		"keys": []interface{}{"lease-xyz"},
	}
	srv := newListTestServer(t, body, http.StatusOK)
	defer srv.Close()

	lister := vault.NewLister(newVaultClient(t, srv.URL))
	// Prefix with leading slash should not cause double-slash in result.
	ids, err := lister.ListLeases(context.Background(), "/aws/creds")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ids) != 1 {
		t.Fatalf("expected 1 lease ID, got %d", len(ids))
	}
}
