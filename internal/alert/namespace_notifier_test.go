package alert

import (
	"bytes"
	"strings"
	"testing"

	"github.com/has/vaultwatch/internal/vault"
)

func makeNamespaceEntry(path, id string) vault.NamespaceEntry {
	return vault.NamespaceEntry{Path: path, ID: id}
}

func TestNamespaceNotifier_NoEntries(t *testing.T) {
	var buf bytes.Buffer
	n := NewNamespaceNotifier(&buf)
	if err := n.Notify(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "no namespaces found") {
		t.Errorf("expected 'no namespaces found', got: %s", buf.String())
	}
}

func TestNamespaceNotifier_WithEntries(t *testing.T) {
	var buf bytes.Buffer
	n := NewNamespaceNotifier(&buf)
	entries := []vault.NamespaceEntry{
		makeNamespaceEntry("team-a/", "team-a/"),
		makeNamespaceEntry("team-b/", "team-b/"),
	}
	if err := n.Notify(entries); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "2 namespace(s)") {
		t.Errorf("expected count in output, got: %s", out)
	}
	if !strings.Contains(out, "team-a/") {
		t.Errorf("expected team-a/ in output, got: %s", out)
	}
}

func TestNamespaceNotifier_NilWriter_UsesStdout(t *testing.T) {
	n := NewNamespaceNotifier(nil)
	if n.w == nil {
		t.Error("expected non-nil writer")
	}
}
