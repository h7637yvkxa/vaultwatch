package alert

import (
	"fmt"
	"io"
	"os"
)

// EGPPolicy represents a single Endpoint Governing Policy entry.
type EGPPolicy struct {
	Name  string
	Paths []string
}

// EGPNotifier prints EGP policy information to a writer.
type EGPNotifier struct {
	w io.Writer
}

// NewEGPNotifier returns an EGPNotifier that writes to w.
// If w is nil, os.Stdout is used.
func NewEGPNotifier(w io.Writer) *EGPNotifier {
	if w == nil {
		w = os.Stdout
	}
	return &EGPNotifier{w: w}
}

// Notify writes a summary of EGP policies to the notifier's writer.
func (n *EGPNotifier) Notify(policies []EGPPolicy) error {
	if len(policies) == 0 {
		fmt.Fprintln(n.w, "[EGP] No endpoint governing policies found.")
		return nil
	}
	fmt.Fprintf(n.w, "[EGP] %d endpoint governing policy(ies):\n", len(policies))
	for _, p := range policies {
		if len(p.Paths) > 0 {
			fmt.Fprintf(n.w, "  - %s (paths: %v)\n", p.Name, p.Paths)
		} else {
			fmt.Fprintf(n.w, "  - %s (no paths)\n", p.Name)
		}
	}
	return nil
}
