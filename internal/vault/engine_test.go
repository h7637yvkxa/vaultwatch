package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newEngineTestServer(t *testing.T, status int, body any) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/sys/mounts" {
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(status)
		_ = json.NewEncoder(w).Encode(body)
	}))
}

func TestListEngines_Success(t *testing.T) {
	body := map[string]any{
		"secret/": map[string]any{"type": "kv", "description": "KV store", "local": false, "seal_wrap": false},
		"pki/":    map[string]any{"type": "pki", "description": "PKI", "local": true, "seal_wrap": true},
	}
	srv := newEngineTestServer(t, http.StatusOK, body)
	defer srv.Close()

	c := newVaultClient(t, srv.URL)
	checker := NewEngineChecker(c)
	engines, err := checker.ListEngines()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(engines) != 2 {
		t.Fatalf("expected 2 engines, got %d", len(engines))
	}
}

func TestListEngines_Empty(t *testing.T) {
	srv := newEngineTestServer(t, http.StatusOK, map[string]any{})
	defer srv.Close()

	c := newVaultClient(t, srv.URL)
	checker := NewEngineChecker(c)
	engines, err := checker.ListEngines()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(engines) != 0 {
		t.Fatalf("expected 0 engines, got %d", len(engines))
	}
}

func TestListEngines_HTTPError(t *testing.T) {
	srv := newEngineTestServer(t, http.StatusForbidden, nil)
	defer srv.Close()

	c := newVaultClient(t, srv.URL)
	checker := NewEngineChecker(c)
	_, err := checker.ListEngines()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestListEngines_InvalidURL(t *testing.T) {
	c := newVaultClient(t, "http://127.0.0.1:0")
	checker := NewEngineChecker(c)
	_, err := checker.ListEngines()
	if err == nil {
		t.Fatal("expected error for invalid URL")
	}
}
