package alert_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/yourusername/vaultwatch/internal/alert"
	"github.com/yourusername/vaultwatch/internal/vault"
)

func makeStatus(leaseID, path string, ttl time.Duration, isExpiring, isCritical bool) vault.LeaseStatus {
	return vault.LeaseStatus{
		LeaseID:    leaseID,
		Path:       path,
		TTL:        ttl,
		ExpiresAt:  time.Now().Add(ttl),
		IsExpiring: isExpiring,
		IsCritical: isCritical,
	}
}

func TestBuildAlerts_SkipsNonExpiring(t *testing.T) {
	statuses := []vault.LeaseStatus{
		makeStatus("id1", "secret/a", 2*time.Hour, false, false),
	}
	alerts := alert.BuildAlerts(statuses)
	if len(alerts) != 0 {
		t.Fatalf("expected 0 alerts, got %d", len(alerts))
	}
}

func TestBuildAlerts_WarningLevel(t *testing.T) {
	statuses := []vault.LeaseStatus{
		makeStatus("id2", "secret/b", 30*time.Minute, true, false),
	}
	alerts := alert.BuildAlerts(statuses)
	if len(alerts) != 1 {
		t.Fatalf("expected 1 alert, got %d", len(alerts))
	}
	if alerts[0].Level != alert.LevelWarning {
		t.Errorf("expected WARNING, got %s", alerts[0].Level)
	}
}

func TestBuildAlerts_CriticalLevel(t *testing.T) {
	statuses := []vault.LeaseStatus{
		makeStatus("id3", "secret/c", 5*time.Minute, true, true),
	}
	alerts := alert.BuildAlerts(statuses)
	if len(alerts) != 1 {
		t.Fatalf("expected 1 alert, got %d", len(alerts))
	}
	if alerts[0].Level != alert.LevelCritical {
		t.Errorf("expected CRITICAL, got %s", alerts[0].Level)
	}
}

func TestStdoutNotifier_Notify(t *testing.T) {
	var buf bytes.Buffer
	n := &alert.StdoutNotifier{Out: &buf}

	alerts := []alert.Alert{
		{
			Level:     alert.LevelWarning,
			LeaseID:   "lease-abc",
			Path:      "secret/myapp",
			ExpiresAt: time.Now().Add(20 * time.Minute),
			TTL:       20 * time.Minute,
		},
	}

	if err := n.Notify(alerts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "WARNING") {
		t.Errorf("expected WARNING in output, got: %s", out)
	}
	if !strings.Contains(out, "lease-abc") {
		t.Errorf("expected lease-abc in output, got: %s", out)
	}
}

func TestAlert_String(t *testing.T) {
	a := alert.Alert{
		Level:     alert.LevelCritical,
		LeaseID:   "lease-xyz",
		Path:      "database/creds/my-role",
		ExpiresAt: time.Now().Add(3 * time.Minute),
		TTL:       3 * time.Minute,
	}
	s := a.String()
	if !strings.HasPrefix(s, "[CRITICAL]") {
		t.Errorf("expected string to start with [CRITICAL], got: %s", s)
	}
}
