package alert

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/drew/vaultwatch/internal/vault"
)

// UserpassNotifier prints userpass users to a writer.
type UserpassNotifier struct {
	w io.Writer
}

// NewUserpassNotifier creates a UserpassNotifier. If w is nil, os.Stdout is used.
func NewUserpassNotifier(w io.Writer) *UserpassNotifier {
	if w == nil {
		w = os.Stdout
	}
	return &UserpassNotifier{w: w}
}

// Notify writes a summary of userpass users.
func (n *UserpassNotifier) Notify(users []vault.UserpassUser) error {
	if len(users) == 0 {
		fmt.Fprintln(n.w, "[userpass] no users found")
		return nil
	}
	fmt.Fprintf(n.w, "[userpass] %d user(s):\n", len(users))
	for _, u := range users {
		fmt.Fprintf(n.w, "  - %s  policies=%s\n", u.Username, formatPolicies(u.Policies))
	}
	return nil
}

// formatPolicies returns a human-readable representation of a policy list.
func formatPolicies(policies []string) string {
	if len(policies) == 0 {
		return "(none)"
	}
	return strings.Join(policies, ", ")
}
