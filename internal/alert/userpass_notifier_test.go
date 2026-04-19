package alert

import (
	"bytes"
	"strings"
	"testing"

	"github.com/drew/vaultwatch/internal/vault"
)

func makeUserpassUser(name string, policies []string) vault.UserpassUser {
	return vault.UserpassUser{Username: name, Policies: policies}
}

func TestUserpassNotifier_NoUsers(t *testing.T) {
	var buf bytes.Buffer
	n := NewUserpassNotifier(&buf)
	if err := n.Notify(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "no users") {
		t.Errorf("expected 'no users', got: %s", buf.String())
	}
}

func TestUserpassNotifier_WithUsers(t *testing.T) {
	var buf bytes.Buffer
	n := NewUserpassNotifier(&buf)
	users := []vault.UserpassUser{
		makeUserpassUser("alice", []string{"admin"}),
		makeUserpassUser("bob", nil),
	}
	if err := n.Notify(users); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "alice") {
		t.Errorf("expected alice in output")
	}
	if !strings.Contains(out, "bob") {
		t.Errorf("expected bob in output")
	}
	if !strings.Contains(out, "admin") {
		t.Errorf("expected admin policy in output")
	}
	if !strings.Contains(out, "(none)") {
		t.Errorf("expected (none) for bob")
	}
}

func TestUserpassNotifier_NilWriter_UsesStdout(t *testing.T) {
	n := NewUserpassNotifier(nil)
	if n.w == nil {
		t.Error("expected non-nil writer")
	}
}
