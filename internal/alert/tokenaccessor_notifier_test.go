package alert_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/danhale-git/vaultwatch/internal/alert"
	"github.com/danhale-git/vaultwatch/internal/vault"
)

func makeAccessorEntry(accessor, display string, ttl int) vault.TokenAccessorEntry {
	return vault.TokenAccessorEntry{
		Accessor:    accessor,
		DisplayName: display,
		TTL:         ttl,
	}
}

func TestTokenAccessorNotifier_NoEntries(t *testing.T) {
	var buf bytes.Buffer
	n := alert.NewTokenAccessorNotifier(&buf)
	if err := n.Notify(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "no accessors found") {
		t.Errorf("expected 'no accessors found', got: %s", buf.String())
	}
}

func TestTokenAccessorNotifier_WithEntries(t *testing.T) {
	var buf bytes.Buffer
	n := alert.NewTokenAccessorNotifier(&buf)
	entries := []vault.TokenAccessorEntry{
		makeAccessorEntry("acc-111", "root", 3600),
		makeAccessorEntry("acc-222", "app-token", 1800),
	}
	if err := n.Notify(entries); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "2 accessor(s) active") {
		t.Errorf("expected count in output, got: %s", out)
	}
	if !strings.Contains(out, "acc-111") {
		t.Errorf("expected acc-111 in output, got: %s", out)
	}
	if !strings.Contains(out, "app-token") {
		t.Errorf("expected app-token in output, got: %s", out)
	}
}

func TestTokenAccessorNotifier_NoDisplayName(t *testing.T) {
	var buf bytes.Buffer
	n := alert.NewTokenAccessorNotifier(&buf)
	entries := []vault.TokenAccessorEntry{
		makeAccessorEntry("acc-333", "", 900),
	}
	if err := n.Notify(entries); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "(no display name)") {
		t.Errorf("expected fallback display name, got: %s", buf.String())
	}
}

func TestTokenAccessorNotifier_NilWriter_UsesStdout(t *testing.T) {
	n := alert.NewTokenAccessorNotifier(nil)
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}
