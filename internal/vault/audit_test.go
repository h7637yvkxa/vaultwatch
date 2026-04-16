package vault

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newAuditTestServer(t *testing.T, status int, payload any) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		if payload != nil {
			_ = json.NewEncoder(w).Encode(payload)
		}
	}))
}

func TestListAuditDevices_Success(t *testing.T) {
	body := map[string]any{
		"data": map[string]any{
			"file/": map[string]any{"type": "file", "description": "file audit"},
			"syslog/": map[string]any{"type": "syslog", "description": ""},
		},
	}
	srv := newAuditTestServer(t, http.StatusOK, body)
	defer srv.Close()

	c := newVaultClient(t, srv.URL)
	checker := NewAuditChecker(c)

	entries, err := checker.ListAuditDevices(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	for _, e := range entries {
		if !e.Enabled {
			t.Errorf("expected entry to be enabled")
		}
		if e.CheckedAt.IsZero() {
			t.Errorf("expected CheckedAt to be set")
		}
	}
}

func TestListAuditDevices_Empty(t *testing.T) {
	body := map[string]any{"data": map[string]any{}}
	srv := newAuditTestServer(t, http.StatusOK, body)
	defer srv.Close()

	c := newVaultClient(t, srv.URL)
	checker := NewAuditChecker(c)

	entries, err := checker.ListAuditDevices(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(entries))
	}
}

func TestListAuditDevices_HTTPError(t *testing.T) {
	srv := newAuditTestServer(t, http.StatusForbidden, nil)
	defer srv.Close()

	c := newVaultClient(t, srv.URL)
	checker := NewAuditChecker(c)

	_, err := checker.ListAuditDevices(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
