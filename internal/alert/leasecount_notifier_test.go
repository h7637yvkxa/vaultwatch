package alert

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourusername/vaultwatch/internal/vault"
)

func makeLeaseCountResult(total int, byMount map[string]int) *vault.LeaseCountResult {
	return &vault.LeaseCountResult{
		Total:   total,
		ByMount: byMount,
	}
}

func TestLeaseCountNotifier_NilResult(t *testing.T) {
	var buf bytes.Buffer
	n := NewLeaseCountNotifier(&buf)
	_ = n.Notify(nil)
	if !strings.Contains(buf.String(), "no data") {
		t.Errorf("expected 'no data' message, got: %s", buf.String())
	}
}

func TestLeaseCountNotifier_WithCounts(t *testing.T) {
	var buf bytes.Buffer
	n := NewLeaseCountNotifier(&buf)
	result := makeLeaseCountResult(55, map[string]int{
		"secret/": 40,
		"aws/":    15,
	})
	_ = n.Notify(result)
	out := buf.String()
	if !strings.Contains(out, "total leases: 55") {
		t.Errorf("expected total leases in output, got: %s", out)
	}
	if !strings.Contains(out, "secret/") {
		t.Errorf("expected secret/ mount in output, got: %s", out)
	}
	if !strings.Contains(out, "aws/") {
		t.Errorf("expected aws/ mount in output, got: %s", out)
	}
}

func TestLeaseCountNotifier_EmptyMounts(t *testing.T) {
	var buf bytes.Buffer
	n := NewLeaseCountNotifier(&buf)
	result := makeLeaseCountResult(0, map[string]int{})
	_ = n.Notify(result)
	out := buf.String()
	if !strings.Contains(out, "total leases: 0") {
		t.Errorf("expected zero total, got: %s", out)
	}
	if !strings.Contains(out, "no per-mount breakdown") {
		t.Errorf("expected no per-mount message, got: %s", out)
	}
}

func TestLeaseCountNotifier_NilWriter_UsesStdout(t *testing.T) {
	n := NewLeaseCountNotifier(nil)
	if n.w == nil {
		t.Error("expected non-nil writer when nil passed")
	}
}

func TestLeaseCountNotifier_SortedOutput(t *testing.T) {
	var buf bytes.Buffer
	n := NewLeaseCountNotifier(&buf)
	result := makeLeaseCountResult(10, map[string]int{
		"zzz/": 1,
		"aaa/": 9,
	})
	_ = n.Notify(result)
	out := buf.String()
	aIdx := strings.Index(out, "aaa/")
	zIdx := strings.Index(out, "zzz/")
	if aIdx > zIdx {
		t.Error("expected sorted output: aaa/ before zzz/")
	}
}
