package alert

import (
	"fmt"
	"io"
	"os"

	"github.com/yourusername/vaultwatch/internal/vault"
)

// RADIUSNotifier prints RADIUS user information to a writer.
type RADIUSNotifier struct {
	w io.Writer
}

// NewRADIUSNotifier creates a RADIUSNotifier. If w is nil, os.Stdout is used.
func NewRADIUSNotifier(w io.Writer) *RADIUSNotifier {
	if w == nil {
		w = os.Stdout
	}
	return &RADIUSNotifier{w: w}
}

// Notify writes a summary of RADIUS users to the configured writer.
func (n *RADIUSNotifier) Notify(users []vault.RADIUSUser) error {
	if len(users) == 0 {
		fmt.Fprintln(n.w, "[RADIUS] No users configured.")
		return nil
	}
	fmt.Fprintf(n.w, "[RADIUS] %d user(s) configured:\n", len(users))
	for _, u := range users {
		policies := "(none)"
		if len(u.Policies) > 0 {
			policies = fmt.Sprintf("%v", u.Policies)
		}
		fmt.Fprintf(n.w, "  - %s  policies=%s\n", u.Username, policies)
	}
	return nil
}
