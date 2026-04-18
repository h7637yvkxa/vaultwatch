package alert_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourusername/vaultwatch/internal/alert"
	"github.com/yourusername/vaultwatch/internal/vault"
)

func makeQuotaRule(name, typ, path string) vault.QuotaRule {
	return vault.QuotaRule{
		Name:      name,
		Type:      typ,
		Path:      path,
		Rate:      100,
		Burst:     200,
		MaxLeases: 500,
	}
}

func TestQuotaNotifier_NoQuotas(t *testing.T) {
	var buf bytes.Buffer
	n := alert.NewQuotaNotifier(&buf)
	if err := n.Notify(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "no quota rules found") {
		t.Errorf("expected no-quota message, got: %s", buf.String())
	}
}

func TestQuotaNotifier_RateLimit(t *testing.T) {
	var buf bytes.Buffer
	n := alert.NewQuotaNotifier(&buf)
	rules := []vault.QuotaRule{makeQuotaRule("rl-global", "rate-limit", "secret/")}
	if err := n.Notify(rules); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "rate-limit") {
		t.Errorf("expected rate-limit in output, got: %s", out)
	}
	if !strings.Contains(out, "rl-global") {
		t.Errorf("expected quota name in output, got: %s", out)
	}
}

func TestQuotaNotifier_LeaseCount(t *testing.T) {
	var buf bytes.Buffer
	n := alert.NewQuotaNotifier(&buf)
	rules := []vault.QuotaRule{makeQuotaRule("lc-root", "lease-count", "")}
	if err := n.Notify(rules); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "max=") {
		t.Errorf("expected max= in lease-count output, got: %s", out)
	}
}

func TestQuotaNotifier_NilWriter_UsesStdout(t *testing.T) {
	n := alert.NewQuotaNotifier(nil)
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}
