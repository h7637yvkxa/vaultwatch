package alert

import (
	"fmt"
	"io"
	"os"

	"github.com/arjunsriva/vaultwatch/internal/vault"
)

// ReplicationNotifier prints replication status to a writer.
type ReplicationNotifier struct {
	w io.Writer
}

// NewReplicationNotifier creates a ReplicationNotifier. If w is nil, os.Stdout is used.
func NewReplicationNotifier(w io.Writer) *ReplicationNotifier {
	if w == nil {
		w = os.Stdout
	}
	return &ReplicationNotifier{w: w}
}

// Notify prints the replication status summary.
func (n *ReplicationNotifier) Notify(status *vault.ReplicationStatus) error {
	if status == nil {
		fmt.Fprintln(n.w, "[replication] no status available")
		return nil
	}

	drRole := role(status.DRPrimary)
	perfRole := role(status.PerfPrimary)

	fmt.Fprintf(n.w, "[replication] DR: mode=%s role=%s | Performance: mode=%s role=%s\n",
		coalesce(status.DRMode, "disabled"),
		drRole,
		coalesce(status.PerformanceMode, "disabled"),
		perfRole,
	)
	return nil
}

func role(primary bool) string {
	if primary {
		return "primary"
	}
	return "secondary"
}

func coalesce(s, fallback string) string {
	if s == "" {
		return fallback
	}
	return s
}
