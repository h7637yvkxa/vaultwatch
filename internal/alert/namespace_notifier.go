package alert

import (
	"fmt"
	"io"
	"os"

	"github.com/has/vaultwatch/internal/vault"
)

// NamespaceNotifier formats and writes namespace information.
type NamespaceNotifier struct {
	w io.Writer
}

// NewNamespaceNotifier creates a NamespaceNotifier writing to w.
// If w is nil, os.Stdout is used.
func NewNamespaceNotifier(w io.Writer) *NamespaceNotifier {
	if w == nil {
		w = os.Stdout
	}
	return &NamespaceNotifier{w: w}
}

// Notify prints namespace entries to the writer.
func (n *NamespaceNotifier) Notify(entries []vault.NamespaceEntry) error {
	if len(entries) == 0 {
		fmt.Fprintln(n.w, "[namespaces] no namespaces found")
		return nil
	}
	fmt.Fprintf(n.w, "[namespaces] %d namespace(s) discovered:\n", len(entries))
	for _, e := range entries {
		fmt.Fprintf(n.w, "  - path=%-30s id=%s\n", e.Path, e.ID)
	}
	return nil
}
