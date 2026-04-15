package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "vaultwatch-*.yaml")
	if err != nil {
		t.Fatalf("creating temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("writing temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

const validYAML = `
vault:
  address: "http://127.0.0.1:8200"
  token: "root"
monitor:
  interval: "30s"
  paths:
    - "secret/myapp"
alerts:
  warn_before:
    - "72h"
    - "24h"
  slack_webhook: "https://hooks.slack.com/xxx"
`

func TestLoad_Valid(t *testing.T) {
	path := writeTemp(t, validYAML)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Vault.Address != "http://127.0.0.1:8200" {
		t.Errorf("expected vault address, got %q", cfg.Vault.Address)
	}
	durations, _ := cfg.Alerts.WarnDurations()
	if len(durations) != 2 || durations[0] != 72*time.Hour {
		t.Errorf("unexpected warn durations: %v", durations)
	}
}

func TestLoad_MissingAddress(t *testing.T) {
	y := `vault:\n  token: "root"\nmonitor:\n  paths: ["secret/x"]\n`
	path := writeTemp(t, "vault:\n  token: root\nmonitor:\n  paths:\n    - secret/x\n")
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for missing vault address")
	}
	_ = y
}

func TestLoad_InvalidWarnBefore(t *testing.T) {
	y := `vault:\n  address: http://x\n  token: t\nmonitor:\n  paths: [x]\nalerts:\n  warn_before: ["notaduration"]\n`
	path := writeTemp(t, y)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for invalid warn_before")
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := Load(filepath.Join(t.TempDir(), "nonexistent.yaml"))
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
