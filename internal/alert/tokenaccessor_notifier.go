package alert

import (
	"fmt"
	"io"
	"os"

	"github.com/danhale-git/vaultwatch/internal/vault"
)

// TokenAccessorNotifier writes token accessor information to a writer.
type TokenAccessorNotifier struct {
	w io.Writer
}

// NewTokenAccessorNotifier creates a new TokenAccessorNotifier.
// If w is nil, os.Stdout is used.
func NewTokenAccessorNotifier(w io.Writer) *TokenAccessorNotifier {
	if w == nil {
		w = os.Stdout
	}
	return &TokenAccessorNotifier{w: w}
}

// Notify writes the list of token accessors to the configured writer.
func (n *TokenAccessorNotifier) Notify(entries []vault.TokenAccessorEntry) error {
	if len(entries) == 0 {
		_, err := fmt.Fprintln(n.w, "[token-accessors] no accessors found")
		return err
	}

	_, err := fmt.Fprintf(n.w, "[token-accessors] %d accessor(s) active:\n", len(entries))
	if err != nil {
		return err
	}

	for _, e := range entries {
		display := e.DisplayName
		if display == "" {
			display = "(no display name)"
		}
		_, err := fmt.Fprintf(n.w, "  accessor=%s display=%s ttl=%ds\n",
			e.Accessor, display, e.TTL)
		if err != nil {
			return err
		}
	}
	return nil
}
