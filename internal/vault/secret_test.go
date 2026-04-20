package vault_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/yourusername/vaultwatch/internal/vault"
)

func newSecretTestServer(t *testing.T, status int, payload any) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		if payload != nil {
			json.NewEncoder(w).Encode(payload)
		}
	}))
}

func TestListSecretVersions_Success(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	payload := map[string]any{
		"data": map[string]any{
			"versions": map[string]any{
				"1": map[string]any{
					"created_time": now.Format(time.RFC3339),
					"destroyed":    false,
				},
				"2": map[string]any{
					"created_time": now.Add(-time.Hour).Format(time.RFC3339),
					"destroyed":    true,
				},
			},
		},
	}
	srv := newSecretTestServer(t, http.StatusOK, payload)
	defer srv.Close()

	c := vault.NewSecretChecker(srv.URL, "test-token")
	versions, err := c.ListSecretVersions(context.Background(), "secret", "myapp/db")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(versions) != 2 {
		t.Fatalf("expected 2 versions, got %d", len(versions))
	}
}

func TestListSecretVersions_NotFound(t *testing.T) {
	srv := newSecretTestServer(t, http.StatusNotFound, nil)
	defer srv.Close()

	c := vault.NewSecretChecker(srv.URL, "test-token")
	versions, err := c.ListSecretVersions(context.Background(), "secret", "missing")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if versions != nil {
		t.Fatalf("expected nil, got %v", versions)
	}
}

func TestListSecretVersions_HTTPError(t *testing.T) {
	srv := newSecretTestServer(t, http.StatusForbidden, nil)
	defer srv.Close()

	c := vault.NewSecretChecker(srv.URL, "bad-token")
	_, err := c.ListSecretVersions(context.Background(), "secret", "myapp/db")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestListSecretVersions_InvalidURL(t *testing.T) {
	c := vault.NewSecretChecker("://invalid", "token")
	_, err := c.ListSecretVersions(context.Background(), "secret", "path")
	if err == nil {
		t.Fatal("expected error for invalid URL")
	}
}
