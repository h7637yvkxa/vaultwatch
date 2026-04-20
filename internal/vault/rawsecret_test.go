package vault_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yourusername/vaultwatch/internal/vault"
)

func newRawSecretTestServer(t *testing.T, status int, payload interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		if payload != nil {
			_ = json.NewEncoder(w).Encode(payload)
		}
	}))
}

func TestReadRawSecret_Success(t *testing.T) {
	payload := map[string]interface{}{
		"data": map[string]interface{}{
			"username": "admin",
			"password": "s3cr3t",
		},
	}
	srv := newRawSecretTestServer(t, http.StatusOK, payload)
	defer srv.Close()

	c := newVaultClient(t, srv.URL)
	checker := vault.NewRawSecretChecker(c)

	entry, err := checker.ReadRawSecret(context.Background(), "secret/myapp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry.Path != "secret/myapp" {
		t.Errorf("expected path %q, got %q", "secret/myapp", entry.Path)
	}
	if entry.Data["username"] != "admin" {
		t.Errorf("expected username=admin, got %v", entry.Data["username"])
	}
}

func TestReadRawSecret_NotFound(t *testing.T) {
	srv := newRawSecretTestServer(t, http.StatusNotFound, nil)
	defer srv.Close()

	c := newVaultClient(t, srv.URL)
	checker := vault.NewRawSecretChecker(c)

	_, err := checker.ReadRawSecret(context.Background(), "secret/missing")
	if err == nil {
		t.Fatal("expected error for 404, got nil")
	}
}

func TestReadRawSecret_HTTPError(t *testing.T) {
	srv := newRawSecretTestServer(t, http.StatusInternalServerError, nil)
	defer srv.Close()

	c := newVaultClient(t, srv.URL)
	checker := vault.NewRawSecretChecker(c)

	_, err := checker.ReadRawSecret(context.Background(), "secret/broken")
	if err == nil {
		t.Fatal("expected error for 500, got nil")
	}
}

func TestReadRawSecret_EmptyPath(t *testing.T) {
	srv := newRawSecretTestServer(t, http.StatusOK, nil)
	defer srv.Close()

	c := newVaultClient(t, srv.URL)
	checker := vault.NewRawSecretChecker(c)

	_, err := checker.ReadRawSecret(context.Background(), "")
	if err == nil {
		t.Fatal("expected error for empty path, got nil")
	}
}

func TestReadRawSecret_InvalidURL(t *testing.T) {
	c := newVaultClient(t, "http://127.0.0.1:0")
	checker := vault.NewRawSecretChecker(c)

	_, err := checker.ReadRawSecret(context.Background(), "secret/test")
	if err == nil {
		t.Fatal("expected connection error, got nil")
	}
}
