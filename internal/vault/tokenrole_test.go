package vault_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yourusername/vaultwatch/internal/vault"
)

func newTokenRoleTestServer(t *testing.T, status int, keys []string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if status == http.StatusNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if status != http.StatusOK {
			w.WriteHeader(status)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{"keys": keys},
		})
	}))
}

func TestListTokenRoles_Success(t *testing.T) {
	srv := newTokenRoleTestServer(t, http.StatusOK, []string{"admin", "readonly"})
	defer srv.Close()

	c := vault.NewTokenRoleChecker(srv.URL, "test-token")
	roles, err := c.ListTokenRoles()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(roles) != 2 {
		t.Fatalf("expected 2 roles, got %d", len(roles))
	}
	if roles[0].Name != "admin" {
		t.Errorf("expected admin, got %s", roles[0].Name)
	}
}

func TestListTokenRoles_Empty(t *testing.T) {
	srv := newTokenRoleTestServer(t, http.StatusNotFound, nil)
	defer srv.Close()

	c := vault.NewTokenRoleChecker(srv.URL, "test-token")
	roles, err := c.ListTokenRoles()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(roles) != 0 {
		t.Fatalf("expected 0 roles, got %d", len(roles))
	}
}

func TestListTokenRoles_HTTPError(t *testing.T) {
	srv := newTokenRoleTestServer(t, http.StatusInternalServerError, nil)
	defer srv.Close()

	c := vault.NewTokenRoleChecker(srv.URL, "test-token")
	_, err := c.ListTokenRoles()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestListTokenRoles_InvalidURL(t *testing.T) {
	c := vault.NewTokenRoleChecker("http://127.0.0.1:0", "tok")
	_, err := c.ListTokenRoles()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
