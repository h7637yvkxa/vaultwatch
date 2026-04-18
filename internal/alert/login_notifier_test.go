package alert

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/densestvoid/vaultwatch/internal/vault"
)

func makeLoginInfo(renewable bool, ttl int) *vault.LoginInfo {
	return &vault.LoginInfo{
		ClientToken:   "tok-test",
		Accessor:      "acc-test",
		Policies:      []string{"default"},
		LeaseDuration: ttl,
		Renewable:     renewable,
		IssuedAt:      time.Unix(1700000000, 0),
	}
}

func TestLoginNotifier_NilInfo(t *testing.T) {
	var buf bytes.Buffer
	n := NewLoginNotifier(&buf)
	if err := n.Notify(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "no token info") {
		t.Errorf("expected 'no token info', got: %s", buf.String())
	}
}

func TestLoginNotifier_WithInfo(t *testing.T) {
	var buf bytes.Buffer
	n := NewLoginNotifier(&buf)
	info := makeLoginInfo(true, 3600)
	if err := n.Notify(info); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	for _, want := range []string{"acc-test", "default", "3600", "renewable : yes"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output: %s", want, out)
		}
	}
}

func TestLoginNotifier_NonRenewable(t *testing.T) {
	var buf bytes.Buffer
	n := NewLoginNotifier(&buf)
	info := makeLoginInfo(false, 1800)
	_ = n.Notify(info)
	if !strings.Contains(buf.String(), "renewable : no") {
		t.Errorf("expected non-renewable output, got: %s", buf.String())
	}
}

func TestLoginNotifier_NilWriter_UsesStdout(t *testing.T) {
	n := NewLoginNotifier(nil)
	if n.w == nil {
		t.Error("expected fallback to stdout")
	}
}
