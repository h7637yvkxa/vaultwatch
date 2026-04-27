package vault

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newRequestCountTestServer(t *testing.T, status int, payload any) (*httptest.Server, *http.Client) {
	t.Helper()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/sys/internal/counters/requests" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_ = json.NewEncoder(w).Encode(payload)
	}))
	t.Cleanup(ts.Close)
	return ts, ts.Client()
}

func TestGetRequestCounts_Success(t *testing.T) {
	payload := map[string]any{
		"data": map[string]any{
			"start_time": "2024-01-01T00:00:00Z",
			"end_time":   "2024-01-31T23:59:59Z",
			"total":      42,
			"by_namespace": []map[string]any{
				{"namespace_id": "root", "namespace_path": "", "counts": 42},
			},
		},
	}
	ts, client := newRequestCountTestServer(t, http.StatusOK, payload)
	checker := NewRequestCountChecker(ts.URL, "test-token", client)

	result, err := checker.GetRequestCounts(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Total != 42 {
		t.Errorf("expected total 42, got %d", result.Total)
	}
	if len(result.ByNamespace) != 1 {
		t.Errorf("expected 1 namespace entry, got %d", len(result.ByNamespace))
	}
	if result.ByNamespace[0].NamespaceID != "root" {
		t.Errorf("expected namespace_id 'root', got %q", result.ByNamespace[0].NamespaceID)
	}
}

func TestGetRequestCounts_HTTPError(t *testing.T) {
	ts, client := newRequestCountTestServer(t, http.StatusInternalServerError, map[string]any{})
	checker := NewRequestCountChecker(ts.URL, "test-token", client)

	_, err := checker.GetRequestCounts(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetRequestCounts_InvalidURL(t *testing.T) {
	checker := NewRequestCountChecker("://bad-url", "token", http.DefaultClient)
	_, err := checker.GetRequestCounts(context.Background())
	if err == nil {
		t.Fatal("expected error for invalid URL")
	}
}

func TestGetRequestCounts_Empty(t *testing.T) {
	payload := map[string]any{
		"data": map[string]any{
			"total":        0,
			"by_namespace": []map[string]any{},
		},
	}
	ts, client := newRequestCountTestServer(t, http.StatusOK, payload)
	checker := NewRequestCountChecker(ts.URL, "test-token", client)

	result, err := checker.GetRequestCounts(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Total != 0 {
		t.Errorf("expected total 0, got %d", result.Total)
	}
	if len(result.ByNamespace) != 0 {
		t.Errorf("expected 0 namespace entries, got %d", len(result.ByNamespace))
	}
}
