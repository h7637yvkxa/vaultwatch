package vault

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func newLicenseTestServer(t *testing.T, status int, body interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		if body != nil {
			_ = json.NewEncoder(w).Encode(body)
		}
	}))
}

func TestGetLicense_Success(t *testing.T) {
	expiry := time.Now().Add(30 * 24 * time.Hour).UTC().Truncate(time.Second)
	payload := map[string]interface{}{
		"data": map[string]interface{}{
			"license_id":      "abc-123",
			"customer_name":   "Acme Corp",
			"expiration_time": expiry.Format(time.RFC3339),
			"features":        []string{"Replication", "HSM"},
		},
	}
	ts := newLicenseTestServer(t, http.StatusOK, payload)
	defer ts.Close()

	c := newVaultClient(t, ts.URL)
	checker := NewLicenseChecker(c)
	info, err := checker.GetLicense(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info.LicenseID != "abc-123" {
		t.Errorf("license_id: got %q, want %q", info.LicenseID, "abc-123")
	}
	if info.CustomerName != "Acme Corp" {
		t.Errorf("customer_name: got %q, want %q", info.CustomerName, "Acme Corp")
	}
	if len(info.Features) != 2 {
		t.Errorf("features: got %d, want 2", len(info.Features))
	}
}

func TestGetLicense_HTTPError(t *testing.T) {
	ts := newLicenseTestServer(t, http.StatusForbidden, nil)
	defer ts.Close()

	c := newVaultClient(t, ts.URL)
	checker := NewLicenseChecker(c)
	_, err := checker.GetLicense(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetLicense_InvalidURL(t *testing.T) {
	c := newVaultClient(t, "://bad-url")
	checker := NewLicenseChecker(c)
	_, err := checker.GetLicense(context.Background())
	if err == nil {
		t.Fatal("expected error for invalid URL")
	}
}
