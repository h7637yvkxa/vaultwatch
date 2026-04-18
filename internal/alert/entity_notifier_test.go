package alert_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/elizabethwanjiku703/vaultwatch/internal/alert"
	"github.com/elizabethwanjiku703/vaultwatch/internal/vault"
)

func makeEntity(name string, disabled bool) vault.EntityEntry {
	return vault.EntityEntry{
		ID:       "id-" + name,
		Name:     name,
		Policies: []string{"default"},
		Disabled: disabled,
	}
}

func TestEntityNotifier_NoEntries(t *testing.T) {
	var buf bytes.Buffer
	n := alert.NewEntityNotifier(&buf)
	if err := n.Notify(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "no identity entities") {
		t.Errorf("expected no-entities message, got: %s", buf.String())
	}
}

func TestEntityNotifier_WithEntries(t *testing.T) {
	var buf bytes.Buffer
	n := alert.NewEntityNotifier(&buf)
	entries := []vault.EntityEntry{makeEntity("alice", false), makeEntity("bob", true)}
	if err := n.Notify(entries); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "alice") {
		t.Errorf("expected alice in output")
	}
	if !strings.Contains(out, "disabled") {
		t.Errorf("expected disabled status in output")
	}
}

func TestEntityNotifier_NilWriter_UsesStdout(t *testing.T) {
	n := alert.NewEntityNotifier(nil)
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}
