package vault_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yourusername/vaultwatch/internal/vault"
)

func newLDAPTestServer(t *testing.T, status int, keys []string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

func TestListLDAPGroups_Success(t *testing.T) {
	srv := newLDAPTestServer(t, http.StatusOK, []string{"admins", "developers"})
	defer srv.Close()

	c := vault.NewLDAPChecker(srv.Client(), srv.URL, "tok", "ldap")
	groups, err := c.ListGroups()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(groups) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(groups))
	}
	if groups[0].Name != "admins" {
		t.Errorf("expected admins, got %s", groups[0].Name)
	}
}

func TestListLDAPGroups_Empty(t *testing.T) {
	srv := newLDAPTestServer(t, http.StatusNotFound, nil)
	defer srv.Close()

	c := vault.NewLDAPChecker(srv.Client(), srv.URL, "tok", "ldap")
	groups, err := c.ListGroups()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(groups) != 0 {
		t.Errorf("expected empty, got %d", len(groups))
	}
}

func TestListLDAPGroups_HTTPError(t *testing.T) {
	srv := newLDAPTestServer(t, http.StatusInternalServerError, nil)
	defer srv.Close()

	c := vault.NewLDAPChecker(srv.Client(), srv.URL, "tok", "ldap")
	_, err := c.ListGroups()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
