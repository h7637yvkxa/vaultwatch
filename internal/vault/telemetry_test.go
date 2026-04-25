package vault

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTelemetryTestServer(t *testing.T, status int, body interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		if body != nil {
			_ = json.NewEncoder(w).Encode(body)
		}
	}))
}

func TestReadTelemetry_Success(t *testing.T) {
	payload := TelemetryResult{
		Gauges: []TelemetryMetric{
			{Name: "vault.core.active", Value: 1},
		},
		Counters: []TelemetryMetric{
			{Name: "vault.token.create", Value: 42},
		},
	}
	srv := newTelemetryTestServer(t, http.StatusOK, payload)
	defer srv.Close()

	c := newVaultClient(t, srv.URL)
	checker := NewTelemetryChecker(c)

	result, err := checker.ReadTelemetry(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Gauges) != 1 {
		t.Errorf("expected 1 gauge, got %d", len(result.Gauges))
	}
	if result.Gauges[0].Name != "vault.core.active" {
		t.Errorf("unexpected gauge name: %s", result.Gauges[0].Name)
	}
	if len(result.Counters) != 1 {
		t.Errorf("expected 1 counter, got %d", len(result.Counters))
	}
}

func TestReadTelemetry_HTTPError(t *testing.T) {
	srv := newTelemetryTestServer(t, http.StatusForbidden, nil)
	defer srv.Close()

	c := newVaultClient(t, srv.URL)
	checker := NewTelemetryChecker(c)

	_, err := checker.ReadTelemetry(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestReadTelemetry_InvalidURL(t *testing.T) {
	c := newVaultClient(t, "http://127.0.0.1:0")
	checker := NewTelemetryChecker(c)

	_, err := checker.ReadTelemetry(context.Background())
	if err == nil {
		t.Fatal("expected error for unreachable server")
	}
}

func TestReadTelemetry_Empty(t *testing.T) {
	payload := TelemetryResult{}
	srv := newTelemetryTestServer(t, http.StatusOK, payload)
	defer srv.Close()

	c := newVaultClient(t, srv.URL)
	checker := NewTelemetryChecker(c)

	result, err := checker.ReadTelemetry(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Gauges) != 0 || len(result.Counters) != 0 || len(result.Samples) != 0 {
		t.Error("expected empty telemetry result")
	}
}
