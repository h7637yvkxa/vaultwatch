package vault

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newSysConfigTestServer(t *testing.T, status int, payload interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		if payload != nil {
			_ = json.NewEncoder(w).Encode(payload)
		}
	}))
}

func TestReadSysConfig_Success(t *testing.T) {
	payload := map[string]interface{}{
		"data": map[string]interface{}{
			"default_lease_ttl": "768h",
			"max_lease_ttl":     "8760h",
			"force_no_cache":    false,
		},
	}
	srv := newSysConfigTestServer(t, http.StatusOK, payload)
	defer srv.Close()

	c := newVaultClient(t, srv.URL)
	checker := NewSysConfigChecker(c)
	cfg, err := checker.ReadSysConfig(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.DefaultLeaseTTL != "768h" {
		t.Errorf("expected default_lease_ttl=768h, got %s", cfg.DefaultLeaseTTL)
	}
	if cfg.MaxLeaseTTL != "8760h" {
		t.Errorf("expected max_lease_ttl=8760h, got %s", cfg.MaxLeaseTTL)
	}
	if cfg.ForceNoCache {
		t.Error("expected force_no_cache=false")
	}
}

func TestReadSysConfig_HTTPError(t *testing.T) {
	srv := newSysConfigTestServer(t, http.StatusForbidden, nil)
	defer srv.Close()

	c := newVaultClient(t, srv.URL)
	checker := NewSysConfigChecker(c)
	_, err := checker.ReadSysConfig(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestReadSysConfig_InvalidURL(t *testing.T) {
	c := newVaultClient(t, "http://127.0.0.1:0")
	checker := NewSysConfigChecker(c)
	_, err := checker.ReadSysConfig(context.Background())
	if err == nil {
		t.Fatal("expected error for unreachable server")
	}
}
