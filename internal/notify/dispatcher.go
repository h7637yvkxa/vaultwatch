package notify

import (
	"context"
	"fmt"

	"github.com/yourusername/vaultwatch/internal/alert"
	"github.com/yourusername/vaultwatch/internal/config"
)

// Notifier is the interface that wraps the Notify method.
type Notifier interface {
	Notify(ctx context.Context, alerts []alert.Alert) error
}

// Dispatcher routes alerts to one or more configured notifiers.
type Dispatcher struct {
	notifiers []Notifier
}

// NewDispatcher builds a Dispatcher from the application config,
// wiring up every enabled notification channel.
func NewDispatcher(cfg *config.Config) (*Dispatcher, error) {
	var notifiers []Notifier

	// Stdout is always enabled.
	notifiers = append(notifiers, alert.NewStdoutNotifier())

	if cfg.Slack.WebhookURL != "" {
		notifiers = append(notifiers, alert.NewSlackNotifier(cfg.Slack.WebhookURL))
	}

	if len(notifiers) == 0 {
		return nil, fmt.Errorf("no notifiers configured")
	}

	return &Dispatcher{notifiers: notifiers}, nil
}

// Dispatch sends alerts through every registered notifier.
// Errors are collected and returned as a combined error.
func (d *Dispatcher) Dispatch(ctx context.Context, alerts []alert.Alert) error {
	if len(alerts) == 0 {
		return nil
	}

	var errs []error
	for _, n := range d.notifiers {
		if err := n.Notify(ctx, alerts); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("dispatch errors: %v", errs)
	}
	return nil
}
