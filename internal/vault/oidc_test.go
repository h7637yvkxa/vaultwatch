package vault_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/morningconsult/vaultwatch/internal/vault"
)

func newOIDCTestServer(t *testing.T, status int, keys []string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		if status == http.StatusOK {
			_ = json.NewEncoder(w).Encode(map[string]any{
				"data": map[string]any{"keys": keys},
			})
		}
	}))
}

func TestListOIDCRoles_Success(t *testing.T) {
	srv := newOIDCTestServer(t, http.StatusOK, []string{"web", "mobile"})
	defer srv.Close()

	c := newVaultClient(t, srv.URL)
	checker := vault.NewOIDCChecker(c)

	roles, err := checker.ListRoles()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(roles) != 2 {
		t.Fatalf("expected 2 roles, got %d", len(roles))
	}
	if roles[0].Name != "web" {
		t.Errorf("expected web, got %s", roles[0].Name)
	}
}

func TestListOIDCRoles_Empty(t *testing.T) {
	srv := newOIDCTestServer(t, http.StatusNotFound, nil)
	defer srv.Close()

	c := newVaultClient(t, srv.URL)
	checker := vault.NewOIDCChecker(c)

	roles, err := checker.ListRoles()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(roles) != 0 {
		t.Fatalf("expected 0 roles, got %d", len(roles))
	}
}

func TestListOIDCRoles_HTTPError(t *testing.T) {
	srv := newOIDCTestServer(t, http.StatusInternalServerError, nil)
	defer srv.Close()

	c := newVaultClient(t, srv.URL)
	checker := vault.NewOIDCChecker(c)

	_, err := checker.ListRoles()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
