package alert

import (
	"fmt"
	"io"
	"os"

	"github.com/yourusername/vaultwatch/internal/vault"
)

// SealNotifier writes seal status information to an io.Writer.
type SealNotifier struct {
	w io.Writer
}

// NewSealNotifier creates a SealNotifier. If w is nil, os.Stdout is used.
func NewSealNotifier(w io.Writer) *SealNotifier {
	if w == nil {
		w = os.Stdout
	}
	return &SealNotifier{w: w}
}

// Notify prints the seal status. Returns an error if Vault is sealed.
func (n *SealNotifier) Notify(status *vault.SealStatus) error {
	if status == nil {
		fmt.Fprintln(n.w, "[seal] no seal status available")
		return nil
	}
	state := "unsealed"
	if status.Sealed {
		state = "SEALED"
	}
	init := "initialized"
	if !status.Initialized {
		init = "NOT initialized"
	}
	fmt.Fprintf(n.w, "[seal] Vault is %s | %s | version=%s | cluster=%s | checked=%s\n",
		state, init, status.Version, status.ClusterName, status.CheckedAt.Format("2006-01-02T15:04:05Z"))
	if status.Sealed {
		return fmt.Errorf("vault is sealed (progress %d/%d)", status.Progress, status.T)
	}
	return nil
}

// NotifyMany calls Notify for each status in the slice, collecting all errors.
// It returns a combined error if any Vault instance is sealed, or nil if all are unsealed.
func (n *SealNotifier) NotifyMany(statuses []*vault.SealStatus) error {
	var errs []error
	for _, s := range statuses {
		if err := n.Notify(s); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) == 0 {
		return nil
	}
	return fmt.Errorf("%d vault instance(s) sealed: %v", len(errs), errs)
}
