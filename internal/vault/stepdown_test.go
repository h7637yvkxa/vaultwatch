package vault

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newStepDownTestServer(statusCode int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		w.WriteHeader(statusCode)
	}))
}

func TestStepDown_Success(t *testing.T) {
	srv := newStepDownTestServer(http.StatusNoContent)
	defer srv.Close()

	c := newVaultClient(srv.URL, "test-token")
	checker := NewStepDownChecker(c)

	result, err := checker.StepDown(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Errorf("expected success, got: %s", result.Message)
	}
}

func TestStepDown_HTTPError(t *testing.T) {
	srv := newStepDownTestServer(http.StatusForbidden)
	defer srv.Close()

	c := newVaultClient(srv.URL, "test-token")
	checker := NewStepDownChecker(c)

	result, err := checker.StepDown(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Success {
		t.Error("expected failure, got success")
	}
	if result.Message == "" {
		t.Error("expected non-empty message on failure")
	}
}

func TestStepDown_InvalidURL(t *testing.T) {
	c := newVaultClient("http://127.0.0.1:0", "tok")
	checker := NewStepDownChecker(c)

	_, err := checker.StepDown(context.Background())
	if err == nil {
		t.Error("expected error for unreachable server")
	}
}
