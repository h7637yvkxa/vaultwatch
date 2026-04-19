package alert

import (
	"fmt"
	"io"
	"os"

	"github.com/rbhz/vaultwatch/internal/vault"
)

// AWSNotifier prints AWS auth roles to a writer.
type AWSNotifier struct {
	w io.Writer
}

// NewAWSNotifier creates a new AWSNotifier. If w is nil, os.Stdout is used.
func NewAWSNotifier(w io.Writer) *AWSNotifier {
	if w == nil {
		w = os.Stdout
	}
	return &AWSNotifier{w: w}
}

// Notify writes AWS role information to the configured writer.
func (n *AWSNotifier) Notify(roles []vault.AWSRole) error {
	if len(roles) == 0 {
		fmt.Fprintln(n.w, "[AWS] No AWS auth roles found.")
		return nil
	}
	fmt.Fprintf(n.w, "[AWS] %d role(s) configured:\n", len(roles))
	for _, r := range roles {
		line := fmt.Sprintf("  - %s", r.Name)
		if r.AuthType != "" {
			line += fmt.Sprintf(" (auth_type=%s)", r.AuthType)
		}
		if r.BoundAMI != "" {
			line += fmt.Sprintf(" bound_ami=%s", r.BoundAMI)
		}
		if len(r.Policies) > 0 {
			line += fmt.Sprintf(" policies=%v", r.Policies)
		}
		fmt.Fprintln(n.w, line)
	}
	return nil
}
