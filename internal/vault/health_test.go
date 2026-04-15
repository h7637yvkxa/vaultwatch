package vault_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/yourusername/vaultwatch/internal/vault"
)

func newHealthServer(statusCode int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/sys/health" {
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(statusCode)
	}))
}

func TestCheck_Healthy(t *testing.T) {
	srv := newHealthServer(http.StatusOK)
	defer srv.Close()

	checker := vault.NewChecker(srv.URL, 5*time.Second)
	status, err := checker.Check(context.Background())

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !status.Initialized {
		t.Error("expected Initialized=true")
	}
	if status.Sealed {
		t.Error("expected Sealed=false")
	}
	if status.Standby {
		t.Error("expected Standby=false")
	}
}

func TestCheck_Standby(t *testing.T) {
	srv := newHealthServer(http.StatusTooManyRequests)
	defer srv.Close()

	checker := vault.NewChecker(srv.URL, 5*time.Second)
	status, err := checker.Check(context.Background())

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !status.Standby {
		t.Error("expected Standby=true")
	}
}

func TestCheck_Sealed(t *testing.T) {
	srv := newHealthServer(http.StatusServiceUnavailable)
	defer srv.Close()

	checker := vault.NewChecker(srv.URL, 5*time.Second)
	status, err := checker.Check(context.Background())

	if err == nil {
		t.Fatal("expected an error for sealed vault")
	}
	if !status.Sealed {
		t.Error("expected Sealed=true")
	}
}

func TestCheck_Uninitialized(t *testing.T) {
	srv := newHealthServer(http.StatusNotImplemented)
	defer srv.Close()

	checker := vault.NewChecker(srv.URL, 5*time.Second)
	_, err := checker.Check(context.Background())

	if err == nil {
		t.Fatal("expected an error for uninitialized vault")
	}
}

func TestCheck_UnreachableServer(t *testing.T) {
	checker := vault.NewChecker("http://127.0.0.1:19999", 500*time.Millisecond)
	_, err := checker.Check(context.Background())

	if err == nil {
		t.Fatal("expected an error for unreachable server")
	}
}
