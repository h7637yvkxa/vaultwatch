package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func newKVMetaTestServer(t *testing.T, status int, payload any) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		if payload != nil {
			_ = json.NewEncoder(w).Encode(payload)
		}
	}))
}

func TestReadMetadata_Success(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	payload := map[string]any{
		"data": map[string]any{
			"current_version":      3,
			"oldest_version":       1,
			"created_time":         now.Format(time.RFC3339),
			"updated_time":         now.Format(time.RFC3339),
			"max_versions":         10,
			"delete_version_after": "0s",
		},
	}
	srv := newKVMetaTestServer(t, http.StatusOK, payload)
	defer srv.Close()

	c := NewKVMetadataChecker(srv.URL, "test-token")
	meta, err := c.ReadMetadata("secret", "myapp/db")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if meta.CurrentVersion != 3 {
		t.Errorf("expected version 3, got %d", meta.CurrentVersion)
	}
	if meta.MaxVersions != 10 {
		t.Errorf("expected max_versions 10, got %d", meta.MaxVersions)
	}
	if meta.Path != "secret/myapp/db" {
		t.Errorf("unexpected path: %s", meta.Path)
	}
}

func TestReadMetadata_HTTPError(t *testing.T) {
	srv := newKVMetaTestServer(t, http.StatusForbidden, nil)
	defer srv.Close()

	c := NewKVMetadataChecker(srv.URL, "bad-token")
	_, err := c.ReadMetadata("secret", "myapp/db")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestReadMetadata_InvalidURL(t *testing.T) {
	c := NewKVMetadataChecker("http://127.0.0.1:0", "tok")
	_, err := c.ReadMetadata("secret", "path")
	if err == nil {
		t.Fatal("expected error for invalid URL")
	}
}
