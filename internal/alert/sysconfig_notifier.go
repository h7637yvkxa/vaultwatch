package alert

import (
	"fmt"
	"io"
	"os"

	"github.com/subtlepseudonym/vaultwatch/internal/vault"
)

// SysConfigNotifier prints Vault system configuration to a writer.
type SysConfigNotifier struct {
	w io.Writer
}

// NewSysConfigNotifier creates a SysConfigNotifier. If w is nil, os.Stdout is used.
func NewSysConfigNotifier(w io.Writer) *SysConfigNotifier {
	if w == nil {
		w = os.Stdout
	}
	return &SysConfigNotifier{w: w}
}

// Notify writes the system configuration details to the writer.
func (n *SysConfigNotifier) Notify(cfg *vault.SysConfig) error {
	if cfg == nil {
		_, err := fmt.Fprintln(n.w, "[sysconfig] no configuration available")
		return err
	}

	fmt.Fprintln(n.w, "[sysconfig] Vault System Configuration:")
	fmt.Fprintf(n.w, "  default_lease_ttl : %s\n", cfg.DefaultLeaseTTL)
	fmt.Fprintf(n.w, "  max_lease_ttl     : %s\n", cfg.MaxLeaseTTL)
	fmt.Fprintf(n.w, "  force_no_cache    : %v\n", cfg.ForceNoCache)
	return nil
}
