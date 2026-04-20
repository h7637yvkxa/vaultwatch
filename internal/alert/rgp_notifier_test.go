package alert_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourusername/vaultwatch/internal/alert"
)

func makeRGPPolicy(name string) alert.RGPPolicy {
	return alert.RGPPolicy{Name: name}
}

func TestRGPNotifier_NoPolicies(t *testing.T) {
	var buf bytes.Buffer
	n := alert.NewRGPNotifier(&buf)
	if err := n.Notify(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No role governing policies") {
		t.Errorf("expected no-policies message, got: %s", buf.String())
	}
}

func TestRGPNotifier_WithPolicies(t *testing.T) {
	var buf bytes.Buffer
	n := alert.NewRGPNotifier(&buf)
	policies := []alert.RGPPolicy{
		makeRGPPolicy("admin-rgp"),
		makeRGPPolicy("readonly-rgp"),
	}
	if err := n.Notify(policies); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "2 role governing") {
		t.Errorf("expected count in output, got: %s", out)
	}
	if !strings.Contains(out, "admin-rgp") {
		t.Errorf("expected policy name in output, got: %s", out)
	}
}

func TestRGPNotifier_NilWriter_UsesStdout(t *testing.T) {
	n := alert.NewRGPNotifier(nil)
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}
