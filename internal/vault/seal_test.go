package vault

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newSealTestServer(t *testing.T, status SealStatus, code int) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/sys/seal-status" {
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(code)
		_ = json.NewEncoder(w).Encode(status)
	}))
}

func TestSealCheck_Unsealed(t *testing.T) {
	srv := newSealTestServer(t, SealStatus{Sealed: false, Initialized: true, Version: "1.14.0"}, http.StatusOK)
	defer srv.Close()
	checker := NewSealChecker(srv.URL, nil)
	got, err := checker.Check(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Sealed {
		t.Error("expected unsealed")
	}
	if got.Version != "1.14.0" {
		t.Errorf("version = %q, want 1.14.0", got.Version)
	}
}

func TestSealCheck_Sealed(t *testing.T) {
	srv := newSealTestServer(t, SealStatus{Sealed: true, Initialized: true}, http.StatusOK)
	defer srv.Close()
	checker := NewSealChecker(srv.URL, nil)
	got, err := checker.Check(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !got.Sealed {
		t.Error("expected sealed")
	}
}

func TestSealCheck_HTTPError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()
	checker := NewSealChecker(srv.URL, nil)
	_, err := checker.Check(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestSealCheck_InvalidURL(t *testing.T) {
	checker := NewSealChecker("http://127.0.0.1:0", nil)
	_, err := checker.Check(context.Background())
	if err == nil {
		t.Fatal("expected error for unreachable server")
	}
}
