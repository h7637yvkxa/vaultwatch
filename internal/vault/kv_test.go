package vault_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/yourusername/vaultwatch/internal/vault"
)

func newKVTestServer(t *testing.T, status int, payload interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		if payload != nil {
			_ = json.NewEncoder(w).Encode(payload)
		}
	}))
}

func TestReadSecret_Success(t *testing.T) {
	created := time.Now().UTC().Truncate(time.Second)
	payload := map[string]interface{}{
		"data": map[string]interface{}{
			"data": map[string]interface{}{"username": "admin", "password": "s3cr3t"},
			"metadata": map[string]interface{}{
				"version":      float64(3),
				"created_time": created.Format(time.RFC3339),
			},
		},
	}
	srv := newKVTestServer(t, http.StatusOK, payload)
	defer srv.Close()

	checker := vault.NewKVChecker(srv.URL, "test-token", nil)
	secret, err := checker.ReadSecret("secret", "myapp/creds")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if secret.Version != 3 {
		t.Errorf("expected version 3, got %d", secret.Version)
	}
	if secret.Data["username"] != "admin" {
		t.Errorf("expected username admin, got %v", secret.Data["username"])
	}
}

func TestReadSecret_HTTPError(t *testing.T) {
	srv := newKVTestServer(t, http.StatusForbidden, nil)
	defer srv.Close()

	checker := vault.NewKVChecker(srv.URL, "bad-token", nil)
	_, err := checker.ReadSecret("secret", "myapp/creds")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestReadSecret_InvalidURL(t *testing.T) {
	checker := vault.NewKVChecker("http://127.0.0.1:0", "token", nil)
	_, err := checker.ReadSecret("secret", "myapp/creds")
	if err == nil {
		t.Fatal("expected error for unreachable server")
	}
}
