package vault

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newEGPTestServer(t *testing.T, status int, keys []string) *httptest.Server {
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

func TestListEGPPolicies_Success(t *testing.T) {
	srv := newEGPTestServer(t, http.StatusOK, []string{"allow-read", "deny-write"})
	defer srv.Close()

	checker := NewEGPChecker(srv.Client(), srv.URL, "test-token")
	policies, err := checker.ListEGPPolicies(context.Background())
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

func TestListEGPPolicies_Empty(t *testing.T) {
	srv := newEGPTestServer(t, http.StatusNotFound, nil)
	defer srv.Close()

	checker := NewEGPChecker(srv.Client(), srv.URL, "test-token")
	policies, err := checker.ListEGPPolicies(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(policies) != 0 {
		t.Fatalf("expected 0 policies, got %d", len(policies))
	}
}

func TestListEGPPolicies_HTTPError(t *testing.T) {
	srv := newEGPTestServer(t, http.StatusInternalServerError, nil)
	defer srv.Close()

	checker := NewEGPChecker(srv.Client(), srv.URL, "test-token")
	_, err := checker.ListEGPPolicies(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestListEGPPolicies_InvalidURL(t *testing.T) {
	checker := NewEGPChecker(http.DefaultClient, "http://127.0.0.1:0", "tok")
	_, err := checker.ListEGPPolicies(context.Background())
	if err == nil {
		t.Fatal("expected error for invalid URL")
	}
}
