package alert

import (
	"fmt"
	"io"
	"os"

	"github.com/arcanericky/vaultwatch/internal/vault"
)

// SSHNotifier prints SSH role information to a writer.
type SSHNotifier struct {
	w io.Writer
}

// NewSSHNotifier creates a new SSHNotifier. If w is nil, os.Stdout is used.
func NewSSHNotifier(w io.Writer) *SSHNotifier {
	if w == nil {
		w = os.Stdout
	}
	return &SSHNotifier{w: w}
}

// Notify writes SSH role details to the configured writer.
func (n *SSHNotifier) Notify(roles []vault.SSHRole) error {
	if len(roles) == 0 {
		fmt.Fprintln(n.w, "[SSH] No roles found.")
		return nil
	}
	fmt.Fprintf(n.w, "[SSH] %d role(s) found:\n", len(roles))
	for _, r := range roles {
		ttl := r.TTL
		if ttl == "" {
			ttl = "(default)"
		}
		maxTTL := r.MaxTTL
		if maxTTL == "" {
			maxTTL = "(default)"
		}
		fmt.Fprintf(n.w, "  - %s  key_type=%s  ttl=%s  max_ttl=%s  allowed_users=%s\n",
			r.Name, r.KeyType, ttl, maxTTL, r.AllowedUsers)
	}
	return nil
}
