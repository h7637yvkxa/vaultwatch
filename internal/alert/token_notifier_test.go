package alert_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/yourusername/vaultwatch/internal/alert"
	"github.com/yourusername/vaultwatch/internal/vault"
)

func makeTokenInfo() *vault.TokenInfo {
	return &vault.TokenInfo{
		ID:          "tok-abc123",
		DisplayName: "ci-runner",
		Policies:    []string{"default", "read-only"},
		TTL:         2 * time.Hour,
		Renewable:   true,
	}
}

func TestTokenNotifier_NilInfo(t *testing.T) {
	var buf bytes.Buffer
	n := alert.NewTokenNotifier(&buf)
	if err := n.Notify(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "no token info") {
		t.Errorf("expected 'no token info' in output, got: %s", buf.String())
	}
}

func TestTokenNotifier_WithInfo(t *testing.T) {
	var buf bytes.Buffer
	n := alert.NewTokenNotifier(&buf)
	info := makeTokenInfo()
	if err := n.Notify(info); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	for _, want := range []string{"tok-abc123", "ci-runner", "2h", "default", "read-only"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output, got: %s", want, out)
		}
	}
}

func TestTokenNotifier_NonRenewable(t *testing.T) {
	var buf bytes.Buffer
	n := alert.NewTokenNotifier(&buf)
	info := makeTokenInfo()
	info.Renewable = false
	if err := n.Notify(info); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "non-renewable") {
		t.Errorf("expected 'non-renewable' in output for non-renewable token, got: %s", buf.String())
	}
}

func TestTokenNotifier_NilWriter_UsesStdout(t *testing.T) {
	n := alert.NewTokenNotifier(nil)
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}
