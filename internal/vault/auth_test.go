package vault

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newAuthTestServer(t *testing.T, status int, body any) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		if body != nil {
			_ = json.NewEncoder(w).Encode(body)
		}
	}))
}

func TestListAuthMethods_Success(t *testing.T) {
	body := map[string]any{
		"data": map[string]any{
			"token/": map[string]any{
				"type": "token", "description": "token auth",
				"accessor": "abc123", "local": false, "seal_wrap": false,
			},
			"approle/": map[string]any{
				"type": "approle", "description": "",
				"accessor": "def456", "local": true, "seal_wrap": false,
			},
		},
	}
	srv := newAuthTestServer(t, http.StatusOK, body)
	defer srv.Close()

	checker := NewAuthChecker(srv.Client(), srv.URL, "test-token")
	methods, err := checker.ListAuthMethods(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(methods) != 2 {
		t.Fatalf("expected 2 methods, got %d", len(methods))
	}
}

func TestListAuthMethods_Empty(t *testing.T) {
	body := map[string]any{"data": map[string]any{}}
	srv := newAuthTestServer(t, http.StatusOK, body)
	defer srv.Close()

	checker := NewAuthChecker(srv.Client(), srv.URL, "test-token")
	methods, err := checker.ListAuthMethods(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(methods) != 0 {
		t.Fatalf("expected 0 methods, got %d", len(methods))
	}
}

func TestListAuthMethods_HTTPError(t *testing.T) {
	srv := newAuthTestServer(t, http.StatusForbidden, nil)
	defer srv.Close()

	checker := NewAuthChecker(srv.Client(), srv.URL, "bad-token")
	_, err := checker.ListAuthMethods(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
