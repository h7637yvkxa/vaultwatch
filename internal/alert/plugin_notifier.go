package alert

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/your-org/vaultwatch/internal/vault"
)

// PluginNotifier writes a summary of registered Vault plugins to a writer.
type PluginNotifier struct {
	w io.Writer
}

// NewPluginNotifier returns a PluginNotifier that writes to w.
// If w is nil, os.Stdout is used.
func NewPluginNotifier(w io.Writer) *PluginNotifier {
	if w == nil {
		w = os.Stdout
	}
	return &PluginNotifier{w: w}
}

// Notify prints plugin information grouped by type.
func (n *PluginNotifier) Notify(plugins []vault.PluginInfo) error {
	if len(plugins) == 0 {
		fmt.Fprintln(n.w, "[plugins] no plugins registered")
		return nil
	}

	groups := make(map[string][]string)
	for _, p := range plugins {
		groups[p.Type] = append(groups[p.Type], p.Name)
	}

	fmt.Fprintf(n.w, "[plugins] %d plugin(s) registered\n", len(plugins))
	for typ, names := range groups {
		fmt.Fprintf(n.w, "  type=%-12s plugins=%s\n", typ, strings.Join(names, ", "))
	}
	return nil
}
