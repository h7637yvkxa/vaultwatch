package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newRoleTestServer(t *testing.T, status int, keys []string) *httptest.Server {
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

func TestListRoles_Success(t *testing.T) {
	srv := newRoleTestServer(t, http.StatusOK, []string{"my-role", "other-role"})
	defer srv.Close()

	c := newVaultClient(t, srv.URL)
	checker := NewRoleChecker(c)
	roles, err := checker.ListRoles("approle")
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

func TestListRoles_Empty(t *testing.T) {
	srv := newRoleTestServer(t, http.StatusNotFound, nil)
	defer srv.Close()

	c := newVaultClient(t, srv.URL)
	checker := NewRoleChecker(c)
	roles, err := checker.ListRoles("approle")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(roles) != 0 {
		t.Errorf("expected empty, got %d", len(roles))
	}
}

func TestListRoles_HTTPError(t *testing.T) {
	srv := newRoleTestServer(t, http.StatusInternalServerError, nil)
	defer srv.Close()

	c := newVaultClient(t, srv.URL)
	checker := NewRoleChecker(c)
	_, err := checker.ListRoles("approle")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestListRoles_InvalidURL(t *testing.T) {
	c := newVaultClient(t, "http://127.0.0.1:0")
	checker := NewRoleChecker(c)
	_, err := checker.ListRoles("approle")
	if err == nil {
		t.Fatal("expected error for invalid URL")
	}
}
