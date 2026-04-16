package vault_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yourusername/vaultwatch/internal/vault"
)

func newTokenTestServer(t *testing.T, data map[string]interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/auth/token/lookup-self" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"data": data})
	}))
}

func TestLookupSelf_Success(t *testing.T) {
	srv := newTokenTestServer(t, map[string]interface{}{
		"id":           "test-token-id",
		"display_name": "mytoken",
		"renewable":    true,
		"ttl":          float64(3600),
		"policies":     []interface{}{"default", "admin"},
	})
	defer srv.Close()

	client := newVaultClient(t, srv.URL)
	checker := vault.NewTokenChecker(client)

	info, err := checker.LookupSelf(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info.ID != "test-token-id" {
		t.Errorf("expected id test-token-id, got %s", info.ID)
	}
	if info.DisplayName != "mytoken" {
		t.Errorf("expected display_name mytoken, got %s", info.DisplayName)
	}
	if !info.Renewable {
		t.Error("expected renewable true")
	}
	if info.TTL.Seconds() != 3600 {
		t.Errorf("expected TTL 3600s, got %v", info.TTL)
	}
	if len(info.Policies) != 2 {
		t.Errorf("expected 2 policies, got %d", len(info.Policies))
	}
}

func TestLookupSelf_HTTPError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer srv.Close()

	client := newVaultClient(t, srv.URL)
	checker := vault.NewTokenChecker(client)

	_, err := checker.LookupSelf(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
