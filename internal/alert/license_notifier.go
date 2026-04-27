package alert

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/yourusername/vaultwatch/internal/vault"
)

// LicenseNotifier formats and writes Vault Enterprise license information.
type LicenseNotifier struct {
	w io.Writer
}

// NewLicenseNotifier creates a LicenseNotifier that writes to w.
// If w is nil, os.Stdout is used.
func NewLicenseNotifier(w io.Writer) *LicenseNotifier {
	if w == nil {
		w = os.Stdout
	}
	return &LicenseNotifier{w: w}
}

// Notify writes license details to the configured writer.
func (n *LicenseNotifier) Notify(info *vault.LicenseInfo) error {
	if info == nil {
		fmt.Fprintln(n.w, "[license] no license information available")
		return nil
	}

	timeUntil := time.Until(info.ExpirationTime)
	status := "OK"
	switch {
	case info.Terminated:
		status = "TERMINATED"
	case timeUntil < 0:
		status = "EXPIRED"
	case timeUntil < 7*24*time.Hour:
		status = "CRITICAL"
	case timeUntil < 30*24*time.Hour:
		status = "WARNING"
	}

	fmt.Fprintf(n.w, "[license] id=%s customer=%q status=%s expires=%s (in %.0f days)\n",
		info.LicenseID,
		info.CustomerName,
		status,
		info.ExpirationTime.Format(time.RFC3339),
		timeUntil.Hours()/24,
	)

	if len(info.Features) > 0 {
		fmt.Fprintf(n.w, "[license] features: %v\n", info.Features)
	}
	return nil
}
