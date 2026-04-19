package alert

import (
	"fmt"
	"io"
	"os"

	"github.com/arturhoo/vaultwatch/internal/vault"
)

// AzureNotifier prints Azure role information to a writer.
type AzureNotifier struct {
	w io.Writer
}

// NewAzureNotifier creates an AzureNotifier. If w is nil, os.Stdout is used.
func NewAzureNotifier(w io.Writer) *AzureNotifier {
	if w == nil {
		w = os.Stdout
	}
	return &AzureNotifier{w: w}
}

// Notify writes Azure role details to the configured writer.
func (n *AzureNotifier) Notify(roles []vault.AzureRole) error {
	if len(roles) == 0 {
		fmt.Fprintln(n.w, "[azure] no roles found")
		return nil
	}
	fmt.Fprintf(n.w, "[azure] %d role(s):\n", len(roles))
	for _, r := range roles {
		ttl := r.TTL
		if ttl == "" {
			ttl = "default"
		}
		maxTTL := r.MaxTTL
		if maxTTL == "" {
			maxTTL = "default"
		}
		fmt.Fprintf(n.w, "  - %s (ttl=%s, max_ttl=%s, app_object_id=%s)\n",
			r.Name, ttl, maxTTL, r.ApplicationObjectID)
	}
	return nil
}
