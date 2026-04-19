package alert

import (
	"fmt"
	"io"
	"os"

	"github.com/wadefletch/vaultwatch/internal/vault"
)

// GroupNotifier prints identity group information to a writer.
type GroupNotifier struct {
	w io.Writer
}

// NewGroupNotifier returns a GroupNotifier writing to w.
// If w is nil, os.Stdout is used.
func NewGroupNotifier(w io.Writer) *GroupNotifier {
	if w == nil {
		w = os.Stdout
	}
	return &GroupNotifier{w: w}
}

// Notify writes a summary of Vault identity groups.
func (n *GroupNotifier) Notify(groups []vault.GroupEntry) {
	if len(groups) == 0 {
		fmt.Fprintln(n.w, "[groups] no identity groups found")
		return
	}
	fmt.Fprintf(n.w, "[groups] %d identity group(s) found:\n", len(groups))
	for _, g := range groups {
		policies := "(none)"
		if len(g.Policies) > 0 {
			policies = fmt.Sprintf("%v", g.Policies)
		}
		fmt.Fprintf(n.w, "  - %s (id=%s, type=%s, policies=%s)\n",
			g.Name, g.ID, g.Type, policies)
	}
}
