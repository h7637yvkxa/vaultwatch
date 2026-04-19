package vault_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/densestvoid/vaultwatch/internal/vault"
)

func newCertStoreTestServer(t *testing.T, status int, keys []string) *httptest.Server {
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

func TestListCertStoreEntries_Success(t *testing.T) {
	srv := newCertStoreTestServer(t, http.StatusOK, []string{"web-cert", "api-cert"})
	defer srv.Close()

	c := vault.NewCertStoreChecker(srv.URL, "token", "cert")
	entries, err := c.ListCerts()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Name != "web-cert" {
		t.Errorf("expected web-cert, got %s", entries[0].Name)
	}
}

func TestListCertStoreEntries_Empty(t *testing.T) {
	srv := newCertStoreTestServer(t, http.StatusNotFound, nil)
	defer srv.Close()

	c := vault.NewCertStoreChecker(srv.URL, "token", "cert")
	entries, err := c.ListCerts()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(entries))
	}
}

func TestListCertStoreEntries_HTTPError(t *testing.T) {
	srv := newCertStoreTestServer(t, http.StatusInternalServerError, nil)
	defer srv.Close()

	c := vault.NewCertStoreChecker(srv.URL, "token", "cert")
	_, err := c.ListCerts()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestListCertStoreEntries_InvalidURL(t *testing.T) {
	c := vault.NewCertStoreChecker("http://127.0.0.1:1", "token", "cert")
	_, err := c.ListCerts()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
