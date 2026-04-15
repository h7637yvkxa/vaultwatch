package alert

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/drengskapur/vaultwatch/internal/vault"
)

func TestRotateNotifier_NoResults(t *testing.T) {
	var buf bytes.Buffer
	n := NewRotateNotifier(&buf)
	if err := n.Notify(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "no leases rotated") {
		t.Errorf("expected 'no leases rotated', got: %s", buf.String())
	}
}

func TestRotateNotifier_AllSuccess(t *testing.T) {
	var buf bytes.Buffer
	n := NewRotateNotifier(&buf)
	at := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	results := []vault.RotateResult{
		{LeaseID: "old-1", NewLeaseID: "new-1", RenewedAt: at},
		{LeaseID: "old-2", NewLeaseID: "new-2", RenewedAt: at},
	}
	if err := n.Notify(results); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	for _, want := range []string{"old-1", "new-1", "old-2", "new-2", "2024-06-01T12:00:00Z"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected output to contain %q, got:\n%s", want, out)
		}
	}
}

func TestRotateNotifier_WithErrors(t *testing.T) {
	var buf bytes.Buffer
	n := NewRotateNotifier(&buf)
	results := []vault.RotateResult{
		{LeaseID: "bad-lease", Err: errors.New("revoke failed")},
		{LeaseID: "old-ok", NewLeaseID: "new-ok", RenewedAt: time.Now()},
	}
	if err := n.Notify(results); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "FAILED") {
		t.Errorf("expected FAILED in output, got: %s", out)
	}
	if !strings.Contains(out, "revoke failed") {
		t.Errorf("expected error text in output, got: %s", out)
	}
	if !strings.Contains(out, "new-ok") {
		t.Errorf("expected successful rotation in output, got: %s", out)
	}
}

func TestRotateNotifier_NilWriter_UsesStdout(t *testing.T) {
	// Ensure constructor does not panic with nil writer.
	n := NewRotateNotifier(nil)
	if n.w == nil {
		t.Error("expected non-nil writer when nil passed to NewRotateNotifier")
	}
}
