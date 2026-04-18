package alert_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/your-org/vaultwatch/internal/alert"
	"github.com/your-org/vaultwatch/internal/vault"
)

func makePlugin(name, typ string) vault.PluginInfo {
	return vault.PluginInfo{Name: name, Type: typ}
}

func TestPluginNotifier_NoPlugins(t *testing.T) {
	var buf bytes.Buffer
	n := alert.NewPluginNotifier(&buf)
	if err := n.Notify(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "no plugins") {
		t.Errorf("expected 'no plugins' message, got: %s", buf.String())
	}
}

func TestPluginNotifier_WithPlugins(t *testing.T) {
	var buf bytes.Buffer
	n := alert.NewPluginNotifier(&buf)
	plugins := []vault.PluginInfo{
		makePlugin("aws", "secret"),
		makePlugin("gcp", "secret"),
		makePlugin("userpass", "auth"),
	}
	if err := n.Notify(plugins); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "3 plugin(s)") {
		t.Errorf("expected count in output, got: %s", out)
	}
	if !strings.Contains(out, "secret") {
		t.Errorf("expected 'secret' type in output, got: %s", out)
	}
	if !strings.Contains(out, "auth") {
		t.Errorf("expected 'auth' type in output, got: %s", out)
	}
}

func TestPluginNotifier_NilWriter_UsesStdout(t *testing.T) {
	n := alert.NewPluginNotifier(nil)
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
	// Just ensure it doesn't panic on a small call
	_ = n.Notify([]vault.PluginInfo{})
}
