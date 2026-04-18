package vault_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/elizabethwanjiku703/vaultwatch/internal/vault"
)

func newEntityTestServer(t *testing.T, status int, payload any) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		if payload != nil {
			_ = json.NewEncoder(w).Encode(payload)
		}
	}))
}

func TestListEntities_Success(t *testing.T) {
	payload := map[string]any{
		"data": map[string]any{
			"key_info": map[string]any{
				"abc-123": map[string]any{"name": "alice", "policies": []string{"default"}, "disabled": false},
				"def-456": map[string]any{"name": "bob", "policies": []string{"admin"}, "disabled": true},
			},
		},
	}
	srv := newEntityTestServer(t, http.StatusOK, payload)
	defer srv.Close()

	c := newVaultClient(t, srv.URL)
	checker := vault.NewEntityChecker(c)
	entries, err := checker.ListEntities(t.Context())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
}

func TestListEntities_Empty(t *testing.T) {
	srv := newEntityTestServer(t, http.StatusNotFound, nil)
	defer srv.Close()

	c := newVaultClient(t, srv.URL)
	checker := vault.NewEntityChecker(c)
	entries, err := checker.ListEntities(t.Context())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(entries))
	}
}

func TestListEntities_HTTPError(t *testing.T) {
	srv := newEntityTestServer(t, http.StatusInternalServerError, nil)
	defer srv.Close()

	c := newVaultClient(t, srv.URL)
	checker := vault.NewEntityChecker(c)
	_, err := checker.ListEntities(t.Context())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
