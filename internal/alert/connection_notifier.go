package alert

import (
	"fmt"
	"io"
	"os"

	"github.com/yourusername/vaultwatch/internal/vault"
)

// ConnectionNotifier formats and writes Vault connectivity status.
type ConnectionNotifier struct {
	w io.Writer
}

// NewConnectionNotifier creates a ConnectionNotifier writing to w.
// If w is nil, os.Stdout is used.
func NewConnectionNotifier(w io.Writer) *ConnectionNotifier {
	if w == nil {
		w = os.Stdout
	}
	return &ConnectionNotifier{w: w}
}

// Notify prints the connection status to the configured writer.
func (n *ConnectionNotifier) Notify(status *vault.ConnectionStatus) error {
	if status == nil {
		fmt.Fprintln(n.w, "[connection] no status available")
		return nil
	}
	if !status.Reachable {
		fmt.Fprintf(n.w, "[connection] UNREACHABLE: %s\n", status.Error)
		return nil
	}
	state := "OK"
	if status.StatusCode >= 500 {
		state = "ERROR"
	} else if status.StatusCode != 200 {
		state = "DEGRADED"
	}
	fmt.Fprintf(n.w, "[connection] %s status=%d cluster=%q version=%s\n",
		state, status.StatusCode, status.ClusterName, status.Version)
	return nil
}
