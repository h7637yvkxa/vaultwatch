package alert

import (
	"fmt"
	"io"
	"os"

	"github.com/yourusername/vaultwatch/internal/vault"
)

// TokenRoleNotifier prints token role information.
type TokenRoleNotifier struct {
	w io.Writer
}

// NewTokenRoleNotifier creates a new TokenRoleNotifier. If w is nil, os.Stdout is used.
func NewTokenRoleNotifier(w io.Writer) *TokenRoleNotifier {
	if w == nil {
		w = os.Stdout
	}
	return &TokenRoleNotifier{w: w}
}

// Notify writes token role details to the configured writer.
func (n *TokenRoleNotifier) Notify(roles []vault.TokenRole) error {
	if len(roles) == 0 {
		fmt.Fprintln(n.w, "[token-roles] no token roles found")
		return nil
	}
	fmt.Fprintf(n.w, "[token-roles] %d role(s) defined:\n", len(roles))
	for _, r := range roles {
		policies := "(none)"
		if len(r.AllowedPolicies) > 0 {
			policies = fmt.Sprintf("%v", r.AllowedPolicies)
		}
		fmt.Fprintf(n.w, "  - %s | orphan=%v renewable=%v max_ttl=%d period=%d policies=%s\n",
			r.Name, r.Orphan, r.Renewable, r.ExplicitMaxTTL, r.Period, policies)
	}
	return nil
}
