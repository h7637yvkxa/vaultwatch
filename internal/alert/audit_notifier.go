package alert

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/yourusername/vaultwatch/internal/vault"
)

// AuditNotifier reports audit device status to a writer.
type AuditNotifier struct {
	w io.Writer
}

// NewAuditNotifier creates an AuditNotifier writing to w.
// If w is nil, os.Stdout is used.
func NewAuditNotifier(w io.Writer) *AuditNotifier {
	if w == nil {
		w = os.Stdout
	}
	return &AuditNotifier{w: w}
}

// Notify prints a summary of audit devices.
func (a *AuditNotifier) Notify(entries []vault.AuditEntry) error {
	if len(entries) == 0 {
		fmt.Fprintln(a.w, "[audit] no audit devices found")
		return nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("[audit] %d device(s) enabled:\n", len(entries)))
	for _, e := range entries {
		desc := e.Description
		if desc == "" {
			desc = "(no description)"
		}
		sb.WriteString(fmt.Sprintf("  - %s [%s] %s\n", e.Path, e.Type, desc))
	}

	_, err := fmt.Fprint(a.w, sb.String())
	return err
}
