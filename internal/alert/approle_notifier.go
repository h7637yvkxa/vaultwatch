package alert

import (
	"fmt"
	"io"
	"os"

	"github.com/subtlepseudonym/vaultwatch/internal/vault"
)

// AppRoleNotifier prints AppRole information to a writer.
type AppRoleNotifier struct {
	w io.Writer
}

// NewAppRoleNotifier creates an AppRoleNotifier. If w is nil, os.Stdout is used.
func NewAppRoleNotifier(w io.Writer) *AppRoleNotifier {
	if w == nil {
		w = os.Stdout
	}
	return &AppRoleNotifier{w: w}
}

// Notify writes AppRole entries to the configured writer.
func (n *AppRoleNotifier) Notify(entries []vault.AppRoleEntry) error {
	if len(entries) == 0 {
		fmt.Fprintln(n.w, "[approle] no AppRoles found")
		return nil
	}

	fmt.Fprintf(n.w, "[approle] %d AppRole(s) found:\n", len(entries))
	for _, e := range entries {
		line := fmt.Sprintf("  - name=%s", e.Name)
		if e.RoleID != "" {
			line += fmt.Sprintf(" role_id=%s", e.RoleID)
		}
		if e.TokenTTL > 0 {
			line += fmt.Sprintf(" token_ttl=%ds", e.TokenTTL)
		}
		if e.TokenMaxTTL > 0 {
			line += fmt.Sprintf(" token_max_ttl=%ds", e.TokenMaxTTL)
		}
		fmt.Fprintln(n.w, line)
	}
	return nil
}
