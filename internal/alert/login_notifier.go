package alert

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/densestvoid/vaultwatch/internal/vault"
)

// LoginNotifier prints token login/lookup info to a writer.
type LoginNotifier struct {
	w io.Writer
}

// NewLoginNotifier constructs a LoginNotifier. If w is nil, os.Stdout is used.
func NewLoginNotifier(w io.Writer) *LoginNotifier {
	if w == nil {
		w = os.Stdout
	}
	return &LoginNotifier{w: w}
}

// Notify writes a summary of the LoginInfo to the writer.
func (n *LoginNotifier) Notify(info *vault.LoginInfo) error {
	if info == nil {
		_, err := fmt.Fprintln(n.w, "[login] no token info available")
		return err
	}

	renewable := "no"
	if info.Renewable {
		renewable = "yes"
	}

	expiry := info.IssuedAt.Add(time.Duration(info.LeaseDuration) * time.Second)
	lines := []string{
		"[login] token info:",
		fmt.Sprintf("  accessor  : %s", info.Accessor),
		fmt.Sprintf("  policies  : %s", strings.Join(info.Policies, ", ")),
		fmt.Sprintf("  ttl       : %ds", info.LeaseDuration),
		fmt.Sprintf("  renewable : %s", renewable),
		fmt.Sprintf("  issued at : %s", info.IssuedAt.UTC().Format(time.RFC3339)),
		fmt.Sprintf("  expires at: %s", expiry.UTC().Format(time.RFC3339)),
	}
	_, err := fmt.Fprintln(n.w, strings.Join(lines, "\n"))
	return err
}

// NotifyExpiring writes a warning to the writer when a token is close to expiry.
// threshold specifies how far in advance of expiry the warning should be emitted.
func (n *LoginNotifier) NotifyExpiring(info *vault.LoginInfo, threshold time.Duration) error {
	if info == nil {
		return nil
	}
	expiry := info.IssuedAt.Add(time.Duration(info.LeaseDuration) * time.Second)
	remaining := time.Until(expiry)
	if remaining > threshold {
		return nil
	}
	_, err := fmt.Fprintf(n.w, "[login] warning: token expires in %s (at %s)\n",
		remaining.Round(time.Second),
		expiry.UTC().Format(time.RFC3339),
	)
	return err
}
