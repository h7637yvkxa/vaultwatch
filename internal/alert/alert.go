package alert

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/yourusername/vaultwatch/internal/vault"
)

// Level represents the severity of an alert.
type Level string

const (
	LevelWarning  Level = "WARNING"
	LevelCritical Level = "CRITICAL"
)

// Alert holds information about an expiring lease.
type Alert struct {
	Level     Level
	LeaseID   string
	Path      string
	ExpiresAt time.Time
	TTL       time.Duration
}

// String returns a human-readable representation of the alert.
func (a Alert) String() string {
	return fmt.Sprintf(
		"[%s] lease=%s path=%s expires_in=%.0fs expires_at=%s",
		a.Level,
		a.LeaseID,
		a.Path,
		a.TTL.Seconds(),
		a.ExpiresAt.Format(time.RFC3339),
	)
}

// Notifier defines the interface for sending alerts.
type Notifier interface {
	Notify(alerts []Alert) error
}

// BuildAlerts converts expiring lease statuses into Alert values.
func BuildAlerts(statuses []vault.LeaseStatus) []Alert {
	alerts := make([]Alert, 0, len(statuses))
	for _, s := range statuses {
		if !s.IsExpiring {
			continue
		}
		lvl := LevelWarning
		if s.IsCritical {
			lvl = LevelCritical
		}
		alerts = append(alerts, Alert{
			Level:     lvl,
			LeaseID:   s.LeaseID,
			Path:      s.Path,
			ExpiresAt: s.ExpiresAt,
			TTL:       s.TTL,
		})
	}
	return alerts
}

// StdoutNotifier writes alerts to an io.Writer (defaults to os.Stdout).
type StdoutNotifier struct {
	Out io.Writer
}

// NewStdoutNotifier creates a StdoutNotifier writing to stdout.
func NewStdoutNotifier() *StdoutNotifier {
	return &StdoutNotifier{Out: os.Stdout}
}

// Notify prints each alert to the configured writer.
func (n *StdoutNotifier) Notify(alerts []Alert) error {
	for _, a := range alerts {
		fmt.Fprintln(n.Out, a.String())
	}
	return nil
}
