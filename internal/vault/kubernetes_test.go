package vault_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dkbrummitt/vaultwatch/internal/vault"
)

func newKubernetesTestServer(t *testing.T, status int, keys []string) *httptest.Server {
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

func TestListKubernetesRoles_Success(t *testing.T) {
	srv := newKubernetesTestServer(t, http.StatusOK, []string{"my-role", "dev-role"})
	defer srv.Close()

	c := vault.NewKubernetesChecker(srv.URL, "test-token")
	roles, err := c.ListKubernetesRoles("kubernetes")
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

func TestListKubernetesRoles_Empty(t *testing.T) {
	srv := newKubernetesTestServer(t, http.StatusNotFound, nil)
	defer srv.Close()

	c := vault.NewKubernetesChecker(srv.URL, "test-token")
	roles, err := c.ListKubernetesRoles("kubernetes")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(roles) != 0 {
		t.Fatalf("expected 0 roles, got %d", len(roles))
	}
}

func TestListKubernetesRoles_HTTPError(t *testing.T) {
	srv := newKubernetesTestServer(t, http.StatusInternalServerError, nil)
	defer srv.Close()

	c := vault.NewKubernetesChecker(srv.URL, "test-token")
	_, err := c.ListKubernetesRoles("kubernetes")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestListKubernetesRoles_InvalidURL(t *testing.T) {
	c := vault.NewKubernetesChecker("http://127.0.0.1:0", "token")
	_, err := c.ListKubernetesRoles("kubernetes")
	if err == nil {
		t.Fatal("expected error for invalid URL")
	}
}
