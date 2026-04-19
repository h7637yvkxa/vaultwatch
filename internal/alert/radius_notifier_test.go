package alert_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourusername/vaultwatch/internal/alert"
	"github.com/yourusername/vaultwatch/internal/vault"
)

func makeRADIUSUser(name string, policies []string) vault.RADIUSUser {
	return vault.RADIUSUser{Username: name, Policies: policies}
}

func TestRADIUSNotifier_NoUsers(t *testing.T) {
	var buf bytes.Buffer
	n := alert.NewRADIUSNotifier(&buf)
	if err := n.Notify(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No users") {
		t.Errorf("expected 'No users', got %q", buf.String())
	}
}

func TestRADIUSNotifier_WithUsers(t *testing.T) {
	var buf bytes.Buffer
	n := alert.NewRADIUSNotifier(&buf)
	users := []vault.RADIUSUser{
		makeRADIUSUser("alice", []string{"default"}),
		makeRADIUSUser("bob", nil),
	}
	if err := n.Notify(users); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "alice") {
		t.Errorf("expected alice in output, got %q", out)
	}
	if !strings.Contains(out, "bob") {
		t.Errorf("expected bob in output, got %q", out)
	}
	if !strings.Contains(out, "(none)") {
		t.Errorf("expected (none) for bob, got %q", out)
	}
}

func TestRADIUSNotifier_NilWriter_UsesStdout(t *testing.T) {
	n := alert.NewRADIUSNotifier(nil)
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}
