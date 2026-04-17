package alert

import (
	"bytes"
	"strings"
	"testing"

	"github.com/rzp-labs/vaultwatch/internal/vault"
)

func makeAuthMethod(path, typ, accessor string, local bool) vault.AuthMethod {
	return vault.AuthMethod{
		Path:     path,
		Type:     typ,
		Accessor: accessor,
		Local:    local,
	}
}

func TestAuthNotifier_NoMethods(t *testing.T) {
	var buf bytes.Buffer
	n := NewAuthNotifier(&buf)
	if err := n.Notify(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "no auth methods") {
		t.Errorf("expected no-methods message, got: %s", buf.String())
	}
}

func TestAuthNotifier_WithMethods(t *testing.T) {
	var buf bytes.Buffer
	n := NewAuthNotifier(&buf)
	methods := []vault.AuthMethod{
		makeAuthMethod("token/", "token", "abc123", false),
		makeAuthMethod("approle/", "approle", "def456", true),
	}
	if err := n.Notify(methods); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "2 auth method") {
		t.Errorf("expected count in output, got: %s", out)
	}
	if !strings.Contains(out, "token/") {
		t.Errorf("expected token/ in output")
	}
	if !strings.Contains(out, "(local)") {
		t.Errorf("expected local marker for approle")
	}
}

func TestAuthNotifier_WithDescription(t *testing.T) {
	var buf bytes.Buffer
	n := NewAuthNotifier(&buf)
	m := makeAuthMethod("github/", "github", "ghi789", false)
	m.Description = "GitHub auth"
	if err := n.Notify([]vault.AuthMethod{m}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "GitHub auth") {
		t.Errorf("expected description in output")
	}
}

func TestAuthNotifier_NilWriter_UsesStdout(t *testing.T) {
	n := NewAuthNotifier(nil)
	if n.w == nil {
		t.Error("expected non-nil writer")
	}
}
