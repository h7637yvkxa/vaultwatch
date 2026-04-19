package alert

import (
	"fmt"
	"io"
	"os"

	"github.com/yourusername/vaultwatch/internal/vault"
)

// LDAPNotifier prints LDAP group mappings to a writer.
type LDAPNotifier struct {
	w io.Writer
}

// NewLDAPNotifier creates a new LDAPNotifier. If w is nil, os.Stdout is used.
func NewLDAPNotifier(w io.Writer) *LDAPNotifier {
	if w == nil {
		w = os.Stdout
	}
	return &LDAPNotifier{w: w}
}

// Notify writes LDAP group information to the writer.
func (n *LDAPNotifier) Notify(groups []vault.LDAPGroup) error {
	if len(groups) == 0 {
		fmt.Fprintln(n.w, "[LDAP] No group mappings found.")
		return nil
	}
	fmt.Fprintf(n.w, "[LDAP] %d group mapping(s):\n", len(groups))
	for _, g := range groups {
		policies := "(none)"
		if len(g.Policies) > 0 {
			policies = fmt.Sprintf("%v", g.Policies)
		}
		fmt.Fprintf(n.w, "  - %s  policies=%s\n", g.Name, policies)
	}
	return nil
}
