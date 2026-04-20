package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newRGPTestServer(t *testing.T, status int, keys []string) *httptest.Server {
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

func TestListRGPPolicies_Success(t *testing.T) {
	srv := newRGPTestServer(t, http.StatusOK, []string{"allow-read", "deny-write"})
	defer srv.Close()

	c := NewRGPChecker(srv.URL, "test-token")
	policies, err := c.ListRGPPolicies()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(policies) != 2 {
		t.Fatalf("expected 2 policies, got %d", len(policies))
	}
	if policies[0].Name != "allow-read" {
		t.Errorf("expected allow-read, got %s", policies[0].Name)
	}
}

func TestListRGPPolicies_Empty(t *testing.T) {
	srv := newRGPTestServer(t, http.StatusNotFound, nil)
	defer srv.Close()

	c := NewRGPChecker(srv.URL, "test-token")
	policies, err := c.ListRGPPolicies()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(policies) != 0 {
		t.Errorf("expected 0 policies, got %d", len(policies))
	}
}

func TestListRGPPolicies_HTTPError(t *testing.T) {
	srv := newRGPTestServer(t, http.StatusInternalServerError, nil)
	defer srv.Close()

	c := NewRGPChecker(srv.URL, "test-token")
	_, err := c.ListRGPPolicies()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestListRGPPolicies_InvalidURL(t *testing.T) {
	c := NewRGPChecker("http://127.0.0.1:0", "tok")
	_, err := c.ListRGPPolicies()
	if err == nil {
		t.Fatal("expected error for invalid URL")
	}
}
