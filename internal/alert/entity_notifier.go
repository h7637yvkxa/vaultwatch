package alert

import (
	"fmt"
	"io"
	"os"

	"github.com/elizabethwanjiku703/vaultwatch/internal/vault"
)

// EntityNotifier prints identity entity information.
type EntityNotifier struct {
	w io.Writer
}

// NewEntityNotifier creates an EntityNotifier writing to w (defaults to stdout).
func NewEntityNotifier(w io.Writer) *EntityNotifier {
	if w == nil {
		w = os.Stdout
	}
	return &EntityNotifier{w: w}
}

// Notify prints a summary of Vault identity entities.
func (n *EntityNotifier) Notify(entries []vault.EntityEntry) error {
	if len(entries) == 0 {
		fmt.Fprintln(n.w, "[entities] no identity entities found")
		return nil
	}

	fmt.Fprintf(n.w, "[entities] %d identity entity/entities:\n", len(entries))
	for _, e := range entries {
		status := "enabled"
		if e.Disabled {
			status = "disabled"
		}
		fmt.Fprintf(n.w, "  - %s (id=%s, policies=%v, status=%s)\n", e.Name, e.ID, e.Policies, status)
	}
	return nil
}
