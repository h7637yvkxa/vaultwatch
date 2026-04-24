package vault_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/danhale-git/vaultwatch/internal/vault"
)

func newTokenAccessorTestServer(t *testing.T, status int, keys []string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "LIST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		w.WriteHeader(status)
		if status == http.StatusOK {
			body := map[string]interface{}{
				"data": map[string]interface{}{"keys": keys},
			}
			_ = json.NewEncoder(w).Encode(body)
		}
	}))
}

func TestListTokenAccessors_Success(t *testing.T) {
	keys := []string{"accessor-aaa", "accessor-bbb"}
	srv := newTokenAccessorTestServer(t, http.StatusOK, keys)
	defer srv.Close()

	checker := vault.NewTokenAccessorChecker(srv.URL, "test-token")
	entries, err := checker.ListTokenAccessors(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Accessor != "accessor-aaa" {
		t.Errorf("expected accessor-aaa, got %s", entries[0].Accessor)
	}
}

func TestListTokenAccessors_Empty(t *testing.T) {
	srv := newTokenAccessorTestServer(t, http.StatusNotFound, nil)
	defer srv.Close()

	checker := vault.NewTokenAccessorChecker(srv.URL, "test-token")
	entries, err := checker.ListTokenAccessors(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(entries))
	}
}

func TestListTokenAccessors_HTTPError(t *testing.T) {
	srv := newTokenAccessorTestServer(t, http.StatusInternalServerError, nil)
	defer srv.Close()

	checker := vault.NewTokenAccessorChecker(srv.URL, "test-token")
	_, err := checker.ListTokenAccessors(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestListTokenAccessors_InvalidURL(t *testing.T) {
	checker := vault.NewTokenAccessorChecker("http://127.0.0.1:0", "test-token")
	_, err := checker.ListTokenAccessors(context.Background())
	if err == nil {
		t.Fatal("expected connection error, got nil")
	}
}
