package alert

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/drengskapur/vaultwatch/internal/vault"
)

// RotateNotifier formats rotation results and writes them to a writer.
type RotateNotifier struct {
	w io.Writer
}

// NewRotateNotifier returns a RotateNotifier that writes to w.
// If w is nil, os.Stdout is used.
func NewRotateNotifier(w io.Writer) *RotateNotifier {
	if w == nil {
		w = os.Stdout
	}
	return &RotateNotifier{w: w}
}

// Notify prints a summary of rotation results to the configured writer.
func (n *RotateNotifier) Notify(results []vault.RotateResult) error {
	if len(results) == 0 {
		_, err := fmt.Fprintln(n.w, "[rotate] no leases rotated")
		return err
	}

	var errs []error
	for _, r := range results {
		if r.Err != nil {
			_, err := fmt.Fprintf(n.w, "[rotate] FAILED lease=%s error=%v\n", r.LeaseID, r.Err)
			if err != nil {
				errs = append(errs, err)
			}
			continue
		}
		_, err := fmt.Fprintf(n.w,
			"[rotate] OK old=%s new=%s at=%s\n",
			r.LeaseID,
			r.NewLeaseID,
			r.RenewedAt.Format(time.RFC3339),
		)
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("rotate notifier encountered %d write error(s): first=%w", len(errs), errs[0])
	}
	return nil
}
