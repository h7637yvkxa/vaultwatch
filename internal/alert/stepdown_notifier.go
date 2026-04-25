package alert

import (
	"fmt"
	"io"
	"os"

	"github.com/wormhole-enterprise/vaultwatch/internal/vault"
)

// StepDownNotifier formats and writes the result of a Vault step-down request.
type StepDownNotifier struct {
	w io.Writer
}

// NewStepDownNotifier creates a StepDownNotifier that writes to w.
// If w is nil, os.Stdout is used.
func NewStepDownNotifier(w io.Writer) *StepDownNotifier {
	if w == nil {
		w = os.Stdout
	}
	return &StepDownNotifier{w: w}
}

// Notify prints the step-down result to the configured writer.
func (n *StepDownNotifier) Notify(result *vault.StepDownResult) error {
	if result == nil {
		fmt.Fprintln(n.w, "[step-down] no result available")
		return nil
	}

	status := "SUCCESS"
	if !result.Success {
		status = "FAILED"
	}

	fmt.Fprintf(n.w, "[step-down] %s: %s\n", status, result.Message)
	return nil
}
