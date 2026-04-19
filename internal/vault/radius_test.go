package vault_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yourusername/vaultwatch/internal/vault"
)

func newRADIUSTestServer(t *testing.T, status int, keys []string) *httptest.Server {
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

func TestListRADIUSUsers_Success(t *testing.T) {
	srv := newRADIUSTestServer(t, http.StatusOK, []string{"alice", "bob"})
	defer srv.Close()

	c := vault.NewRADIUSChecker(srv.Client(), srv.URL, "tok", "radius")
	users, err := c.ListUsers()
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

func TestListRADIUSUsers_Empty(t *testing.T) {
	srv := newRADIUSTestServer(t, http.StatusNotFound, nil)
	defer srv.Close()

	c := vault.NewRADIUSChecker(srv.Client(), srv.URL, "tok", "radius")
	users, err := c.ListUsers()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(users) != 0 {
		t.Fatalf("expected 0 users, got %d", len(users))
	}
}

func TestListRADIUSUsers_HTTPError(t *testing.T) {
	srv := newRADIUSTestServer(t, http.StatusInternalServerError, nil)
	defer srv.Close()

	c := vault.NewRADIUSChecker(srv.Client(), srv.URL, "tok", "radius")
	_, err := c.ListUsers()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
