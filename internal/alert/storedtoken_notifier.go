package alert

import (
	"fmt"
	"io"
	"os"

	"github.com/has-bii/vaultwatch/internal/vault"
)

// StoredTokenNotifier formats and writes stored token accessor information.
type StoredTokenNotifier struct {
	result *vault.StoredTokenResult
	w      io.Writer
}

// NewStoredTokenNotifier creates a new StoredTokenNotifier.
func NewStoredTokenNotifier(result *vault.StoredTokenResult, w io.Writer) *StoredTokenNotifier {
	if w == nil {
		w = os.Stdout
	}
	return &StoredTokenNotifier{result: result, w: w}
}

// Notify writes stored token accessor entries to the configured writer.
func (n *StoredTokenNotifier) Notify() error {
	if n.result == nil {
		fmt.Fprintln(n.w, "[stored-tokens] no result available")
		return nil
	}
	if len(n.result.Entries) == 0 {
		fmt.Fprintln(n.w, "[stored-tokens] no token accessors found")
		return nil
	}
	fmt.Fprintf(n.w, "[stored-tokens] %d token accessor(s) found:\n", len(n.result.Entries))
	for _, e := range n.result.Entries {
		line := fmt.Sprintf("  accessor=%s", e.Accessor)
		if e.DisplayName != "" {
			line += fmt.Sprintf(" display_name=%s", e.DisplayName)
		}
		if len(e.Policies) > 0 {
			line += fmt.Sprintf(" policies=%v", e.Policies)
		}
		if e.TTL > 0 {
			line += fmt.Sprintf(" ttl=%ds", e.TTL)
		}
		if e.Renewable {
			line += " renewable=true"
		}
		fmt.Fprintln(n.w, line)
	}
	return nil
}
