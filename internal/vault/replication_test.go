package vault_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arjunsriva/vaultwatch/internal/vault"
)

func newReplicationTestServer(t *testing.T, status int, body any) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		if body != nil {
			_ = json.NewEncoder(w).Encode(body)
		}
	}))
}

func TestReplication_Success(t *testing.T) {
	body := map[string]any{
		"data": map[string]any{
			"dr":          map[string]any{"mode": "primary", "primary": true},
			"performance": map[string]any{"mode": "secondary", "primary": false},
		},
	}
	srv := newReplicationTestServer(t, http.StatusOK, body)
	defer srv.Close()

	c := vault.NewReplicationChecker(srv.URL, "tok", nil)
	st, err := c.Check(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if st.DRMode != "primary" {
		t.Errorf("expected dr_mode primary, got %s", st.DRMode)
	}
	if !st.DRPrimary {
		t.Error("expected DR primary true")
	}
	if st.PerformanceMode != "secondary" {
		t.Errorf("expected perf secondary, got %s", st.PerformanceMode)
	}
}

func TestReplication_HTTPError(t *testing.T) {
	srv := newReplicationTestServer(t, http.StatusForbidden, nil)
	defer srv.Close()

	c := vault.NewReplicationChecker(srv.URL, "bad", nil)
	_, err := c.Check(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestReplication_InvalidURL(t *testing.T) {
	c := vault.NewReplicationChecker("http://127.0.0.1:1", "tok", nil)
	_, err := c.Check(context.Background())
	if err == nil {
		t.Fatal("expected error for invalid URL")
	}
}
