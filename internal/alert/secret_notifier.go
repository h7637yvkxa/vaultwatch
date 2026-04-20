package alert

import (
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/yourusername/vaultwatch/internal/vault"
)

// SecretNotifier prints KV v2 secret version metadata to a writer.
type SecretNotifier struct {
	w io.Writer
}

// NewSecretNotifier creates a SecretNotifier. If w is nil, os.Stdout is used.
func NewSecretNotifier(w io.Writer) *SecretNotifier {
	if w == nil {
		w = os.Stdout
	}
	return &SecretNotifier{w: w}
}

// Notify writes a summary of secret versions to the configured writer.
func (n *SecretNotifier) Notify(versions []vault.SecretVersion) error {
	if len(versions) == 0 {
		fmt.Fprintln(n.w, "[secret-versions] no versions found")
		return nil
	}

	sort.Slice(versions, func(i, j int) bool {
		return versions[i].Version < versions[j].Version
	})

	fmt.Fprintf(n.w, "[secret-versions] path=%s versions=%d\n", versions[0].Path, len(versions))
	for _, v := range versions {
		status := "active"
		if v.Destroyed {
			status = "destroyed"
		} else if v.DeletedAt != nil {
			status = "deleted"
		}
		fmt.Fprintf(n.w, "  version=%d created=%s status=%s\n",
			v.Version,
			v.CreatedAt.Format("2006-01-02T15:04:05Z"),
			status,
		)
	}
	return nil
}
