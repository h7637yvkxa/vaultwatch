package alert

import (
	"fmt"
	"io"
	"os"

	"github.com/yourusername/vaultwatch/internal/vault"
)

// RGPNotifier prints RGP policy information to a writer.
type RGPNotifier struct {
	w io.Writer
}

// NewRGPNotifier creates a new RGPNotifier. If w is nil, os.Stdout is used.
func NewRGPNotifier(w io.Writer) *RGPNotifier {
	if w == nil {
		w = os.Stdout
	}
	return &RGPNotifier{w: w}
}

// Notify writes RGP policy details to the writer.
func (n *RGPNotifier) Notify(policies []vault.RGPPolicy) error {
	if len(policies) == 0 {
		fmt.Fprintln(n.w, "[RGP] No RGP policies found.")
		return nil
	}
	fmt.Fprintf(n.w, "[RGP] %d policy(ies) found:\n", len(policies))
	for _, p := range policies {
		level := p.EnforcementLevel
		if level == "" {
			level = "unknown"
		}
		fmt.Fprintf(n.w, "  - %s (enforcement: %s)\n", p.Name, level)
	}
	return nil
}
