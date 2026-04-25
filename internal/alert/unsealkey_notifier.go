package alert

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/yourusername/vaultwatch/internal/vault"
)

// UnsealKeyNotifier formats and writes unseal key status to a writer.
type UnsealKeyNotifier struct {
	w io.Writer
}

// NewUnsealKeyNotifier creates a new UnsealKeyNotifier. If w is nil, os.Stdout is used.
func NewUnsealKeyNotifier(w io.Writer) *UnsealKeyNotifier {
	if w == nil {
		w = os.Stdout
	}
	return &UnsealKeyNotifier{w: w}
}

// Notify writes the unseal key status report to the writer.
func (n *UnsealKeyNotifier) Notify(status *vault.UnsealKeyStatus) error {
	if status == nil {
		fmt.Fprintln(n.w, "[unseal-key] no status available")
		return nil
	}

	fmt.Fprintln(n.w, "[unseal-key] Unseal Key Configuration")
	fmt.Fprintf(n.w, "  Secret Shares    : %d\n", status.SecretShares)
	fmt.Fprintf(n.w, "  Secret Threshold : %d\n", status.SecretThreshold)
	fmt.Fprintf(n.w, "  Stored Shares    : %d\n", status.StoredShares)

	if status.Nonce != "" {
		fmt.Fprintf(n.w, "  Nonce            : %s\n", status.Nonce)
	}

	if len(status.PGPFingerprints) > 0 {
		fmt.Fprintf(n.w, "  PGP Fingerprints : %s\n", strings.Join(status.PGPFingerprints, ", "))
	} else {
		fmt.Fprintln(n.w, "  PGP Fingerprints : none")
	}

	if status.SecretThreshold > 0 && status.SecretShares > 0 {
		fmt.Fprintf(n.w, "  Quorum           : %d of %d shares required\n",
			status.SecretThreshold, status.SecretShares)
	}
	return nil
}
