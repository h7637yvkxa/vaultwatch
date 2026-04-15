package alert

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/your-org/vaultwatch/internal/vault"
)

// RenewNotifier reports renewal results to a writer (defaults to os.Stdout).
type RenewNotifier struct {
	out io.Writer
}

// NewRenewNotifier creates a RenewNotifier that writes to w.
// If w is nil, os.Stdout is used.
func NewRenewNotifier(w io.Writer) *RenewNotifier {
	if w == nil {
		w = os.Stdout
	}
	return &RenewNotifier{out: w}
}

// Notify prints a summary of renewal results.
func (n *RenewNotifier) Notify(_ context.Context, results []vault.RenewResult) error {
	if len(results) == 0 {
		fmt.Fprintln(n.out, "[renew] no leases required renewal")
		return nil
	}

	var errs []string
	for _, r := range results {
		if r.Error != nil {
			fmt.Fprintf(n.out, "[renew] FAILED  lease=%s error=%v\n", r.LeaseID, r.Error)
			errs = append(errs, r.Error.Error())
			continue
		}
		fmt.Fprintf(n.out, "[renew] OK      lease=%s new_ttl=%s\n", r.LeaseID, r.NewTTL)
	}

	if len(errs) > 0 {
		return fmt.Errorf("renewal errors: %s", strings.Join(errs, "; "))
	}
	return nil
}
