package alert_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/eliziario/vaultwatch/internal/alert"
	"github.com/eliziario/vaultwatch/internal/vault"
)

func makeGitHubTeam(name, policy string) vault.GitHubTeam {
	return vault.GitHubTeam{Name: name, Policy: policy}
}

func TestGitHubNotifier_NoTeams(t *testing.T) {
	var buf bytes.Buffer
	n := alert.NewGitHubNotifier(&buf)
	if err := n.Notify(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "no teams") {
		t.Errorf("expected 'no teams' in output, got: %s", buf.String())
	}
}

func TestGitHubNotifier_WithTeams(t *testing.T) {
	var buf bytes.Buffer
	n := alert.NewGitHubNotifier(&buf)
	teams := []vault.GitHubTeam{
		makeGitHubTeam("dev", "dev-policy"),
		makeGitHubTeam("ops", "ops-policy"),
	}
	if err := n.Notify(teams); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "dev") {
		t.Errorf("expected 'dev' in output")
	}
	if !strings.Contains(out, "2 team") {
		t.Errorf("expected count in output, got: %s", out)
	}
}

func TestGitHubNotifier_NilWriter_UsesStdout(t *testing.T) {
	n := alert.NewGitHubNotifier(nil)
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestGitHubNotifier_SingleTeam(t *testing.T) {
	var buf bytes.Buffer
	n := alert.NewGitHubNotifier(&buf)
	teams := []vault.GitHubTeam{
		makeGitHubTeam("admin", "admin-policy"),
	}
	if err := n.Notify(teams); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "admin") {
		t.Errorf("expected 'admin' in output, got: %s", out)
	}
	if !strings.Contains(out, "1 team") {
		t.Errorf("expected '1 team' in output, got: %s", out)
	}
}
