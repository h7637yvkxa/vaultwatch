package vault_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arcanericky/vaultwatch/internal/vault"
)

func newHAStateTestServer(t *testing.T, status int, body interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		if body != nil {
			_ = json.NewEncoder(w).Encode(body)
		}
	}))
}

func TestHAState_Success(t *testing.T) {
	payload := map[string]interface{}{
		"ha_enabled":     true,
		"is_self":        true,
		"leader_address": "https://vault.example.com:8200",
		"active_time":    "2024-01-01T00:00:00Z",
	}
	srv := newHAStateTestServer(t, http.StatusOK, payload)
	defer srv.Close()

	checker := vault.NewHAStateChecker(srv.URL, "test-token")
	state, err := checker.Check(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !state.HAEnabled {
		t.Error("expected HA to be enabled")
	}
	if state.LeaderAddress != "https://vault.example.com:8200" {
		t.Errorf("unexpected leader address: %s", state.LeaderAddress)
	}
}

func TestHAState_HTTPError(t *testing.T) {
	srv := newHAStateTestServer(t, http.StatusInternalServerError, nil)
	defer srv.Close()

	checker := vault.NewHAStateChecker(srv.URL, "test-token")
	_, err := checker.Check(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestHAState_InvalidURL(t *testing.T) {
	checker := vault.NewHAStateChecker("http://127.0.0.1:0", "test-token")
	_, err := checker.Check(context.Background())
	if err == nil {
		t.Fatal("expected error for invalid URL")
	}
}
