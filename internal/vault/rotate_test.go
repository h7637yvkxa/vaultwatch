package vault

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func newRotateTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPut && r.URL.Path == "/v1/sys/revoke":
			w.WriteHeader(http.StatusNoContent)
		case r.Method == http.MethodGet && r.URL.Path == "/v1/secret/data/myapp/db":
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"lease_id":       "new-lease-id-123",
				"lease_duration": 3600,
				"renewable":      true,
				"data":           map[string]interface{}{"username": "app"},
			})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
}

func TestRotateLease_EmptyLeaseID(t *testing.T) {
	svr := newRotateTestServer(t)
	defer svr.Close()
	c := newVaultClient(t, svr.URL)
	rot := NewRotator(c)
	res := rot.RotateLease(context.Background(), "", "secret/data/myapp/db")
	if res.Err == nil {
		t.Fatal("expected error for empty leaseID")
	}
}

func TestRotateLease_EmptyPath(t *testing.T) {
	svr := newRotateTestServer(t)
	defer svr.Close()
	c := newVaultClient(t, svr.URL)
	rot := NewRotator(c)
	res := rot.RotateLease(context.Background(), "old-lease", "")
	if res.Err == nil {
		t.Fatal("expected error for empty secretPath")
	}
}

func TestRotateLease_Success(t *testing.T) {
	svr := newRotateTestServer(t)
	defer svr.Close()
	c := newVaultClient(t, svr.URL)
	rot := NewRotator(c)
	before := time.Now()
	res := rot.RotateLease(context.Background(), "old-lease", "secret/data/myapp/db")
	if res.Err != nil {
		t.Fatalf("unexpected error: %v", res.Err)
	}
	if res.NewLeaseID != "new-lease-id-123" {
		t.Errorf("expected new-lease-id-123, got %s", res.NewLeaseID)
	}
	if res.RenewedAt.Before(before) {
		t.Error("RenewedAt should be after test start")
	}
}

func TestRotateExpiring_SkipsOK(t *testing.T) {
	svr := newRotateTestServer(t)
	defer svr.Close()
	c := newVaultClient(t, svr.URL)
	rot := NewRotator(c)
	statuses := []LeaseStatus{
		{LeaseID: "lease-ok", Level: LevelOK},
	}
	results := rot.RotateExpiring(context.Background(), statuses, map[string]string{})
	if len(results) != 0 {
		t.Errorf("expected 0 results for OK leases, got %d", len(results))
	}
}

func TestRotateExpiring_MissingPathMap(t *testing.T) {
	svr := newRotateTestServer(t)
	defer svr.Close()
	c := newVaultClient(t, svr.URL)
	rot := NewRotator(c)
	statuses := []LeaseStatus{
		{LeaseID: "lease-warn", Level: LevelWarning},
	}
	results := rot.RotateExpiring(context.Background(), statuses, map[string]string{})
	if len(results) != 1 || results[0].Err == nil {
		t.Error("expected error when path not in map")
	}
}
