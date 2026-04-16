package alert

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"

	"github.com/your-org/vaultwatch/internal/vault"
)

// MountNotifier prints mount information to a writer.
type MountNotifier struct {
	w io.Writer
}

// NewMountNotifier creates a MountNotifier. If w is nil, os.Stdout is used.
func NewMountNotifier(w io.Writer) *MountNotifier {
	if w == nil {
		w = os.Stdout
	}
	return &MountNotifier{w: w}
}

// Notify prints a summary table of all mounted secret engines.
func (n *MountNotifier) Notify(mounts []vault.MountEntry) error {
	if len(mounts) == 0 {
		fmt.Fprintln(n.w, "[mounts] no secret engines found")
		return nil
	}

	tw := tabwriter.NewWriter(n.w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "PATH\tTYPE\tDESCRIPTION\tACCESSOR")
	for _, m := range mounts {
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", m.Path, m.Type, m.Description, m.Accessor)
	}
	return tw.Flush()
}
