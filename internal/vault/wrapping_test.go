package vault

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newWrappingTestServer(t *testing.T, status int, payload any) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_ = json.NewEncoder(w).Encode(payload)
	}))
}

func TestWrappingLookup_Success(t *testing.T) {
	payload := map[string]any{
		"data": map[string]any{
			"token":         "s.abc123",
			"accessor":      "acc456",
			"ttl":           300,
			"creation_time": "2024-01-01T00:00:00Z",
			"creation_path": "auth/token/create",
		},
	}
	srv := newWrappingTestServer(t, http.StatusOK, payload)
	defer srv.Close()

	c := newVaultClient(t, srv.URL)
	wc := NewWrappingChecker(c)

	info, err := wc.Lookup(context.Background(), "s.abc123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info.Accessor != "acc456" {
		t.Errorf("expected accessor acc456, got %s", info.Accessor)
	}
	if info.TTL.Seconds() != 300 {
		t.Errorf("expected TTL 300s, got %v", info.TTL)
	}
	if info.CreationPath != "auth/token/create" {
		t.Errorf("unexpected creation_path: %s", info.CreationPath)
	}
}

func TestWrappingLookup_EmptyToken(t *testing.T) {
	c := newVaultClient(t, "http://127.0.0.1")
	wc := NewWrappingChecker(c)

	_, err := wc.Lookup(context.Background(), "")
	if err == nil {
		t.Fatal("expected error for empty token")
	}
}

func TestWrappingLookup_HTTPError(t *testing.T) {
	srv := newWrappingTestServer(t, http.StatusForbidden, map[string]string{"errors": "permission denied"})
	defer srv.Close()

	c := newVaultClient(t, srv.URL)
	wc := NewWrappingChecker(c)

	_, err := wc.Lookup(context.Background(), "s.bad")
	if err == nil {
		t.Fatal("expected error for non-200 status")
	}
}

func TestWrappingLookup_ZeroTTL(t *testing.T) {
	payload := map[string]any{
		"data": map[string]any{
			"token":         "s.expired",
			"accessor":      "acc789",
			"ttl":           0,
			"creation_time": "2024-01-01T00:00:00Z",
			"creation_path": "auth/token/create",
		},
	}
	srv := newWrappingTestServer(t, http.StatusOK, payload)
	defer srv.Close()

	c := newVaultClient(t, srv.URL)
	wc := NewWrappingChecker(c)

	info, err := wc.Lookup(context.Background(), "s.expired")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info.TTL != 0 {
		t.Errorf("expected zero TTL, got %v", info.TTL)
	}
}
