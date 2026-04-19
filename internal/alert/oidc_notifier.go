package alert

import (
	"fmt"
	"io"
	"os"

	"github.com/morningconsult/vaultwatch/internal/vault"
)

// OIDCNotifier writes OIDC role information to a writer.
type OIDCNotifier struct {
	w io.Writer
}

// NewOIDCNotifier returns a new OIDCNotifier. If w is nil, os.Stdout is used.
func NewOIDCNotifier(w io.Writer) *OIDCNotifier {
	if w == nil {
		w = os.Stdout
	}
	return &OIDCNotifier{w: w}
}

// Notify prints OIDC role details.
func (n *OIDCNotifier) Notify(roles []vault.OIDCRole) error {
	if len(roles) == 0 {
		fmt.Fprintln(n.w, "[oidc] no roles configured")
		return nil
	}
	fmt.Fprintf(n.w, "[oidc] %d role(s) configured:\n", len(roles))
	for _, r := range roles {
		ttl := r.TTL
		if ttl == "" {
			ttl = "default"
		}
		fmt.Fprintf(n.w, "  - name=%-20s user_claim=%-15s ttl=%s\n",
			r.Name, r.UserClaim, ttl)
	}
	return nil
}
