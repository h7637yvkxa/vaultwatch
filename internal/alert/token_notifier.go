package alert

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/yourusername/vaultwatch/internal/vault"
)

// TokenNotifier writes token metadata to a writer.
type TokenNotifier struct {
	w io.Writer
}

// NewTokenNotifier returns a TokenNotifier that writes to w.
// If w is nil, os.Stdout is used.
func NewTokenNotifier(w io.Writer) *TokenNotifier {
	if w == nil {
		w = os.Stdout
	}
	return &TokenNotifier{w: w}
}

// Notify prints token information to the configured writer.
func (tn *TokenNotifier) Notify(info *vault.TokenInfo) error {
	if info == nil {
		fmt.Fprintln(tn.w, "[token] no token info available")
		return nil
	}

	fmt.Fprintln(tn.w, "[token] current token info:")
	fmt.Fprintf(tn.w, "  id:           %s\n", info.ID)
	fmt.Fprintf(tn.w, "  display_name: %s\n", info.DisplayName)
	fmt.Fprintf(tn.w, "  ttl:          %s\n", info.TTL)
	fmt.Fprintf(tn.w, "  renewable:    %v\n", info.Renewable)
	if len(info.Policies) > 0 {
		fmt.Fprintf(tn.w, "  policies:     %s\n", strings.Join(info.Policies, ", "))
	}
	return nil
}
