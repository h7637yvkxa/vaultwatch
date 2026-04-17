package alert_test

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/yourusername/vaultwatch/internal/alert"
)

func makeSnapshotResult(size int64, err error) *alert.SnapshotResult {
	return &alert.SnapshotResult{
		TakenAt: time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
		Size:    size,
		Err:     err,
	}
}

func TestSnapshotNotifier_NilResult(t *testing.T) {
	var buf bytes.Buffer
	n := alert.NewSnapshotNotifier(&buf)
	if err := n.Notify(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "no result") {
		t.Errorf("expected 'no result', got: %s", buf.String())
	}
}

func TestSnapshotNotifier_Success(t *testing.T) {
	var buf bytes.Buffer
	n := alert.NewSnapshotNotifier(&buf)
	if err := n.Notify(makeSnapshotResult(4096, nil)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "OK") {
		t.Errorf("expected OK, got: %s", out)
	}
	if !strings.Contains(out, "4096") {
		t.Errorf("expected size 4096, got: %s", out)
	}
}

func TestSnapshotNotifier_Failure(t *testing.T) {
	var buf bytes.Buffer
	n := alert.NewSnapshotNotifier(&buf)
	if err := n.Notify(makeSnapshotResult(0, errors.New("disk full"))); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "FAILED") {
		t.Errorf("expected FAILED, got: %s", out)
	}
	if !strings.Contains(out, "disk full") {
		t.Errorf("expected error message, got: %s", out)
	}
}

func TestSnapshotNotifier_NilWriter_UsesStdout(t *testing.T) {
	n := alert.NewSnapshotNotifier(nil)
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}
