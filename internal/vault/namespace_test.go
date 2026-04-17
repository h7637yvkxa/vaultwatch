package vault

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newNamespaceTestServer(t *testing.T, status int, keys []string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if status == http.StatusNotFound {
			w.WriteHeader(http.StatusNotFound)
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

func TestListNamespaces_Success(t *testing.T) {
	srv := newNamespaceTestServer(t, http.StatusOK, []string{"team-a/", "team-b/"})
	defer srv.Close()

	nc := NewNamespaceChecker(srv.URL, "tok", nil)
	entries, err := nc.ListNamespaces(context.Background(), "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].ID != "team-a/" {
		t.Errorf("unexpected ID: %s", entries[0].ID)
	}
}

func TestListNamespaces_Empty(t *testing.T) {
	srv := newNamespaceTestServer(t, http.StatusNotFound, nil)
	defer srv.Close()

	nc := NewNamespaceChecker(srv.URL, "tok", nil)
	entries, err := nc.ListNamespaces(context.Background(), "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(entries))
	}
}

func TestListNamespaces_HTTPError(t *testing.T) {
	srv := newNamespaceTestServer(t, http.StatusInternalServerError, nil)
	defer srv.Close()

	nc := NewNamespaceChecker(srv.URL, "tok", nil)
	_, err := nc.ListNamespaces(context.Background(), "")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
