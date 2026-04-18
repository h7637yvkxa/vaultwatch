package vault

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newLoginTestServer(t *testing.T, status int, payload any) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		if payload != nil {
			_ = json.NewEncoder(w).Encode(payload)
		}
	}))
}

func TestLookupToken_Success(t *testing.T) {
	payload := map[string]any{
		"data": map[string]any{
			"id":            "tok-abc",
			"accessor":      "acc-123",
			"policies":      []string{"default", "admin"},
			"ttl":           3600,
			"renewable":     true,
			"creation_time": 1700000000,
		},
	}
	srv := newLoginTestServer(t, http.StatusOK, payload)
	defer srv.Close()

	c := newVaultClient(t, srv.URL, "test-token")
	lc := NewLoginChecker(c)
	info, err := lc.LookupToken(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info.ClientToken != "tok-abc" {
		t.Errorf("expected tok-abc, got %s", info.ClientToken)
	}
	if !info.Renewable {
		t.Error("expected renewable")
	}
	if len(info.Policies) != 2 {
		t.Errorf("expected 2 policies, got %d", len(info.Policies))
	}
}

func TestLookupToken_HTTPError(t *testing.T) {
	srv := newLoginTestServer(t, http.StatusForbidden, nil)
	defer srv.Close()

	c := newVaultClient(t, srv.URL, "bad-token")
	lc := NewLoginChecker(c)
	_, err := lc.LookupToken(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestLookupToken_InvalidURL(t *testing.T) {
	c := newVaultClient(t, "http://127.0.0.1:0", "tok")
	lc := NewLoginChecker(c)
	_, err := lc.LookupToken(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
}
