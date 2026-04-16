package vault

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newMountTestServer(t *testing.T, status int, body interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		if body != nil {
			_ = json.NewEncoder(w).Encode(body)
		}
	}))
}

func TestListMounts_Success(t *testing.T) {
	body := map[string]interface{}{
		"secret/": map[string]string{"type": "kv", "description": "KV store", "accessor": "abc123"},
		"pki/":    map[string]string{"type": "pki", "description": "PKI engine", "accessor": "def456"},
	}
	srv := newMountTestServer(t, http.StatusOK, body)
	defer srv.Close()

	checker := NewMountChecker(srv.URL, "test-token", nil)
	mounts, err := checker.ListMounts(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(mounts) != 2 {
		t.Fatalf("expected 2 mounts, got %d", len(mounts))
	}
	types := map[string]bool{}
	for _, m := range mounts {
		types[m.Type] = true
	}
	if !types["kv"] || !types["pki"] {
		t.Errorf("expected kv and pki types, got %v", types)
	}
}

func TestListMounts_Empty(t *testing.T) {
	srv := newMountTestServer(t, http.StatusOK, map[string]interface{}{})
	defer srv.Close()

	checker := NewMountChecker(srv.URL, "test-token", nil)
	mounts, err := checker.ListMounts(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(mounts) != 0 {
		t.Errorf("expected 0 mounts, got %d", len(mounts))
	}
}

func TestListMounts_HTTPError(t *testing.T) {
	srv := newMountTestServer(t, http.StatusForbidden, nil)
	defer srv.Close()

	checker := NewMountChecker(srv.URL, "bad-token", nil)
	_, err := checker.ListMounts(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
