package alert_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/densestvoid/vaultwatch/internal/alert"
	"github.com/densestvoid/vaultwatch/internal/vault"
)

func makeCertEntry(name string) vault.CertStoreEntry {
	return vault.CertStoreEntry{Name: name}
}

func TestCertStoreNotifier_NoEntries(t *testing.T) {
	var buf bytes.Buffer
	n := alert.NewCertStoreNotifier(&buf, "cert")
	if err := n.Notify(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "no certificate entries") {
		t.Errorf("expected no-entries message, got: %s", buf.String())
	}
}

func TestCertStoreNotifier_WithEntries(t *testing.T) {
	var buf bytes.Buffer
	n := alert.NewCertStoreNotifier(&buf, "cert")
	entries := []vault.CertStoreEntry{
		makeCertEntry("web-cert"),
		makeCertEntry("api-cert"),
	}
	if err := n.Notify(entries); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "web-cert") {
		t.Errorf("expected web-cert in output, got: %s", out)
	}
	if !strings.Contains(out, "2 certificate") {
		t.Errorf("expected count in output, got: %s", out)
	}
}

func TestCertStoreNotifier_NilWriter_UsesStdout(t *testing.T) {
	n := alert.NewCertStoreNotifier(nil, "cert")
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}
