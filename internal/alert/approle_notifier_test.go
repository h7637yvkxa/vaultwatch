package alert

import (
	"bytes"
	"strings"
	"testing"

	"github.com/subtlepseudonym/vaultwatch/internal/vault"
)

func makeAppRole(name, roleID string, ttl, maxTTL int) vault.AppRoleEntry {
	return vault.AppRoleEntry{
		Name:        name,
		RoleID:      roleID,
		TokenTTL:    ttl,
		TokenMaxTTL: maxTTL,
	}
}

func TestAppRoleNotifier_NoEntries(t *testing.T) {
	var buf bytes.Buffer
	n := NewAppRoleNotifier(&buf)
	if err := n.Notify(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "no AppRoles found") {
		t.Errorf("expected no-roles message, got: %s", buf.String())
	}
}

func TestAppRoleNotifier_WithEntries(t *testing.T) {
	var buf bytes.Buffer
	n := NewAppRoleNotifier(&buf)
	entries := []vault.AppRoleEntry{
		makeAppRole("web", "abc-123", 3600, 7200),
		makeAppRole("api", "", 0, 0),
	}
	if err := n.Notify(entries); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "2 AppRole(s)") {
		t.Errorf("expected count in output, got: %s", out)
	}
	if !strings.Contains(out, "web") {
		t.Errorf("expected 'web' in output, got: %s", out)
	}
	if !strings.Contains(out, "role_id=abc-123") {
		t.Errorf("expected role_id in output, got: %s", out)
	}
	if !strings.Contains(out, "token_ttl=3600s") {
		t.Errorf("expected token_ttl in output, got: %s", out)
	}
}

func TestAppRoleNotifier_NilWriter_UsesStdout(t *testing.T) {
	n := NewAppRoleNotifier(nil)
	if n.w == nil {
		t.Error("expected non-nil writer when nil passed")
	}
}
