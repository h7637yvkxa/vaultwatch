package alert

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// PolicyReport summarises the result of a policy listing operation.
type PolicyReport struct {
	Policies []string
	Error    error
}

// PolicyNotifier writes a policy report to a writer.
type PolicyNotifier struct {
	w io.Writer
}

// NewPolicyNotifier creates a PolicyNotifier. If w is nil, os.Stdout is used.
func NewPolicyNotifier(w io.Writer) *PolicyNotifier {
	if w == nil {
		w = os.Stdout
	}
	return &PolicyNotifier{w: w}
}

// Notify writes the policy report to the configured writer.
func (n *PolicyNotifier) Notify(report PolicyReport) error {
	if report.Error != nil {
		_, err := fmt.Fprintf(n.w, "[policy] error fetching policies: %v\n", report.Error)
		return err
	}
	if len(report.Policies) == 0 {
		_, err := fmt.Fprintln(n.w, "[policy] no policies found")
		return err
	}
	fmt.Fprintf(n.w, "[policy] %d policies found:\n", len(report.Policies))
	for _, p := range report.Policies {
		if _, err := fmt.Fprintf(n.w, "  - %s\n", p); err != nil {
			return err
		}
	}
	return nil
}

// Summary returns a single-line summary string of the report.
func (n *PolicyNotifier) Summary(report PolicyReport) string {
	if report.Error != nil {
		return fmt.Sprintf("policy fetch error: %v", report.Error)
	}
	return fmt.Sprintf("%d policies: %s", len(report.Policies), strings.Join(report.Policies, ", "))
}
