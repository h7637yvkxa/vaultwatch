package alert

import (
	"fmt"
	"io"
	"os"

	"github.com/arcanericky/vaultwatch/internal/vault"
)

// HAStateNotifier writes Vault HA state information to a writer.
type HAStateNotifier struct {
	w io.Writer
}

// NewHAStateNotifier creates a new HAStateNotifier. If w is nil, os.Stdout is used.
func NewHAStateNotifier(w io.Writer) *HAStateNotifier {
	if w == nil {
		w = os.Stdout
	}
	return &HAStateNotifier{w: w}
}

// Notify writes the HA state to the configured writer.
func (n *HAStateNotifier) Notify(state *vault.HAState) error {
	if state == nil {
		fmt.Fprintln(n.w, "[HA State] no data available")
		return nil
	}

	fmt.Fprintln(n.w, "[HA State]")
	fmt.Fprintf(n.w, "  HA Enabled:       %v\n", state.HAEnabled)
	fmt.Fprintf(n.w, "  Is Self Leader:   %v\n", state.IsSelf)
	fmt.Fprintf(n.w, "  Leader Address:   %s\n", coalesceStr(state.LeaderAddress, "(none)"))
	fmt.Fprintf(n.w, "  Leader Cluster:   %s\n", coalesceStr(state.LeaderCluster, "(none)"))
	fmt.Fprintf(n.w, "  Perf Standby:     %v\n", state.PerfStandby)
	if state.ActiveTime != "" {
		fmt.Fprintf(n.w, "  Active Since:     %s\n", state.ActiveTime)
	}
	return nil
}

func coalesceStr(s, fallback string) string {
	if s == "" {
		return fallback
	}
	return s
}
