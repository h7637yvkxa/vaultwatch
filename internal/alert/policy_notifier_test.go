package alert

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

func TestPolicyNotifier_NoPolicies(t *testing.T) {
	var buf bytes.Buffer
	n := NewPolicyNotifier(&buf)
	if err := n.Notify(PolicyReport{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "no policies found") {
		t.Errorf("expected 'no policies found', got: %s", buf.String())
	}
}

func TestPolicyNotifier_WithPolicies(t *testing.T) {
	var buf bytes.Buffer
	n := NewPolicyNotifier(&buf)
	report := PolicyReport{Policies: []string{"default", "root", "app-policy"}}
	if err := n.Notify(report); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "3 policies") {
		t.Errorf("expected count in output, got: %s", out)
	}
	if !strings.Contains(out, "app-policy") {
		t.Errorf("expected app-policy in output, got: %s", out)
	}
}

func TestPolicyNotifier_WithError(t *testing.T) {
	var buf bytes.Buffer
	n := NewPolicyNotifier(&buf)
	report := PolicyReport{Error: errors.New("permission denied")}
	if err := n.Notify(report); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "permission denied") {
		t.Errorf("expected error in output, got: %s", buf.String())
	}
}

func TestPolicyNotifier_NilWriter_UsesStdout(t *testing.T) {
	n := NewPolicyNotifier(nil)
	if n.w == nil {
		t.Error("expected non-nil writer")
	}
}

func TestPolicyNotifier_Summary(t *testing.T) {
	n := NewPolicyNotifier(nil)
	s := n.Summary(PolicyReport{Policies: []string{"a", "b"}})
	if !strings.Contains(s, "2 policies") {
		t.Errorf("unexpected summary: %s", s)
	}
}
