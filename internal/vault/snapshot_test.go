package vault_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/yourusername/vaultwatch/internal/vault"
)

func newSnapshotServer(statusCode int, body string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/sys/storage/raft/snapshot" {
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(statusCode)
		_, _ = w.Write([]byte(body))
	}))
}

func TestTakeSnapshot_Success(t *testing.T) {
	payload := strings.Repeat("x", 512)
	srv := newSnapshotServer(http.StatusOK, payload)
	defer srv.Close()

	checker := vault.NewSnapshotChecker(srv.URL, "test-token")
	var buf bytes.Buffer
	meta, err := checker.TakeSnapshot(context.Background(), &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if meta.Size != int64(len(payload)) {
		t.Errorf("expected size %d, got %d", len(payload), meta.Size)
	}
	if buf.String() != payload {
		t.Errorf("unexpected snapshot content")
	}
	if meta.TakenAt.IsZero() {
		t.Errorf("expected non-zero TakenAt")
	}
}

func TestTakeSnapshot_HTTPError(t *testing.T) {
	srv := newSnapshotServer(http.StatusForbidden, "permission denied")
	defer srv.Close()

	checker := vault.NewSnapshotChecker(srv.URL, "bad-token")
	var buf bytes.Buffer
	_, err := checker.TakeSnapshot(context.Background(), &buf)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "403") {
		t.Errorf("expected 403 in error, got: %v", err)
	}
}

func TestTakeSnapshot_InvalidURL(t *testing.T) {
	checker := vault.NewSnapshotChecker("http://127.0.0.1:0", "token")
	var buf bytes.Buffer
	_, err := checker.TakeSnapshot(context.Background(), &buf)
	if err == nil {
		t.Fatal("expected connection error")
	}
}
