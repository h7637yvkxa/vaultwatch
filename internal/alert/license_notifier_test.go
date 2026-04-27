package alert

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/yourusername/vaultwatch/internal/vault"
)

func makeLicenseInfo(expiresIn time.Duration, terminated bool) *vault.LicenseInfo {
	return &vault.LicenseInfo{
		LicenseID:      "lic-001",
		CustomerName:   "Test Corp",
		ExpirationTime: time.Now().Add(expiresIn),
		Terminated:     terminated,
		Features:       []string{"Replication"},
	}
}

func TestLicenseNotifier_NilInfo(t *testing.T) {
	var buf bytes.Buffer
	n := NewLicenseNotifier(&buf)
	if err := n.Notify(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "no license information") {
		t.Errorf("expected no-info message, got: %s", buf.String())
	}
}

func TestLicenseNotifier_OK(t *testing.T) {
	var buf bytes.Buffer
	n := NewLicenseNotifier(&buf)
	info := makeLicenseInfo(60*24*time.Hour, false)
	_ = n.Notify(info)
	out := buf.String()
	if !strings.Contains(out, "status=OK") {
		t.Errorf("expected OK status, got: %s", out)
	}
}

func TestLicenseNotifier_Warning(t *testing.T) {
	var buf bytes.Buffer
	n := NewLicenseNotifier(&buf)
	info := makeLicenseInfo(15*24*time.Hour, false)
	_ = n.Notify(info)
	out := buf.String()
	if !strings.Contains(out, "status=WARNING") {
		t.Errorf("expected WARNING status, got: %s", out)
	}
}

func TestLicenseNotifier_Critical(t *testing.T) {
	var buf bytes.Buffer
	n := NewLicenseNotifier(&buf)
	info := makeLicenseInfo(3*24*time.Hour, false)
	_ = n.Notify(info)
	out := buf.String()
	if !strings.Contains(out, "status=CRITICAL") {
		t.Errorf("expected CRITICAL status, got: %s", out)
	}
}

func TestLicenseNotifier_Terminated(t *testing.T) {
	var buf bytes.Buffer
	n := NewLicenseNotifier(&buf)
	info := makeLicenseInfo(-1*time.Hour, true)
	_ = n.Notify(info)
	out := buf.String()
	if !strings.Contains(out, "status=TERMINATED") {
		t.Errorf("expected TERMINATED status, got: %s", out)
	}
}

func TestLicenseNotifier_NilWriter_UsesStdout(t *testing.T) {
	n := NewLicenseNotifier(nil)
	if n.w == nil {
		t.Error("expected non-nil writer when nil passed")
	}
}
