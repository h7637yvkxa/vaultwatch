package vault_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arcanericky/vaultwatch/internal/vault"
)

func newDatabaseTestServer(t *testing.T, status int, keys []string) *httptest.Server {
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

func TestListDatabaseRoles_Success(t *testing.T) {
	srv := newDatabaseTestServer(t, http.StatusOK, []string{"readonly", "readwrite"})
	defer srv.Close()

	checker := vault.NewDatabaseChecker(srv.URL, "token", "")
	roles, err := checker.ListDatabaseRoles()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(roles) != 2 {
		t.Fatalf("expected 2 roles, got %d", len(roles))
	}
	if roles[0].Name != "readonly" {
		t.Errorf("expected readonly, got %s", roles[0].Name)
	}
}

func TestListDatabaseRoles_Empty(t *testing.T) {
	srv := newDatabaseTestServer(t, http.StatusNotFound, nil)
	defer srv.Close()

	checker := vault.NewDatabaseChecker(srv.URL, "token", "database")
	roles, err := checker.ListDatabaseRoles()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(roles) != 0 {
		t.Fatalf("expected 0 roles, got %d", len(roles))
	}
}

func TestListDatabaseRoles_HTTPError(t *testing.T) {
	srv := newDatabaseTestServer(t, http.StatusInternalServerError, nil)
	defer srv.Close()

	checker := vault.NewDatabaseChecker(srv.URL, "token", "")
	_, err := checker.ListDatabaseRoles()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestListDatabaseRoles_InvalidURL(t *testing.T) {
	checker := vault.NewDatabaseChecker("://bad-url", "token", "")
	_, err := checker.ListDatabaseRoles()
	if err == nil {
		t.Fatal("expected error for invalid URL")
	}
}
