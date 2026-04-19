package vault_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arturhoo/vaultwatch/internal/vault"
)

func newAzureTestServer(t *testing.T, keys []string, status int) *httptest.Server {
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

func TestListAzureRoles_Success(t *testing.T) {
	srv := newAzureTestServer(t, []string{"role-a", "role-b"}, http.StatusOK)
	defer srv.Close()

	c := newVaultClient(t, srv.URL)
	checker := vault.NewAzureChecker(c)
	roles, err := checker.ListAzureRoles("azure")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(roles) != 2 {
		t.Fatalf("expected 2 roles, got %d", len(roles))
	}
	if roles[0].Name != "role-a" {
		t.Errorf("expected role-a, got %s", roles[0].Name)
	}
}

func TestListAzureRoles_Empty(t *testing.T) {
	srv := newAzureTestServer(t, nil, http.StatusNotFound)
	defer srv.Close()

	c := newVaultClient(t, srv.URL)
	checker := vault.NewAzureChecker(c)
	roles, err := checker.ListAzureRoles("azure")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(roles) != 0 {
		t.Fatalf("expected 0 roles, got %d", len(roles))
	}
}

func TestListAzureRoles_HTTPError(t *testing.T) {
	srv := newAzureTestServer(t, nil, http.StatusInternalServerError)
	defer srv.Close()

	c := newVaultClient(t, srv.URL)
	checker := vault.NewAzureChecker(c)
	_, err := checker.ListAzureRoles("azure")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestListAzureRoles_InvalidURL(t *testing.T) {
	c := newVaultClient(t, "http://127.0.0.1:0")
	checker := vault.NewAzureChecker(c)
	_, err := checker.ListAzureRoles("azure")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
