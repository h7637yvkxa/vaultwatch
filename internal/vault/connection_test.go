package vault_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yourusername/vaultwatch/internal/vault"
)

func newConnectionTestServer(statusCode int, clusterName, version string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"cluster_name": clusterName,
			"version":      version,
		})
	}))
}

func TestConnectionCheck_Reachable(t *testing.T) {
	srv := newConnectionTestServer(http.StatusOK, "vault-cluster", "1.15.0")
	defer srv.Close()

	checker := vault.NewConnectionChecker(srv.URL, "test-token", nil)
	status, err := checker.Check(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !status.Reachable {
		t.Error("expected reachable=true")
	}
	if status.ClusterName != "vault-cluster" {
		t.Errorf("expected cluster_name=vault-cluster, got %s", status.ClusterName)
	}
	if status.Version != "1.15.0" {
		t.Errorf("expected version=1.15.0, got %s", status.Version)
	}
	if status.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", status.StatusCode)
	}
}

func TestConnectionCheck_Unreachable(t *testing.T) {
	checker := vault.NewConnectionChecker("http://127.0.0.1:19999", "", nil)
	status, err := checker.Check(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if status.Reachable {
		t.Error("expected reachable=false")
	}
	if status.Error == "" {
		t.Error("expected non-empty error message")
	}
}

func TestConnectionCheck_NonOKStatus(t *testing.T) {
	srv := newConnectionTestServer(http.StatusServiceUnavailable, "", "")
	defer srv.Close()

	checker := vault.NewConnectionChecker(srv.URL, "", nil)
	status, err := checker.Check(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !status.Reachable {
		t.Error("expected reachable=true even on non-200")
	}
	if status.StatusCode != http.StatusServiceUnavailable {
		t.Errorf("expected 503, got %d", status.StatusCode)
	}
}
