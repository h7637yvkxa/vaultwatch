package alert_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/arussellsaw/vaultwatch/internal/alert"
	"github.com/arussellsaw/vaultwatch/internal/vault"
)

func makeGCPRole(name, secretType, project string) vault.GCPRole {
	return vault.GCPRole{Name: name, SecretType: secretType, Project: project}
}

func TestGCPNotifier_NoRoles(t *testing.T) {
	var buf bytes.Buffer
	n := alert.NewGCPNotifier(&buf)
	if err := n.Notify(nil); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "no roles found") {
		t.Errorf("unexpected output: %s", buf.String())
	}
}

func TestGCPNotifier_WithRoles(t *testing.T) {
	var buf bytes.Buffer
	n := alert.NewGCPNotifier(&buf)
	roles := []vault.GCPRole{
		makeGCPRole("prod-role", "access_token", "my-project"),
		makeGCPRole("dev-role", "", ""),
	}
	if err := n.Notify(roles); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "prod-role") {
		t.Errorf("expected prod-role in output: %s", out)
	}
	if !strings.Contains(out, "my-project") {
		t.Errorf("expected project in output: %s", out)
	}
	if !strings.Contains(out, "access_token") {
		t.Errorf("expected secret type in output: %s", out)
	}
	if !strings.Contains(out, "dev-role") {
		t.Errorf("expected dev-role in output: %s", out)
	}
}

func TestGCPNotifier_NilWriter_UsesStdout(t *testing.T) {
	n := alert.NewGCPNotifier(nil)
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}
