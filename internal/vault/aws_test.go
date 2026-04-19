package vault_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rbhz/vaultwatch/internal/vault"
)

func newAWSTestServer(t *testing.T, status int, keys []string) *httptest.Server {
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

func TestListAWSRoles_Success(t *testing.T) {
	srv := newAWSTestServer(t, http.StatusOK, []string{"role1", "role2"})
	defer srv.Close()

	c := vault.NewAWSChecker(srv.URL, "token", "")
	roles, err := c.ListAWSRoles()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(roles) != 2 {
		t.Fatalf("expected 2 roles, got %d", len(roles))
	if roles[0].Name != "role1" {
		t.Errorf("expected role1, got %s", roles[0].Name)
	}
}

func TestListAWSRoles_Empty(t *testing.T) {
	srv := newAWSTestServer(t, http.StatusNotFound, nil)
	defer srv.Close()

	c := vault.NewAWSChecker(srv.URL, "token", "")
	roles, err := c.ListAWSRoles()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(roles) != 0 {
		t.Fatalf("expected 0 roles, got %d", len(roles))
	}
}

func TestListAWSRoles_HTTPError(t *testing.T) {
	srv := newAWSTestServer(t, http.StatusInternalServerError, nil)
	defer srv.Close()

	c := vault.NewAWSChecker(srv.URL, "token", "")
	_, err := c.ListAWSRoles()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestListAWSRoles_InvalidURL(t *testing.T) {
	c := vault.NewAWSChecker("http://127.0.0.1:0", "token", "")
	_, err := c.ListAWSRoles()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
