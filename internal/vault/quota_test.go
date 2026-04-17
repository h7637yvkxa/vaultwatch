package vault

import (
/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newQuotaTestServer(t *testing.T, keys []string, entries map[string]QuotaEntry) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("list") == "true" {
			json.NewEncoder(w).Encode(map[string]interface{}{
				"data": map[string]interface{}{"keys": keys},
			})
			return
		}
		name := r.URL.Path[len("/v1/sys/quotas/rate-limit/"):]
		if entry, ok := entries[name]; ok {
			json.NewEncoder(w).Encode(map[string]interface{}{"data": entry})
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
}

func TestListQuotas_Empty(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	checker := NewQuotaChecker(srv.URL, "token", srv.Client())
	entries, err := checker.ListQuotas()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(entries))
	}
}

func TestListQuotas_Success(t *testing.T) {
	keys := []string{"global-rl", "kv-rl"}
	entries := map[string]QuotaEntry{
		"global-rl": {Name: "global-rl", Type: "rate-limit", Rate: 100, Interval: 1},
		"kv-rl":     {Name: "kv-rl", Type: "rate-limit", Path: "secret/", Rate: 50, Interval: 1},
	}
	srv := newQuotaTestServer(t, keys, entries)
	defer srv.Close()

	checker := NewQuotaChecker(srv.URL, "token", srv.Client())
	result, err := checker.ListQuotas()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result))
	}
	if result[0].Name != "global-rl" {
		t.Errorf("expected global-rl, got %s", result[0].Name)
	}
}

func TestListQuotas_HTTPError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	checker := NewQuotaChecker(srv.URL, "token", srv.Client())
	_, err := checker.ListQuotas()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestListQuotas_InvalidURL(t *testing.T) {
	checker := NewQuotaChecker("http://127.0.0.1:0", "token", nil)
	_, err := checker.ListQuotas()
	if err == nil {
		t.Fatal("expected error for invalid URL")
	}
}
