package alert_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourusername/vaultwatch/internal/alert"
)

func makeEGPPolicy(name string, paths []string) alert.EGPPolicy {
	return alert.EGPPolicy{Name: name, Paths: paths}
}

func TestEGPNotifier_NoPolicies(t *testing.T) {
	var buf bytes.Buffer
	n := alert.NewEGPNotifier(&buf)
	if err := n.Notify(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No endpoint governing policies") {
		t.Errorf("expected no-policies message, got: %s", buf.String())
	}
}

func TestEGPNotifier_WithPolicies(t *testing.T) {
	var buf bytes.Buffer
	n := alert.NewEGPNotifier(&buf)
	policies := []alert.EGPPolicy{
		makeEGPPolicy("allow-read", []string{"secret/*", "kv/*"}),
		makeEGPPolicy("deny-admin", []string{"sys/"}),
	}
	if err := n.Notify(policies); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "2 endpoint governing") {
		t.Errorf("expected count in output, got: %s", out)
	}
	if !strings.Contains(out, "allow-read") {
		t.Errorf("expected policy name in output, got: %s", out)
	}
	if !strings.Contains(out, "secret/*") {
		t.Errorf("expected path in output, got: %s", out)
	}
}

func TestEGPNotifier_NoPaths(t *testing.T) {
	var buf bytes.Buffer
	n := alert.NewEGPNotifier(&buf)
	policies := []alert.EGPPolicy{
		makeEGPPolicy("empty-policy", nil),
	}
	if err := n.Notify(policies); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "no paths") {
		t.Errorf("expected 'no paths' in output, got: %s", out)
	}
}

func TestEGPNotifier_NilWriter_UsesStdout(t *testing.T) {
	n := alert.NewEGPNotifier(nil)
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}
