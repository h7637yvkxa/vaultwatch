package config

import (
	"errors"
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds the full vaultwatch configuration.
type Config struct {
	Vault   VaultConfig   `yaml:"vault"`
	Alerts  AlertsConfig  `yaml:"alerts"`
	Notify  NotifyConfig  `yaml:"notify"`
}

// VaultConfig contains Vault connection settings.
type VaultConfig struct {
	Address string `yaml:"address"`
	Token   string `yaml:"token"`
}

// AlertsConfig defines thresholds for warning and critical alerts.
type AlertsConfig struct {
	WarnBefore     time.Duration `yaml:"warn_before"`
	CriticalBefore time.Duration `yaml:"critical_before"`
	PollInterval   time.Duration `yaml:"poll_interval"`
}

// NotifyConfig holds notifier-specific settings.
type NotifyConfig struct {
	SlackWebhook string `yaml:"slack_webhook"`
	Stdout       bool   `yaml:"stdout"`
}

// Load reads and validates a Config from the given YAML file path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("config: read file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("config: parse yaml: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("config: validation: %w", err)
	}

	return &cfg, nil
}

func (c *Config) validate() error {
	if c.Vault.Address == "" {
		return errors.New("vault.address is required")
	}
	if c.Alerts.WarnBefore <= 0 {
		return errors.New("alerts.warn_before must be a positive duration")
	}
	if c.Alerts.CriticalBefore <= 0 {
		return errors.New("alerts.critical_before must be a positive duration")
	}
	if c.Alerts.CriticalBefore >= c.Alerts.WarnBefore {
		return errors.New("alerts.critical_before must be less than alerts.warn_before")
	}
	if c.Alerts.PollInterval <= 0 {
		c.Alerts.PollInterval = 60 * time.Second
	}
	return nil
}
