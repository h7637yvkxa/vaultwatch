package alert

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/yourusername/vaultwatch/internal/vault"
)

func makeAuditEntry(path, typ, desc string) vault.AuditEntry {
	return vault.AuditEntry{
		Path:        path,
		Type:        typ,
		Description: desc,
		Enabled:     true,
		CheckedAt:   time.Now(),
	}
}

func TestAuditNotifier_NoDevices(t *testing.T) {
	var buf bytes.Buffer
	n := NewAuditNotifier(&buf)
	if err := n.Notify(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "no audit devices") {
		t.Errorf("expected no-devices message, got: %s", buf.String())
	}
}

func TestAuditNotifier_WithDevices(t *testing.T) {
	var buf bytes.Buffer
	n := NewAuditNotifier(&buf)
	entries := []vault.AuditEntry{
		makeAuditEntry("file/", "file", "file audit log"),
		makeAuditEntry("syslog/", "syslog", ""),
	}
	if err := n.Notify(entries); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "2 device(s)") {
		t.Errorf("expected count in output, got: %s", out)
	}
	if !strings.Contains(out, "file/") {
		t.Errorf("expected file/ in output, got: %s", out)
	}
	if !strings.Contains(out, "(no description)") {
		t.Errorf("expected fallback description, got: %s", out)
	}
}

func TestAuditNotifier_NilWriter_UsesStdout(t *testing.T) {
	n := NewAuditNotifier(nil)
	if n.w == nil {
		t.Error("expected fallback to stdout, got nil writer")
	}
}
