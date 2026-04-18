package vault

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newPluginTestServer(t *testing.T, status int, plugins []PluginEntry) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/sys/plugins/catalog" {
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(status)
		if status == http.StatusOK {
			body := map[string]interface{}{
				"data": map[string]interface{}{
					"detailed": plugins,
				},
			}
			_ = json.NewEncoder(w).Encode(body)
		}
	}))
}

func TestListPlugins_Success(t *testing.T) {
	plugins := []PluginEntry{
		{Name: "aws", Type: "secret", Version: "v1.0.0", Builtin: true},
		{Name: "my-plugin", Type: "auth", Version: "v0.2.1", Builtin: false},
	}
	srv := newPluginTestServer(t, http.StatusOK, plugins)
	defer srv.Close()

	checker := NewPluginChecker(srv.URL, "test-token", srv.Client())
	result, err := checker.ListPlugins(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 plugins, got %d", len(result))
	}
	if result[0].Name != "aws" {
		t.Errorf("expected aws, got %s", result[0].Name)
	}
}

func TestListPlugins_Empty(t *testing.T) {
	srv := newPluginTestServer(t, http.StatusOK, []PluginEntry{})
	defer srv.Close()

	checker := NewPluginChecker(srv.URL, "test-token", srv.Client())
	result, err := checker.ListPlugins(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected empty, got %d", len(result))
	}
}

func TestListPlugins_HTTPError(t *testing.T) {
	srv := newPluginTestServer(t, http.StatusForbidden, nil)
	defer srv.Close()

	checker := NewPluginChecker(srv.URL, "bad-token", srv.Client())
	_, err := checker.ListPlugins(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestListPlugins_InvalidURL(t *testing.T) {
	checker := NewPluginChecker("http://127.0.0.1:0", "token", nil)
	_, err := checker.ListPlugins(context.Background())
	if err == nil {
		t.Fatal("expected error for invalid URL")
	}
}
