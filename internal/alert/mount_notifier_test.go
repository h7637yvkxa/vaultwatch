package alert_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/your-org/vaultwatch/internal/alert"
	"github.com/your-org/vaultwatch/internal/vault"
)

func makeMountEntry(path, typ, desc, accessor string) vault.MountEntry {
	return vault.MountEntry{
		Path:        path,
		Type:        typ,
		Description: desc,
		Accessor:    accessor,
	}
}

func TestMountNotifier_NoMounts(t *testing.T) {
	var buf bytes.Buffer
	n := alert.NewMountNotifier(&buf)
	if err := n.Notify(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "no secret engines") {
		t.Errorf("expected no-mounts message, got: %s", buf.String())
	}
}

func TestMountNotifier_WithMounts(t *testing.T) {
	var buf bytes.Buffer
	n := alert.NewMountNotifier(&buf)
	mounts := []vault.MountEntry{
		makeMountEntry("secret/", "kv", "key/value store", "kv_abc123"),
		makeMountEntry("pki/", "pki", "PKI engine", "pki_def456"),
	}
	if err := n.Notify(mounts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	for _, want := range []string{"secret/", "kv", "pki/", "PKI engine", "PATH", "TYPE"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output, got:\n%s", want, out)
		}
	}
}

func TestMountNotifier_NilWriter_UsesStdout(t *testing.T) {
	n := alert.NewMountNotifier(nil)
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
	// Just ensure it doesn't panic on notify with empty list
	_ = n.Notify([]vault.MountEntry{})
}
