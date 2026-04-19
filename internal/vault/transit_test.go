package vault_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/wes/vaultwatch/internal/vault"
)

func newTransitTestServer(t *testing.T, status int, keys []string) *httptest.Server {
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

func TestListTransitKeys_Success(t *testing.T) {
	srv := newTransitTestServer(t, http.StatusOK, []string{"my-key", "other-key"})
	defer srv.Close()

	c := vault.NewTransitChecker(srv.URL, "token", nil)
	keys, err := c.ListTransitKeys("transit")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(keys) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(keys))
	}
	if keys[0].Name != "my-key" {
		t.Errorf("expected my-key, got %s", keys[0].Name)
	}
}

func TestListTransitKeys_Empty(t *testing.T) {
	srv := newTransitTestServer(t, http.StatusNotFound, nil)
	defer srv.Close()

	c := vault.NewTransitChecker(srv.URL, "token", nil)
	keys, err := c.ListTransitKeys("transit")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(keys) != 0 {
		t.Errorf("expected empty, got %d", len(keys))
	}
}

func TestListTransitKeys_HTTPError(t *testing.T) {
	srv := newTransitTestServer(t, http.StatusInternalServerError, nil)
	defer srv.Close()

	c := vault.NewTransitChecker(srv.URL, "token", nil)
	_, err := c.ListTransitKeys("transit")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestListTransitKeys_InvalidURL(t *testing.T) {
	c := vault.NewTransitChecker("http://127.0.0.1:0", "token", nil)
	_, err := c.ListTransitKeys("transit")
	if err == nil {
		t.Fatal("expected error for invalid URL")
	}
}
