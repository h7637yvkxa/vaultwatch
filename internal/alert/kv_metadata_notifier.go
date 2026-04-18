package alert

import (
	"fmt"
	"io"
	"os"

	"github.com/rzp/vaultwatch/internal/vault"
)

// KVMetadataNotifier prints KV v2 secret metadata to a writer.
type KVMetadataNotifier struct {
	w io.Writer
}

// NewKVMetadataNotifier creates a KVMetadataNotifier. If w is nil, os.Stdout is used.
func NewKVMetadataNotifier(w io.Writer) *KVMetadataNotifier {
	if w == nil {
		w = os.Stdout
	}
	return &KVMetadataNotifier{w: w}
}

// Notify writes metadata details for each entry.
func (n *KVMetadataNotifier) Notify(entries []*vault.KVMetadata) {
	if len(entries) == 0 {
		fmt.Fprintln(n.w, "[kv-metadata] no entries to report")
		return
	}
	fmt.Fprintf(n.w, "[kv-metadata] %d secret(s)\n", len(entries))
	for _, m := range entries {
		if m == nil {
			continue
		}
		fmt.Fprintf(n.w, "  path=%-40s current_version=%-4d max_versions=%-4d updated=%s\n",
			m.Path,
			m.CurrentVersion,
			m.MaxVersions,
			m.UpdatedTime.Format("2006-01-02T15:04:05Z"),
		)
		if m.DeleteVersionAfter != "" && m.DeleteVersionAfter != "0s" {
			fmt.Fprintf(n.w, "    delete_version_after=%s\n", m.DeleteVersionAfter)
		}
	}
}
