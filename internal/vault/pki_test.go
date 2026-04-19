package vault_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yourusername/vaultwatch/internal/vault"
)

func newPKITestServer(t *testing.T, status int, keys []string) *httptest.Server {
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

func TestListCerts_Success(t *testing.T) {
	serials := []string{"1a:2b:3c", "4d:5e:6f"}
	srv := newPKITestServer(t, http.StatusOK, serials)
	defer srv.Close()

	checker := vault.NewPKIChecker(srv.URL, "test-token")
	certs, err := checker.ListCerts("pki")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(certs) != 2 {
		t.Fatalf("expected 2 certs, got %d", len(certs))
	}
	if certs[0].SerialNumber != "1a:2b:3c" {
		t.Errorf("expected serial 1a:2b:3c, got %s", certs[0].SerialNumber)
	}
}

func TestListCerts_Empty(t *testing.T) {
	srv := newPKITestServer(t, http.StatusNotFound, nil)
	defer srv.Close()

	checker := vault.NewPKIChecker(srv.URL, "test-token")
	certs, err := checker.ListCerts("pki")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(certs) != 0 {
		t.Errorf("expected 0 certs, got %d", len(certs))
	}
}

func TestListCerts_HTTPError(t *testing.T) {
	srv := newPKITestServer(t, http.StatusInternalServerError, nil)
	defer srv.Close()

	checker := vault.NewPKIChecker(srv.URL, "test-token")
	_, err := checker.ListCerts("pki")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestListCerts_InvalidURL(t *testing.T) {
	checker := vault.NewPKIChecker("http://127.0.0.1:0", "token")
	_, err := checker.ListCerts("pki")
	if err == nil {
		t.Fatal("expected error for invalid URL")
	}
}
