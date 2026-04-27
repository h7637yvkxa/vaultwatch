package alert

import (
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/yourusername/vaultwatch/internal/vault"
)

// LogicalBackendNotifier formats and writes logical backend information.
type LogicalBackendNotifier struct {
	w io.Writer
}

// NewLogicalBackendNotifier creates a new LogicalBackendNotifier.
// If w is nil, os.Stdout is used.
func NewLogicalBackendNotifier(w io.Writer) *LogicalBackendNotifier {
	if w == nil {
		w = os.Stdout
	}
	return &LogicalBackendNotifier{w: w}
}

// Notify writes logical backend details to the configured writer.
func (n *LogicalBackendNotifier) Notify(backends []vault.LogicalBackend) error {
	if len(backends) == 0 {
		fmt.Fprintln(n.w, "[logical-backends] no backends found")
		return nil
	}

	sort.Slice(backends, func(i, j int) bool {
		return backends[i].Path < backends[j].Path
	})

	fmt.Fprintf(n.w, "[logical-backends] %d backend(s) mounted:\n", len(backends))
	for _, b := range backends {
		desc := b.Description
		if desc == "" {
			desc = "(no description)"
		}
		flags := ""
		if b.Local {
			flags += " local=true"
		}
		if b.SealWrap {
			flags += " seal_wrap=true"
		}
		fmt.Fprintf(n.w, "  path=%-30s type=%-12s desc=%s%s\n", b.Path, b.Type, desc, flags)
	}
	return nil
}
