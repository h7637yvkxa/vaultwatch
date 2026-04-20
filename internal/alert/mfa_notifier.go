package alert

import (
	"fmt"
	"io"
	"os"

	"github.com/arjunsriva/vaultwatch/internal/vault"
)

// MFANotifier prints MFA method information to a writer.
type MFANotifier struct {
	w io.Writer
}

// NewMFANotifier creates a new MFANotifier. If w is nil, os.Stdout is used.
func NewMFANotifier(w io.Writer) *MFANotifier {
	if w == nil {
		w = os.Stdout
	}
	return &MFANotifier{w: w}
}

// Notify writes MFA method details to the configured writer.
func (n *MFANotifier) Notify(methods []vault.MFAMethod) {
	if len(methods) == 0 {
		fmt.Fprintln(n.w, "[MFA] No MFA methods configured.")
		return
	}
	fmt.Fprintf(n.w, "[MFA] %d method(s) configured:\n", len(methods))
	for _, m := range methods {
		name := m.Name
		if name == "" {
			name = "(unnamed)"
		}
		fmt.Fprintf(n.w, "  - id=%-36s  type=%-10s  name=%s\n", m.ID, m.Type, name)
	}
}
