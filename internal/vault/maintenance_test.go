package vault

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newMaintenanceTestServer(t *testing.T, status int, body interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		if body != nil {
			_ = json.NewEncoder(w).Encode(body)
		}
	}))
}

func TestMaintenanceCheck_Disabled(t *testing.T) {
	payload := map[string]interface{}{"enabled": false, "message": "", "request_id": "abc"}
	srv := newMaintenanceTestServer(t, http.StatusOK, payload)
	defer srv.Close()

	c := newVaultClient(t, srv.URL)
	checker := NewMaintenanceChecker(c)

	status, err := checker.Check(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if status.Enabled {
		t.Error("expected maintenance to be disabled")
	}
}

func TestMaintenanceCheck_Enabled(t *testing.T) {
	payload := map[string]interface{}{"enabled": true, "message": "scheduled downtime", "request_id": "xyz"}
	srv := newMaintenanceTestServer(t, http.StatusOK, payload)
	defer srv.Close()

	c := newVaultClient(t, srv.URL)
	checker := NewMaintenanceChecker(c)

	status, err := checker.Check(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !status.Enabled {
		t.Error("expected maintenance to be enabled")
	}
	if status.Message != "scheduled downtime" {
		t.Errorf("unexpected message: %q", status.Message)
	}
}

func TestMaintenanceCheck_HTTPError(t *testing.T) {
	srv := newMaintenanceTestServer(t, http.StatusInternalServerError, nil)
	defer srv.Close()

	c := newVaultClient(t, srv.URL)
	checker := NewMaintenanceChecker(c)

	_, err := checker.Check(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestMaintenanceCheck_InvalidURL(t *testing.T) {
	c := newVaultClient(t, "http://127.0.0.1:0")
	checker := NewMaintenanceChecker(c)

	_, err := checker.Check(context.Background())
	if err == nil {
		t.Fatal("expected error for unreachable server")
	}
}
