package vault_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arjunsriva/vaultwatch/internal/vault"
)

func newTOTPTestServer(t *testing.T, status int, keys []string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if status != http.StatusOK {
			w.WriteHeader(status)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{"keys": keys},
		})
	}))
}

func TestListTOTPKeys_Success(t *testing.T) {
	srv := newTOTPTestServer(t, http.StatusOK, []string{"my-key", "other-key"})
	defer srv.Close()

	c := vault.NewTOTPChecker(srv.Client(), srv.URL, "test-token")
	keys, err := c.ListKeys("totp")
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

func TestListTOTPKeys_Empty(t *testing.T) {
	srv := newTOTPTestServer(t, http.StatusNotFound, nil)
	defer srv.Close()

	c := vault.NewTOTPChecker(srv.Client(), srv.URL, "test-token")
	keys, err := c.ListKeys("totp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(keys) != 0 {
		t.Errorf("expected empty, got %d", len(keys))
	}
}

func TestListTOTPKeys_HTTPError(t *testing.T) {
	srv := newTOTPTestServer(t, http.StatusInternalServerError, nil)
	defer srv.Close()

	c := vault.NewTOTPChecker(srv.Client(), srv.URL, "test-token")
	_, err := c.ListKeys("totp")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestListTOTPKeys_InvalidURL(t *testing.T) {
	c := vault.NewTOTPChecker(http.DefaultClient, "http://127.0.0.1:0", "tok")
	_, err := c.ListKeys("totp")
	if err == nil {
		t.Fatal("expected error for invalid URL")
	}
}
