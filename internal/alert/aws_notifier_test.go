package alert_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/rbhz/vaultwatch/internal/alert"
	"github.com/rbhz/vaultwatch/internal/vault"
)

func makeAWSRole(name, authType string) vault.AWSRole {
	return vault.AWSRole{Name: name, AuthType: authType}
}

func TestAWSNotifier_NoRoles(t *testing.T) {
	var buf bytes.Buffer
	n := alert.NewAWSNotifier(&buf)
	if err := n.Notify(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No AWS auth roles") {
		t.Errorf("expected no-roles message, got: %s", buf.String())
	}
}

func TestAWSNotifier_WithRoles(t *testing.T) {
	var buf bytes.Buffer
	n := alert.NewAWSNotifier(&buf)
	roles := []vault.AWSRole{
		makeAWSRole("dev-role", "iam"),
		makeAWSRole("prod-role", "ec2"),
	}
	if err := n.Notify(roles); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "dev-role") {
		t.Errorf("expected dev-role in output")
	}
	if !strings.Contains(out, "auth_type=iam") {
		t.Errorf("expected auth_type=iam in output")
	}
	if !strings.Contains(out, "2 role(s)") {
		t.Errorf("expected count in output")
	}
}

func TestAWSNotifier_NilWriter_UsesStdout(t *testing.T) {
	n := alert.NewAWSNotifier(nil)
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
	// Should not panic
	_ = n.Notify(nil)
}
