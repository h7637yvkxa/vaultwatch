package alert

import (
	"bytes"
	"strings"
	"testing"

	"github.com/wormhole-enterprise/vaultwatch/internal/vault"
)

func makeStepDownResult(success bool, msg string) *vault.StepDownResult {
	return &vault.StepDownResult{Success: success, Message: msg}
}

func TestStepDownNotifier_NilResult(t *testing.T) {
	var buf bytes.Buffer
	n := NewStepDownNotifier(&buf)
	if err := n.Notify(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "no result") {
		t.Errorf("expected 'no result' in output, got: %s", buf.String())
	}
}

func TestStepDownNotifier_Success(t *testing.T) {
	var buf bytes.Buffer
	n := NewStepDownNotifier(&buf)
	result := makeStepDownResult(true, "step-down accepted")
	if err := n.Notify(result); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "SUCCESS") {
		t.Errorf("expected SUCCESS in output, got: %s", out)
	}
	if !strings.Contains(out, "step-down accepted") {
		t.Errorf("expected message in output, got: %s", out)
	}
}

func TestStepDownNotifier_Failure(t *testing.T) {
	var buf bytes.Buffer
	n := NewStepDownNotifier(&buf)
	result := makeStepDownResult(false, "unexpected status: 403")
	if err := n.Notify(result); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "FAILED") {
		t.Errorf("expected FAILED in output, got: %s", out)
	}
}

func TestStepDownNotifier_NilWriter_UsesStdout(t *testing.T) {
	n := NewStepDownNotifier(nil)
	if n.w == nil {
		t.Error("expected writer to default to stdout, got nil")
	}
}
