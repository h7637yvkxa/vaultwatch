package alert_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourusername/vaultwatch/internal/alert"
	"github.com/yourusername/vaultwatch/internal/vault"
)

func makeReplicationStatus(mode, state string, primary bool) *vault.ReplicationStatus {
	return &vault.ReplicationStatus{
		Mode:      mode,
		State:     state,
		Primary:   primary,
		ClusterID: "cluster-abc",
	}
}

func TestReplicationNotifier_NilStatus(t *testing.T) {
	var buf bytes.Buffer
	n := alert.NewReplicationNotifier(nil, &buf)
	if err := n.Notify(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "unavailable") {
		t.Errorf("expected 'unavailable' in output, got: %s", buf.String())
	}
}

func TestReplicationNotifier_PrimaryActive(t *testing.T) {
	var buf bytes.Buffer
	status := makeReplicationStatus("performance", "running", true)
	n := alert.NewReplicationNotifier(nil, &buf)
	if err := n.Notify(status); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "primary") {
		t.Errorf("expected 'primary' in output, got: %s", out)
	}
	if !strings.Contains(out, "running") {
		t.Errorf("expected 'running' in output, got: %s", out)
	}
}

func TestReplicationNotifier_SecondaryDegraded(t *testing.T) {
	var buf bytes.Buffer
	status := makeReplicationStatus("dr", "degraded", false)
	n := alert.NewReplicationNotifier(nil, &buf)
	if err := n.Notify(status); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "secondary") {
		t.Errorf("expected 'secondary' in output, got: %s", out)
	}
	if !strings.Contains(out, "degraded") {
		t.Errorf("expected 'degraded' in output, got: %s", out)
	}
}

func TestReplicationNotifier_NilWriter_UsesStdout(t *testing.T) {
	status := makeReplicationStatus("performance", "running", true)
	n := alert.NewReplicationNotifier(nil, nil)
	if err := n.Notify(status); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
