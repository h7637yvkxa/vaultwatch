package vault_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arussellsaw/vaultwatch/internal/vault"
)

func newGCPTestServer(t *testing.T, keys []string, status int) *httptest.Server {
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

func TestListGCPRoles_Success(t *testing.T) {
	srv := newGCPTestServer(t, []string{"my-role", "other-role"}, http.StatusOK)
	defer srv.Close()

	c := newVaultClient(t, srv.URL)
	checker := vault.NewGCPChecker(c, "gcp")
	roles, err := checker.ListRoles()
	if err != nil {
		t.Fatal(err)
	}
	if len(roles) != 2 {
		t.Fatalf("expected 2 roles, got %d", len(roles))
	}
	if roles[0].Name != "my-role" {
		t.Errorf("unexpected role name: %s", roles[0].Name)
	}
}

func TestListGCPRoles_Empty(t *testing.T) {
	srv := newGCPTestServer(t, nil, http.StatusNotFound)
	defer srv.Close()

	c := newVaultClient(t, srv.URL)
	checker := vault.NewGCPChecker(c, "gcp")
	roles, err := checker.ListRoles()
	if err != nil {
		t.Fatal(err)
	}
	if len(roles) != 0 {
		t.Fatalf("expected 0 roles, got %d", len(roles))
	}
}

func TestListGCPRoles_HTTPError(t *testing.T) {
	srv := newGCPTestServer(t, nil, http.StatusInternalServerError)
	defer srv.Close()

	c := newVaultClient(t, srv.URL)
	checker := vault.NewGCPChecker(c, "gcp")
	_, err := checker.ListRoles()
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestListGCPRoles_InvalidURL(t *testing.T) {
	c := newVaultClient(t, "http://127.0.0.1:1")
	checker := vault.NewGCPChecker(c, "gcp")
	_, err := checker.ListRoles()
	if err == nil {
		t.Fatal("expected error")
	}
}
