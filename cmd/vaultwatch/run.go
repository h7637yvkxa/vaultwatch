package main

import (
	"context"
	"fmt"
	"log"

	"github.com/yourusername/vaultwatch/internal/alert"
	"github.com/yourusername/vaultwatch/internal/config"
	"github.com/yourusername/vaultwatch/internal/notify"
	"github.com/yourusername/vaultwatch/internal/vault"
)

// run is the main application entry-point, separated from main() for
// testability.
func run(ctx context.Context, cfgPath string) error {
	cfg, err := config.Load(cfgPath)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	client, err := vault.NewClient(cfg)
	if err != nil {
		return fmt.Errorf("create vault client: %w", err)
	}

	leases, err := client.ListLeases(ctx)
	if err != nil {
		return fmt.Errorf("list leases: %w", err)
	}

	log.Printf("found %d leases", len(leases))

	statuses := vault.CheckExpiry(leases, cfg.WarnBefore, cfg.CriticalBefore)
	expiring := vault.FilterExpiring(statuses)
	alerts := alert.BuildAlerts(expiring)

	dispatcher, err := notify.NewDispatcher(cfg)
	if err != nil {
		return fmt.Errorf("build dispatcher: %w", err)
	}

	if err := dispatcher.Dispatch(ctx, alerts); err != nil {
		return fmt.Errorf("dispatch alerts: %w", err)
	}

	log.Printf("dispatched %d alert(s)", len(alerts))
	return nil
}
