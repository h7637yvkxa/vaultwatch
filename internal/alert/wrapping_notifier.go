package alert

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/subtlepseudonym/vaultwatch/internal/vault"
)

// WrappingNotifier prints wrapping token info to a writer.
type WrappingNotifier struct {
	w io.Writer
}

// NewWrappingNotifier creates a WrappingNotifier. If w is nil, os.Stdout is used.
func NewWrappingNotifier(w io.Writer) *WrappingNotifier {
	if w == nil {
		w = os.Stdout
	}
	return &WrappingNotifier{w: w}
}

// Notify prints wrapping token metadata, warning if TTL is below threshold.
func (n *WrappingNotifier) Notify(info *vault.WrappingInfo, warnBelow time.Duration) error {
	if info == nil {
		fmt.Fprintln(n.w, "[wrapping] no wrapping info available")
		return nil
	}

	status := "OK"
	if info.TTL <= warnBelow {
		status = "WARNING"
	}

	fmt.Fprintf(n.w, "[wrapping] status=%-8s token=%s accessor=%s ttl=%s path=%s created=%s\n",
		status,
		info.Token,
		info.Accessor,
		info.TTL.Round(time.Second),
		info.CreationPath,
		info.CreationTime,
	)
	return nil
}
