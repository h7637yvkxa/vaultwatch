package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newUserpassTestServer(t *testing.T, status int, keys []string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if status == http.StatusNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		}
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

func TestListUserpassUsers_Success(t *testing.T) {
	srv := newUserpassTestServer(t, http.StatusOK, []string{"alice", "bob"})
	defer srv.Close()
	c := newVaultClient(t, srv.URL)
	checker := NewUserpassChecker(c, "")
	users, err := checker.ListUsers()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(users) != 2 {
		t.Fatalf("expected 2 users, got %d", len(users))
	}
	if users[0].Username != "alice" {
		t.Errorf("expected alice, got %s", users[0].Username)
	}
}

func TestListUserpassUsers_Empty(t *testing.T) {
	srv := newUserpassTestServer(t, http.StatusNotFound, nil)
	defer srv.Close()
	c := newVaultClient(t, srv.URL)
	checker := NewUserpassChecker(c, "userpass")
	users, err := checker.ListUsers()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(users) != 0 {
		t.Errorf("expected empty, got %d", len(users))
	}
}

func TestListUserpassUsers_HTTPError(t *testing.T) {
	srv := newUserpassTestServer(t, http.StatusInternalServerError, nil)
	defer srv.Close()
	c := newVaultClient(t, srv.URL)
	checker := NewUserpassChecker(c, "userpass")
	_, err := checker.ListUsers()
	if err == nil {
		t.Fatal("expected error")
	}
}
