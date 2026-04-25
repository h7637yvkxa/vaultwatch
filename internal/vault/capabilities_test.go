package vault_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yourusername/vaultwatch/internal/vault"
)

func newCapabilitiesTestServer(t *testing.T, payload map[string]interface{}, status int) *httptest.Server {
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

func TestCheckCapabilities_Success(t *testing.T) {
	payload := map[string]interface{}{
		"secret/data/foo": []interface{}{"read", "list"},
		"secret/data/bar": []interface{}{"deny"},
	}
	srv := newCapabilitiesTestServer(t, payload, http.StatusOK)
	defer srv.Close()

	checker := vault.NewCapabilityChecker(srv.URL, "test-token", nil)
	results, err := checker.CheckCapabilities(context.Background(), []string{"secret/data/foo", "secret/data/bar"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if results[0].Path != "secret/data/foo" {
		t.Errorf("expected path secret/data/foo, got %s", results[0].Path)
	}
	if len(results[0].Capabilities) != 2 {
		t.Errorf("expected 2 caps, got %d", len(results[0].Capabilities))
	}
}

func TestCheckCapabilities_Empty(t *testing.T) {
	checker := vault.NewCapabilityChecker("http://localhost", "tok", nil)
	results, err := checker.CheckCapabilities(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected empty results")
	}
}

func TestCheckCapabilities_HTTPError(t *testing.T) {
	srv := newCapabilitiesTestServer(t, nil, http.StatusForbidden)
	defer srv.Close()

	checker := vault.NewCapabilityChecker(srv.URL, "bad-token", nil)
	_, err := checker.CheckCapabilities(context.Background(), []string{"secret/data/foo"})
	if err == nil {
		t.Fatal("expected error for non-200 status")
	}
}

func TestCheckCapabilities_InvalidURL(t *testing.T) {
	checker := vault.NewCapabilityChecker("://bad-url", "tok", nil)
	_, err := checker.CheckCapabilities(context.Background(), []string{"secret/"})
	if err == nil {
		t.Fatal("expected error for invalid URL")
	}
}
