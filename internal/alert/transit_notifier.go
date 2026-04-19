package alert

import (
	"fmt"
	"io"
	"os"

	"github.com/wes/vaultwatch/internal/vault"
)

// TransitNotifier reports transit key information to a writer.
type TransitNotifier struct {
	w io.Writer
}

// NewTransitNotifier creates a TransitNotifier writing to w.
// If w is nil, os.Stdout is used.
func NewTransitNotifier(w io.Writer) *TransitNotifier {
	if w == nil {
		w = os.Stdout
	}
	return &TransitNotifier{w: w}
}

// Notify prints a summary of the provided transit keys.
func (n *TransitNotifier) Notify(keys []vault.TransitKey) error {
	if len(keys) == 0 {
		fmt.Fprintln(n.w, "[transit] no transit keys found")
		return nil
	}
	fmt.Fprintf(n.w, "[transit] %d key(s) found:\n", len(keys))
	for _, k := range keys {
		deletion := "no"
		if k.DeletionAllowed {
			deletion = "yes"
		}
		exportable := "no"
		if k.Exportable {
			exportable = "yes"
		}
		fmt.Fprintf(n.w, "  - %s (type=%s, latest_version=%d, deletion_allowed=%s, exportable=%s)\n",
			k.Name, k.Type, k.LatestVersion, deletion, exportable)
	}
	return nil
}
