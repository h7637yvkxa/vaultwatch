package alert

import (
	"fmt"
	"io"
	"os"

	"github.com/yourusername/vaultwatch/internal/vault"
)

// PKINotifier prints PKI certificate information to a writer.
type PKINotifier struct {
	w io.Writer
}

// NewPKINotifier creates a PKINotifier. If w is nil, os.Stdout is used.
func NewPKINotifier(w io.Writer) *PKINotifier {
	if w == nil {
		w = os.Stdout
	}
	return &PKINotifier{w: w}
}

// Notify writes PKI certificate details to the configured writer.
func (n *PKINotifier) Notify(mount string, certs []vault.PKICert) error {
	if len(certs) == 0 {
		fmt.Fprintf(n.w, "[PKI] No certificates found under mount: %s\n", mount)
		return nil
	}
	fmt.Fprintf(n.w, "[PKI] Certificates under mount '%s': %d\n", mount, len(certs))
	for _, c := range certs {
		expiry := "unknown"
		if !c.Expiry.IsZero() {
			expiry = c.Expiry.Format("2006-01-02 15:04:05 UTC")
		}
		ca := c.IssuingCA
		if ca == "" {
			ca = "n/a"
		}
		fmt.Fprintf(n.w, "  serial=%-30s expiry=%-25s ca=%s\n", c.SerialNumber, expiry, ca)
	}
	return nil
}
