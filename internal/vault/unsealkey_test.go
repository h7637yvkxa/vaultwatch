package vault_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/yourusername/vaultwatch/internal/vault"
)

func newUnsealKeyTestServer(t *testing.T, status int, body interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/sys/rekey/init" {
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(status)
		if body != nil {
			_ = json.NewEncoder(w).Encode(body)
		}
	}))
}

func TestGetUnsealKeyStatus_Success(t *testing.T) {
	payload := map[string]interface{}{
		"secret_shares":    5,
		"secret_threshold": 3,
		"stored_shares":    0,
		"nonce":            "abc123",
		"pgp_fingerprints": []string{"fp1", "fp2"},
	}
	srv := newUnsealKeyTestServer(t, http.StatusOK, payload)
	defer srv.Close()

	checker := vault.NewUnsealKeyChecker(srv.URL, "test-token")
	result, err := checker.GetUnsealKeyStatus(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.SecretShares != 5 {
		t.Errorf("expected SecretShares=5, got %d", result.SecretShares)
	}
	if result.SecretThreshold != 3 {
		t.Errorf("expected SecretThreshold=3, got %d", result.SecretThreshold)
	}
	if len(result.PGPFingerprints) != 2 {
		t.Errorf("expected 2 PGP fingerprints, got %d", len(result.PGPFingerprints))
	}
}

func TestGetUnsealKeyStatus_HTTPError(t *testing.T) {
	srv := newUnsealKeyTestServer(t, http.StatusInternalServerError, nil)
	defer srv.Close()

	checker := vault.NewUnsealKeyChecker(srv.URL, "test-token")
	_, err := checker.GetUnsealKeyStatus(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "unexpected status 500") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestGetUnsealKeyStatus_InvalidURL(t *testing.T) {
	checker := vault.NewUnsealKeyChecker("http://127.0.0.1:0", "token")
	_, err := checker.GetUnsealKeyStatus(context.Background())
	if err == nil {
		t.Fatal("expected error for unreachable server")
	}
}
