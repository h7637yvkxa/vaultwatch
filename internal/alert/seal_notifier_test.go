package alert

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/yourusername/vaultwatch/internal/vault"
)

func makeSealStatus(sealed bool) *vault.SealStatus {
	return &vault.SealStatus{
		Sealed:      sealed,
		Initialized: true,
		Version:     "1.14.0",
		ClusterName: "vault-cluster",
		T:           3,
		N:           5,
		Progress:    1,
		CheckedAt:   time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
	}
}

func TestSealNotifier_Unsealed(t *testing.T) {
	var buf bytes.Buffer
	n := NewSealNotifier(&buf)
	err := n.Notify(makeSealStatus(false))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "unsealed") {
		t.Errorf("expected 'unsealed' in output, got: %s", buf.String())
	}
}

func TestSealNotifier_Sealed(t *testing.T) {
	var buf bytes.Buffer
	n := NewSealNotifier(&buf)
	err := n.Notify(makeSealStatus(true))
	if err == nil {
		t.Error("expected error for sealed vault")
	}
	if !strings.Contains(buf.String(), "SEALED") {
		t.Errorf("expected 'SEALED' in output, got: %s", buf.String())
	}
}

func TestSealNotifier_NilStatus(t *testing.T) {
	var buf bytes.Buffer
	n := NewSealNotifier(&buf)
	err := n.Notify(nil)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "no seal status") {
		t.Errorf("expected 'no seal status' in output, got: %s", buf.String())
	}
}

func TestSealNotifier_NilWriter_UsesStdout(t *testing.T) {
	n := NewSealNotifier(nil)
	if n.w == nil {
		t.Error("expected non-nil writer")
	}
}
