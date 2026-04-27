package vault_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yourusername/vaultwatch/internal/vault"
)

func newLogicalBackendTestServer(t *testing.T, status int, body any) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_ = json.NewEncoder(w).Encode(body)
	}))
}

func TestListLogicalBackends_Success(t *testing.T) {
	body := map[string]any{
		"secret/": map[string]any{"type": "kv", "description": "key/value store", "local": false, "seal_wrap": false},
		"pki/":    map[string]any{"type": "pki", "description": "PKI backend", "local": false, "seal_wrap": true},
	}
	srv := newLogicalBackendTestServer(t, http.StatusOK, body)
	defer srv.Close()

	checker := vault.NewLogicalBackendChecker(srv.URL, "test-token")
	results, err := checker.ListLogicalBackends(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 backends, got %d", len(results))
	}
}

func TestListLogicalBackends_Empty(t *testing.T) {
	srv := newLogicalBackendTestServer(t, http.StatusOK, map[string]any{})
	defer srv.Close()

	checker := vault.NewLogicalBackendChecker(srv.URL, "test-token")
	results, err := checker.ListLogicalBackends(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Fatalf("expected 0 backends, got %d", len(results))
	}
}

func TestListLogicalBackends_HTTPError(t *testing.T) {
	srv := newLogicalBackendTestServer(t, http.StatusForbidden, map[string]any{})
	defer srv.Close()

	checker := vault.NewLogicalBackendChecker(srv.URL, "bad-token")
	_, err := checker.ListLogicalBackends(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestListLogicalBackends_InvalidURL(t *testing.T) {
	checker := vault.NewLogicalBackendChecker("://invalid", "token")
	_, err := checker.ListLogicalBackends(context.Background())
	if err == nil {
		t.Fatal("expected error for invalid URL")
	}
}
