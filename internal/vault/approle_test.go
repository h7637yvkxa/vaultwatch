package vault

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newAppRoleTestServer(t *testing.T, status int, keys []string) *httptest.Server {
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

func TestListAppRoles_Success(t *testing.T) {
	srv := newAppRoleTestServer(t, http.StatusOK, []string{"my-role", "other-role"})
	defer srv.Close()

	c := NewAppRoleChecker(srv.URL, "tok", "approle")
	entries, err := c.ListAppRoles(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Name != "my-role" {
		t.Errorf("expected my-role, got %s", entries[0].Name)
	}
}

func TestListAppRoles_Empty(t *testing.T) {
	srv := newAppRoleTestServer(t, http.StatusNotFound, nil)
	defer srv.Close()

	c := NewAppRoleChecker(srv.URL, "tok", "approle")
	entries, err := c.ListAppRoles(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected empty, got %d", len(entries))
	}
}

func TestListAppRoles_HTTPError(t *testing.T) {
	srv := newAppRoleTestServer(t, http.StatusForbidden, nil)
	defer srv.Close()

	c := NewAppRoleChecker(srv.URL, "tok", "approle")
	_, err := c.ListAppRoles(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestListAppRoles_InvalidURL(t *testing.T) {
	c := NewAppRoleChecker("http://127.0.0.1:0", "tok", "approle")
	_, err := c.ListAppRoles(context.Background())
	if err == nil {
		t.Fatal("expected error for unreachable server")
	}
}
