package alert

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/yourusername/vaultwatch/internal/vault"
)

// CapabilitiesNotifier prints token capability results to a writer.
type CapabilitiesNotifier struct {
	w io.Writer
}

// NewCapabilitiesNotifier creates a CapabilitiesNotifier. If w is nil, os.Stdout is used.
func NewCapabilitiesNotifier(w io.Writer) *CapabilitiesNotifier {
	if w == nil {
		w = os.Stdout
	}
	return &CapabilitiesNotifier{w: w}
}

// Notify writes a formatted summary of capability results.
func (n *CapabilitiesNotifier) Notify(results []vault.CapabilityResult) error {
	if len(results) == 0 {
		fmt.Fprintln(n.w, "[capabilities] no paths checked")
		return nil
	}

	fmt.Fprintln(n.w, "[capabilities] token capability report:")
	for _, r := range results {
		caps := "(none)"
		if len(r.Capabilities) > 0 {
			caps = strings.Join(r.Capabilities, ", ")
		}
		denied := false
		for _, c := range r.Capabilities {
			if c == "deny" {
				denied = true
				break
			}
		}
		prefix := "  "
		if denied {
			prefix = "  [DENIED] "
		}
		fmt.Fprintf(n.w, "%spath=%s caps=[%s]\n", prefix, r.Path, caps)
	}
	return nil
}
