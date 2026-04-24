package alert

import (
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/yourusername/vaultwatch/internal/vault"
)

// LeaseCountNotifier formats and writes lease count results.
type LeaseCountNotifier struct {
	w io.Writer
}

// NewLeaseCountNotifier creates a new LeaseCountNotifier.
// If w is nil, os.Stdout is used.
func NewLeaseCountNotifier(w io.Writer) *LeaseCountNotifier {
	if w == nil {
		w = os.Stdout
	}
	return &LeaseCountNotifier{w: w}
}

// Notify writes a summary of the lease count result to the writer.
func (n *LeaseCountNotifier) Notify(result *vault.LeaseCountResult) error {
	if result == nil {
		fmt.Fprintln(n.w, "[lease-count] no data available")
		return nil
	}

	fmt.Fprintf(n.w, "[lease-count] total leases: %d\n", result.Total)

	if len(result.ByMount) == 0 {
		fmt.Fprintln(n.w, "[lease-count] no per-mount breakdown available")
		return nil
	}

	mounts := make([]string, 0, len(result.ByMount))
	for m := range result.ByMount {
		mounts = append(mounts, m)
	}
	sort.Strings(mounts)

	for _, m := range mounts {
		fmt.Fprintf(n.w, "  mount=%-30s count=%d\n", m, result.ByMount[m])
	}
	return nil
}
