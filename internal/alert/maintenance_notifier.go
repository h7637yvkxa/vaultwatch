package alert

import (
	"fmt"
	"io"
	"os"

	"github.com/wiliamsouza/vaultwatch/internal/vault"
)

// MaintenanceNotifier formats and writes Vault maintenance status output.
type MaintenanceNotifier struct {
	w io.Writer
}

// NewMaintenanceNotifier creates a MaintenanceNotifier that writes to w.
// If w is nil, os.Stdout is used.
func NewMaintenanceNotifier(w io.Writer) *MaintenanceNotifier {
	if w == nil {
		w = os.Stdout
	}
	return &MaintenanceNotifier{w: w}
}

// Notify writes the maintenance status to the configured writer.
func (n *MaintenanceNotifier) Notify(status *vault.MaintenanceStatus) error {
	if status == nil {
		fmt.Fprintln(n.w, "[maintenance] no status available")
		return nil
	}

	if !status.Enabled {
		fmt.Fprintln(n.w, "[maintenance] status: disabled — Vault is operating normally")
		return nil
	}

	msg := status.Message
	if msg == "" {
		msg = "(no message provided)"
	}
	fmt.Fprintf(n.w, "[maintenance] status: ENABLED — %s\n", msg)
	return nil
}
