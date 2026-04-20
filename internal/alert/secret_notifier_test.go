package alert_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/yourusername/vaultwatch/internal/alert"
	"github.com/yourusername/vaultwatch/internal/vault"
)

func makeSecretVersion(path string, ver int, destroyed bool) vault.SecretVersion {
	return vault.SecretVersion{
		Path:      path,
		Version:   ver,
		CreatedAt: time.Now(),
		Destroyed: destroyed,
	}
}

func TestSecretNotifier_NoVersions(t *testing.T) {
	var buf bytes.Buffer
	n := alert.NewSecretNotifier(&buf)
	if err := n.Notify(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "no versions found") {
		t.Errorf("expected 'no versions found', got: %s", buf.String())
	}
}

func TestSecretNotifier_WithVersions(t *testing.T) {
	var buf bytes.Buffer
	n := alert.NewSecretNotifier(&buf)
	versions := []vault.SecretVersion{
		makeSecretVersion("myapp/db", 2, false),
		makeSecretVersion("myapp/db", 1, true),
	}
	if err := n.Notify(versions); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "path=myapp/db") {
		t.Errorf("expected path in output, got: %s", out)
	}
	if !strings.Contains(out, "version=1") {
		t.Errorf("expected version=1 in output, got: %s", out)
	}
	if !strings.Contains(out, "destroyed") {
		t.Errorf("expected destroyed status in output, got: %s", out)
	}
}

func TestSecretNotifier_DeletedVersion(t *testing.T) {
	var buf bytes.Buffer
	n := alert.NewSecretNotifier(&buf)
	delTime := time.Now()
	versions := []vault.SecretVersion{
		{Path: "app/creds", Version: 3, CreatedAt: time.Now(), DeletedAt: &delTime, Destroyed: false},
	}
	if err := n.Notify(versions); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "deleted") {
		t.Errorf("expected 'deleted' status, got: %s", buf.String())
	}
}

func TestSecretNotifier_NilWriter_UsesStdout(t *testing.T) {
	n := alert.NewSecretNotifier(nil)
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}
