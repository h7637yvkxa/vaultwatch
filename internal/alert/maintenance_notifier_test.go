package alert

import (
	"bytes"
	"strings"
	"testing"

	"github.com/wiliamsouza/vaultwatch/internal/vault"
)

func makeMaintenanceStatus(enabled bool, message string) *vault.MaintenanceStatus {
	return &vault.MaintenanceStatus{
		Enabled: enabled,
		Message: message,
	}
}

func TestMaintenanceNotifier_NilStatus(t *testing.T) {
	var buf bytes.Buffer
	n := NewMaintenanceNotifier(&buf)
	if err := n.Notify(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "no status available") {
		t.Errorf("expected 'no status available', got: %q", buf.String())
	}
}

func TestMaintenanceNotifier_Disabled(t *testing.T) {
	var buf bytes.Buffer
	n := NewMaintenanceNotifier(&buf)
	if err := n.Notify(makeMaintenanceStatus(false, "")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "disabled") {
		t.Errorf("expected 'disabled' in output, got: %q", out)
	}
}

func TestMaintenanceNotifier_Enabled(t *testing.T) {
	var buf bytes.Buffer
	n := NewMaintenanceNotifier(&buf)
	if err := n.Notify(makeMaintenanceStatus(true, "scheduled downtime")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "ENABLED") {
		t.Errorf("expected 'ENABLED' in output, got: %q", out)
	}
	if !strings.Contains(out, "scheduled downtime") {
		t.Errorf("expected message in output, got: %q", out)
	}
}

func TestMaintenanceNotifier_EnabledNoMessage(t *testing.T) {
	var buf bytes.Buffer
	n := NewMaintenanceNotifier(&buf)
	if err := n.Notify(makeMaintenanceStatus(true, "")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "no message provided") {
		t.Errorf("expected fallback message, got: %q", out)
	}
}

func TestMaintenanceNotifier_NilWriter_UsesStdout(t *testing.T) {
	n := NewMaintenanceNotifier(nil)
	if n.w == nil {
		t.Error("expected non-nil writer when nil passed")
	}
}
