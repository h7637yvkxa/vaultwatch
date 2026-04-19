package alert_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/arcanericky/vaultwatch/internal/alert"
	"github.com/arcanericky/vaultwatch/internal/vault"
)

func makeSSHRole(name, keyType string) vault.SSHRole {
	return vault.SSHRole{Name: name, KeyType: keyType, TTL: "1h", MaxTTL: "24h", AllowedUsers: "ubuntu"}
}

func TestSSHNotifier_NoRoles(t *testing.T) {
	var buf bytes.Buffer
	n := alert.NewSSHNotifier(&buf)
	if err := n.Notify(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No roles") {
		t.Errorf("expected 'No roles' message, got: %s", buf.String())
	}
}

func TestSSHNotifier_WithRoles(t *testing.T) {
	var buf bytes.Buffer
	n := alert.NewSSHNotifier(&buf)
	roles := []vault.SSHRole{
		makeSSHRole("prod-role", "ca"),
		makeSSHRole("dev-role", "otp"),
	}
	if err := n.Notify(roles); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "prod-role") {
		t.Errorf("expected prod-role in output")
	}
	if !strings.Contains(out, "dev-role") {
		t.Errorf("expected dev-role in output")
	}
	if !strings.Contains(out, "2 role(s)") {
		t.Errorf("expected role count in output, got: %s", out)
	}
}

func TestSSHNotifier_NilWriter_UsesStdout(t *testing.T) {
	n := alert.NewSSHNotifier(nil)
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}
