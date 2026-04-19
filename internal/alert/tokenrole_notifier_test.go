package alert_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourusername/vaultwatch/internal/alert"
	"github.com/yourusername/vaultwatch/internal/vault"
)

func makeTokenRole(name string, policies []string) vault.TokenRole {
	return vault.TokenRole{
		Name:            name,
		AllowedPolicies: policies,
		Renewable:       true,
		Orphan:          false,
		ExplicitMaxTTL:  3600,
	}
}

func TestTokenRoleNotifier_NoRoles(t *testing.T) {
	var buf bytes.Buffer
	n := alert.NewTokenRoleNotifier(&buf)
	if err := n.Notify(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "no token roles") {
		t.Errorf("expected 'no token roles', got: %s", buf.String())
	}
}

func TestTokenRoleNotifier_WithRoles(t *testing.T) {
	var buf bytes.Buffer
	n := alert.NewTokenRoleNotifier(&buf)
	roles := []vault.TokenRole{
		makeTokenRole("admin", []string{"default", "admin"}),
		makeTokenRole("readonly", nil),
	}
	if err := n.Notify(roles); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "admin") {
		t.Errorf("expected 'admin' in output")
	}
	if !strings.Contains(out, "readonly") {
		t.Errorf("expected 'readonly' in output")
	}
	if !strings.Contains(out, "(none)") {
		t.Errorf("expected '(none)' for role with no policies")
	}
}

func TestTokenRoleNotifier_NilWriter_UsesStdout(t *testing.T) {
	n := alert.NewTokenRoleNotifier(nil)
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}
