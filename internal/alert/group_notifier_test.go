package alert

import (
	"bytes"
	"strings"
	"testing"

	"github.com/wadefletch/vaultwatch/internal/vault"
)

func makeGroup(id, name, groupType string, policies []string) vault.GroupEntry {
	return vault.GroupEntry{
		ID:       id,
		Name:     name,
		Type:     groupType,
		Policies: policies,
	}
}

func TestGroupNotifier_NoGroups(t *testing.T) {
	var buf bytes.Buffer
	NewGroupNotifier(&buf).Notify(nil)
	if !strings.Contains(buf.String(), "no identity groups") {
		t.Errorf("expected no-groups message, got: %s", buf.String())
	}
}

func TestGroupNotifier_WithGroups(t *testing.T) {
	var buf bytes.Buffer
	groups := []vault.GroupEntry{
		makeGroup("id1", "admins", "internal", []string{"default", "admin"}),
		makeGroup("id2", "readers", "external", nil),
	}
	NewGroupNotifier(&buf).Notify(groups)
	out := buf.String()
	if !strings.Contains(out, "admins") {
		t.Errorf("expected admins in output, got: %s", out)
	}
	if !strings.Contains(out, "readers") {
		t.Errorf("expected readers in output, got: %s", out)
	}
	if !strings.Contains(out, "2 identity group") {
		t.Errorf("expected count in output, got: %s", out)
	}
}

func TestGroupNotifier_NoPolicies(t *testing.T) {
	var buf bytes.Buffer
	groups := []vault.GroupEntry{makeGroup("id1", "empty", "internal", nil)}
	NewGroupNotifier(&buf).Notify(groups)
	if !strings.Contains(buf.String(), "(none)") {
		t.Errorf("expected (none) for empty policies, got: %s", buf.String())
	}
}

func TestGroupNotifier_NilWriter_UsesStdout(t *testing.T) {
	n := NewGroupNotifier(nil)
	if n.w == nil {
		t.Error("expected non-nil writer")
	}
}
