package vault

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func newResponseWrapTestServer(t *testing.T, status int, payload interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		if payload != nil {
			_ = json.NewEncoder(w).Encode(payload)
		}
	}))
}

func TestResponseWrapLookup_Success(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	payload := map[string]interface{}{
		"data": map[string]interface{}{
			"token":           "s.abc123",
			"accessor":        "acc-xyz",
			"ttl":             300,
			"creation_time":   now,
			"creation_path":   "secret/data/foo",
			"wrapped_accessor": "wacc-xyz",
		},
	}
	srv := newResponseWrapTestServer(t, http.StatusOK, payload)
	defer srv.Close()

	c := newVaultClient(t, srv.URL)
	checker := NewResponseWrapChecker(c)

	info, err := checker.Lookup(context.Background(), "s.abc123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info.Token != "s.abc123" {
		t.Errorf("expected token s.abc123, got %s", info.Token)
	}
	if info.TTL != 300 {
		t.Errorf("expected TTL 300, got %d", info.TTL)
	}
	if info.CreationPath != "secret/data/foo" {
		t.Errorf("unexpected creation_path: %s", info.CreationPath)
	}
}

func TestResponseWrapLookup_EmptyToken(t *testing.T) {
	srv := newResponseWrapTestServer(t, http.StatusOK, nil)
	defer srv.Close()

	c := newVaultClient(t, srv.URL)
	checker := NewResponseWrapChecker(c)

	_, err := checker.Lookup(context.Background(), "")
	if err == nil {
		t.Fatal("expected error for empty token")
	}
}

func TestResponseWrapLookup_HTTPError(t *testing.T) {
	srv := newResponseWrapTestServer(t, http.StatusForbidden, nil)
	defer srv.Close()

	c := newVaultClient(t, srv.URL)
	checker := NewResponseWrapChecker(c)

	_, err := checker.Lookup(context.Background(), "s.badtoken")
	if err == nil {
		t.Fatal("expected error on 403")
	}
}

func TestResponseWrapLookup_InvalidURL(t *testing.T) {
	c := newVaultClient(t, "http://127.0.0.1:0")
	checker := NewResponseWrapChecker(c)

	_, err := checker.Lookup(context.Background(), "s.sometoken")
	if err == nil {
		t.Fatal("expected connection error")
	}
}
