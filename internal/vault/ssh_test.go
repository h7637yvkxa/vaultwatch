package vault_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arcanericky/vaultwatch/internal/vault"
)

func newSSHTestServer(t *testing.T, status int, keys []string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if status != http.StatusOK {
			w.WriteHeader(status)
			return
		}
		body := map[string]interface{}{
			"data": map[string]interface{}{"keys": keys},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(body)
	}))
}

func TestListSSHRoles_Success(t *testing.T) {
	srv := newSSHTestServer(t, http.StatusOK, []string{"my-role", "dev-role"})
	defer srv.Close()

	c := vault.NewSSHChecker(srv.URL, "token")
	roles, err := c.ListRoles("ssh")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(roles) != 2 {
		t.Fatalf("expected 2 roles, got %d", len(roles))
	}
	if roles[0].Name != "my-role" {
		t.Errorf("expected my-role, got %s", roles[0].Name)
	}
}

func TestListSSHRoles_Empty(t *testing.T) {
	srv := newSSHTestServer(t, http.StatusNotFound, nil)
	defer srv.Close()

	c := vault.NewSSHChecker(srv.URL, "token")
	roles, err := c.ListRoles("ssh")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(roles) != 0 {
		t.Errorf("expected 0 roles, got %d", len(roles))
	}
}

func TestListSSHRoles_HTTPError(t *testing.T) {
	srv := newSSHTestServer(t, http.StatusInternalServerError, nil)
	defer srv.Close()

	c := vault.NewSSHChecker(srv.URL, "token")
	_, err := c.ListRoles("ssh")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestListSSHRoles_InvalidURL(t *testing.T) {
	c := vault.NewSSHChecker("http://127.0.0.1:0", "token")
	_, err := c.ListRoles("ssh")
	if err == nil {
		t.Fatal("expected error for invalid URL")
	}
}
