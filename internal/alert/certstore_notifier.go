package alert

import (
	"fmt"
	"io"
	"os"

	"github.com/densestvoid/vaultwatch/internal/vault"
)

// CertStoreNotifier prints cert auth entries to a writer.
type CertStoreNotifier struct {
	writer io.Writer
	mount  string
}

// NewCertStoreNotifier creates a CertStoreNotifier. If w is nil, os.Stdout is used.
func NewCertStoreNotifier(w io.Writer, mount string) *CertStoreNotifier {
	if w == nil {
		w = os.Stdout
	}
	return &CertStoreNotifier{writer: w, mount: mount}
}

// Notify writes cert store entries to the configured writer.
func (n *CertStoreNotifier) Notify(entries []vault.CertStoreEntry) error {
	if len(entries) == 0 {
		fmt.Fprintf(n.writer, "[cert-store/%s] no certificate entries found\n", n.mount)
		return nil
	}
	fmt.Fprintf(n.writer, "[cert-store/%s] %d certificate(s):\n", n.mount, len(entries))
	for _, e := range entries {
		fmt.Fprintf(n.writer, "  - %s\n", e.Name)
	}
	return nil
}
