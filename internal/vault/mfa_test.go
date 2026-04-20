package vault_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arjunsriva/vaultwatch/internal/vault"
)

func newMFATestServer(t *testing.T, status int, payload any) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		if payload != nil {
			_ = json.NewEncoder(w).Encode(payload)
		}
	}))
}

func TestListMFAMethods_Success(t *testing.T) {
	payload := map[string]any{
		"data": map[string]any{
			"keys": []string{"abc123"},
			"key_info": map[string]any{
				"abc123": map[string]any{
					"id":   "abc123",
					"name": "my-totp",
					"type": "totp",
				},
			},
		},
	}
	srv := newMFATestServer(t, http.StatusOK, payload)
	defer srv.Close()

	c := newVaultClient(t, srv.URL)
	checker := vault.NewMFAChecker(c)
	methods, err := checker.ListMFAMethods(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(methods) != 1 {
		t.Fatalf("expected 1 method, got %d", len(methods))
	}
	if methods[0].Type != "totp" {
		t.Errorf("expected type totp, got %s", methods[0].Type)
	}
}

func TestListMFAMethods_Empty(t *testing.T) {
	payload := map[string]any{
		"data": map[string]any{
			"keys":     []string{},
			"key_info": map[string]any{},
		},
	}
	srv := newMFATestServer(t, http.StatusOK, payload)
	defer srv.Close()

	c := newVaultClient(t, srv.URL)
	checker := vault.NewMFAChecker(c)
	methods, err := checker.ListMFAMethods(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(methods) != 0 {
		t.Errorf("expected 0 methods, got %d", len(methods))
	}
}

func TestListMFAMethods_HTTPError(t *testing.T) {
	srv := newMFATestServer(t, http.StatusForbidden, nil)
	defer srv.Close()

	c := newVaultClient(t, srv.URL)
	checker := vault.NewMFAChecker(c)
	_, err := checker.ListMFAMethods(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
