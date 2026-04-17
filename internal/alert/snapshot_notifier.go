package alert

import (
	"fmt"
	"io"
	"os"
	"time"
)

// SnapshotResult holds the outcome of a snapshot operation.
type SnapshotResult struct {
	TakenAt time.Time
	Size    int64
	Err     error
}

// SnapshotNotifier writes snapshot results to a writer.
type SnapshotNotifier struct {
	w io.Writer
}

// NewSnapshotNotifier returns a SnapshotNotifier writing to w.
// If w is nil, os.Stdout is used.
func NewSnapshotNotifier(w io.Writer) *SnapshotNotifier {
	if w == nil {
		w = os.Stdout
	}
	return &SnapshotNotifier{w: w}
}

// Notify writes a summary of the snapshot result.
func (n *SnapshotNotifier) Notify(result *SnapshotResult) error {
	if result == nil {
		fmt.Fprintln(n.w, "[snapshot] no result available")
		return nil
	}
	if result.Err != nil {
		fmt.Fprintf(n.w, "[snapshot] FAILED at %s: %v\n",
			result.TakenAt.Format(time.RFC3339), result.Err)
		return nil
	}
	fmt.Fprintf(n.w, "[snapshot] OK — taken at %s, size %d bytes\n",
		result.TakenAt.Format(time.RFC3339), result.Size)
	return nil
}
