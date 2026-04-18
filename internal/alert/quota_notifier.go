package alert

import (
	"fmt"
	"io"
	"os"

	"github.com/yourusername/vaultwatch/internal/vault"
)

// QuotaNotifier prints quota rule information to a writer.
type QuotaNotifier struct {
	w io.Writer
}

// NewQuotaNotifier returns a QuotaNotifier writing to w.
// If w is nil, os.Stdout is used.
func NewQuotaNotifier(w io.Writer) *QuotaNotifier {
	if w == nil {
		w = os.Stdout
	}
	return &QuotaNotifier{w: w}
}

// Notify writes a summary of quota rules to the configured writer.
func (n *QuotaNotifier) Notify(quotas []vault.QuotaRule) error {
	if len(quotas) == 0 {
		fmt.Fprintln(n.w, "[quotas] no quota rules found")
		return nil
	}

	fmt.Fprintf(n.w, "[quotas] %d rule(s) found:\n", len(quotas))
	for _, q := range quotas {
		switch q.Type {
		case "rate-limit":
			fmt.Fprintf(n.w, "  - [%s] %s: rate=%.0f/s burst=%.0f path=%q\n",
				q.Type, q.Name, q.Rate, q.Burst, q.Path)
		case "lease-count":
			fmt.Fprintf(n.w, "  - [%s] %s: max=%d path=%q\n",
				q.Type, q.Name, q.MaxLeases, q.Path)
		default:
			fmt.Fprintf(n.w, "  - [%s] %s: path=%q\n",
				q.Type, q.Name, q.Path)
		}
	}
	return nil
}
