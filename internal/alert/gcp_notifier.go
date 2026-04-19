package alert

import (
	"fmt"
	"io"
	"os"

	"github.com/arussellsaw/vaultwatch/internal/vault"
)

// GCPNotifier prints GCP role information to a writer.
type GCPNotifier struct {
	w io.Writer
}

// NewGCPNotifier creates a GCPNotifier. If w is nil, os.Stdout is used.
func NewGCPNotifier(w io.Writer) *GCPNotifier {
	if w == nil {
		w = os.Stdout
	}
	return &GCPNotifier{w: w}
}

// Notify writes GCP role details to the configured writer.
func (n *GCPNotifier) Notify(roles []vault.GCPRole) error {
	if len(roles) == 0 {
		fmt.Fprintln(n.w, "[gcp] no roles found")
		return nil
	}
	fmt.Fprintf(n.w, "[gcp] %d role(s):\n", len(roles))
	for _, r := range roles {
		line := fmt.Sprintf("  - %s", r.Name)
		if r.SecretType != "" {
			line += fmt.Sprintf(" (type: %s)", r.SecretType)
		}
		if r.Project != "" {
			line += fmt.Sprintf(" project=%s", r.Project)
		}
		fmt.Fprintln(n.w, line)
	}
	return nil
}
