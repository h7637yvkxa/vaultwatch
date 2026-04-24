package vault

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newLeaseCountTestServer(t *testing.T, status int, payload any) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Vault-Token") == "" {
			http.Error(w, "missing token", http.StatusForbidden)
			return
		}
		w.WriteHeader(status)
		if payload != nil {
			_ = json.NewEncoder(w).Encode(payload)
		}
	}))
}

func TestLeaseCount_Success(t *testing.T) {
	payload := map[string]any{
		"data": map[string]any{
			"lease_count": 42,
			"count_per_mount": map[string]any{
				"secret/": 30,
				"aws/": 12,
			},
		},
	}
	srv := newLeaseCountTestServer(t, http.StatusOK, payload)
	defer srv.Close()

	checker := NewLeaseCountChecker(srv.URL, "test-token")
	result, err := checker.Count(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Total != 42 {
		t.Errorf("expected total 42, got %d", result.Total)
	}
	if result.ByMount["secret/"] != 30 {
		t.Errorf("expected secret/ count 30, got %d", result.ByMount["secret/"])
	}
}

func TestLeaseCount_HTTPError(t *testing.T) {
	srv := newLeaseCountTestServer(t, http.StatusInternalServerError, nil)
	defer srv.Close()

	checker := NewLeaseCountChecker(srv.URL, "test-token")
	_, err := checker.Count(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestLeaseCount_InvalidURL(t *testing.T) {
	checker := NewLeaseCountChecker("://bad-url", "token")
	_, err := checker.Count(context.Background())
	if err == nil {
		t.Fatal("expected error for invalid URL")
	}
}

func TestLeaseCount_Empty(t *testing.T) {
	payload := map[string]any{
		"data": map[string]any{
			"lease_count":    0,
			"count_per_mount": map[string]any{},
		},
	}
	srv := newLeaseCountTestServer(t, http.StatusOK, payload)
	defer srv.Close()

	checker := NewLeaseCountChecker(srv.URL, "test-token")
	result, err := checker.Count(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Total != 0 {
		t.Errorf("expected total 0, got %d", result.Total)
	}
}
