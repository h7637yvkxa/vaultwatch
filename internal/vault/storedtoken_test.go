package vault

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newStoredTokenTestServer(t *testing.T, status int, keys []string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("list") != "true" {
			http.Error(w, "missing list param", http.StatusBadRequest)
			return
		}
		w.WriteHeader(status)
		if status == http.StatusOK {
			_ = json.NewEncoder(w).Encode(map[string]any{
				"data": map[string]any{"keys": keys},
			})
		}
	}))
}

func TestListStoredTokens_Success(t *testing.T) {
	srv := newStoredTokenTestServer(t, http.StatusOK, []string{"abc123", "def456"})
	defer srv.Close()

	c := NewStoredTokenChecker(srv.Client(), srv.URL, "test-token")
	res, err := c.List(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(res.Entries))
	}
	if res.Entries[0].Accessor != "abc123" {
		t.Errorf("expected accessor abc123, got %s", res.Entries[0].Accessor)
	}
}

func TestListStoredTokens_Empty(t *testing.T) {
	srv := newStoredTokenTestServer(t, http.StatusNotFound, nil)
	defer srv.Close()

	c := NewStoredTokenChecker(srv.Client(), srv.URL, "test-token")
	res, err := c.List(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(res.Entries))
	}
}

func TestListStoredTokens_HTTPError(t *testing.T) {
	srv := newStoredTokenTestServer(t, http.StatusInternalServerError, nil)
	defer srv.Close()

	c := NewStoredTokenChecker(srv.Client(), srv.URL, "test-token")
	_, err := c.List(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestListStoredTokens_InvalidURL(t *testing.T) {
	c := NewStoredTokenChecker(http.DefaultClient, "http://127.0.0.1:0", "tok")
	_, err := c.List(context.Background())
	if err == nil {
		t.Fatal("expected error for invalid URL")
	}
}
