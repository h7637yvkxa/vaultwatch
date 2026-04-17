package alert

import (
	"fmt"
	"io"
	"os"

	"github.com/rzp-labs/vaultwatch/internal/vault"
)

// AuthNotifier prints enabled auth methods to a writer.
type AuthNotifier struct {
	w io.Writer
}

// NewAuthNotifier creates a new AuthNotifier. If w is nil, os.Stdout is used.
func NewAuthNotifier(w io.Writer) *AuthNotifier {
	if w == nil {
		w = os.Stdout
	}
	return &AuthNotifier{w: w}
}

// Notify writes auth method information to the writer.
func (n *AuthNotifier) Notify(methods []vault.AuthMethod) error {
	if len(methods) == 0 {
		fmt.Fprintln(n.w, "[auth] no auth methods found")
		return nil
	}

	fmt.Fprintf(n.w, "[auth] %d auth method(s) enabled:\n", len(methods))
	for _, m := range methods {
		local := ""
		if m.Local {
			local = " (local)"
		}
		fmt.Fprintf(n.w, "  %-20s type=%-12s accessor=%s%s\n",
			m.Path, m.Type, m.Accessor, local)
		if m.Description != "" {
			fmt.Fprintf(n.w, "    description: %s\n", m.Description)
		}
	}
	return nil
}
