package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds the top-level vaultwatch configuration.
type Config struct {
	Vault   VaultConfig   `yaml:"vault"`
	Alerts  AlertsConfig  `yaml:"alerts"`
	Monitor MonitorConfig `yaml:"monitor"`
}

// VaultConfig contains connection details for HashiCorp Vault.
type VaultConfig struct {
	Address string `yaml:"address"`
	Token   string `yaml:"token"`
	CACert  string `yaml:"ca_cert"`
}

// AlertsConfig defines how and when alerts are sent.
type AlertsConfig struct {
	WarnBefore  []string `yaml:"warn_before"`
	SlackWebhook string  `yaml:"slack_webhook"`
	Email        string  `yaml:"email"`
}

// MonitorConfig controls polling behaviour.
type MonitorConfig struct {
	Interval string   `yaml:"interval"`
	Paths    []string `yaml:"paths"`
}

// WarnDurations parses WarnBefore strings into time.Duration values.
func (a AlertsConfig) WarnDurations() ([]time.Duration, error) {
	var durations []time.Duration
	for _, s := range a.WarnBefore {
		d, err := time.ParseDuration(s)
		if err != nil {
			return nil, fmt.Errorf("invalid warn_before value %q: %w", s, err)
		}
		durations = append(durations, d)
	}
	return durations, nil
}

// Load reads and parses a YAML config file at the given path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config file: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (c *Config) validate() error {
	if c.Vault.Address == "" {
		return fmt.Errorf("vault.address is required")
	}
	if c.Vault.Token == "" {
		return fmt.Errorf("vault.token is required")
	}
	if len(c.Monitor.Paths) == 0 {
		return fmt.Errorf("monitor.paths must contain at least one path")
	}
	if _, err := c.Alerts.WarnDurations(); err != nil {
		return err
	}
	return nil
}
