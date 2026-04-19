package alert

import (
	"fmt"
	"io"
	"os"

	"github.com/subtlepseudonym/vaultwatch/internal/vault"
)

// RoleNotifier prints role information for a given auth mount.
type RoleNotifier struct {
	Mount  string
	Roles  []vault.RoleEntry
	Writer io.Writer
}

// NewRoleNotifier creates a RoleNotifier writing to w (defaults to stdout).
func NewRoleNotifier(mount string, roles []vault.RoleEntry, w io.Writer) *RoleNotifier {
	if w == nil {
		w = os.Stdout
	}
	return &RoleNotifier{Mount: mount, Roles: roles, Writer: w}
}

// Notify writes a summary of roles to the configured writer.
func (n *RoleNotifier) Notify() error {
	if len(n.Roles) == 0 {
		fmt.Fprintf(n.Writer, "[roles] mount=%s no roles found\n", n.Mount)
		return nil
	}
	fmt.Fprintf(n.Writer, "[roles] mount=%s count=%d\n", n.Mount, len(n.Roles))
	for _, r := range n.Roles {
		fmt.Fprintf(n.Writer, "  - name=%-30s path=%s token_ttl=%d max_ttl=%d\n",
			r.Name, r.Path, r.TokenTTL, r.MaxTTL)
	}
	return nil
}
