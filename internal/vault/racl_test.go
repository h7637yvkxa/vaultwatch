package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newRACLTestServer(t *testing.T, status int, rules []RACLEntry) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		if status == http.StatusOK {
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"rules": rules})
		}
	}))
}

func TestListACLPaths_Success(t *testing.T) {
	rules := []RACLEntry{
		{Path: "secret/*", Capabilities: []string{"read", "list"}},
		{Path: "auth/token/create", Capabilities: []string{"create"}},
	}
	srv := newRACLTestServer(t, http.StatusOK, rules)
	defer srv.Close()

	c := newVaultClient(t, srv.URL)
	checker := NewRACLChecker(c)

	entries, err := checker.ListACLPaths("default")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Path != "secret/*" {
		t.Errorf("expected path 'secret/*', got %q", entries[0].Path)
	}
}

func TestListACLPaths_EmptyPolicy(t *testing.T) {
	c := newVaultClient(t, "http://127.0.0.1")
	checker := NewRACLChecker(c)

	_, err := checker.ListACLPaths("")
	if err == nil {
		t.Fatal("expected error for empty policy name")
	}
}

func TestListACLPaths_HTTPError(t *testing.T) {
	srv := newRACLTestServer(t, http.StatusForbidden, nil)
	defer srv.Close()

	c := newVaultClient(t, srv.URL)
	checker := NewRACLChecker(c)

	_, err := checker.ListACLPaths("restricted")
	if err == nil {
		t.Fatal("expected error on non-200 status")
	}
}

func TestListACLPaths_InvalidURL(t *testing.T) {
	c := newVaultClient(t, "http://127.0.0.1:0")
	checker := NewRACLChecker(c)

	_, err := checker.ListACLPaths("default")
	if err == nil {
		t.Fatal("expected error for unreachable server")
	}
}
