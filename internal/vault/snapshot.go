package vault

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

// SnapshotMeta holds metadata about a Vault raft snapshot.
type SnapshotMeta struct {
	TakenAt time.Time
	Size    int64
}

// SnapshotChecker retrieves raft snapshot info from Vault.
type SnapshotChecker struct {
	address string
	token   string
	client  *http.Client
}

// NewSnapshotChecker creates a new SnapshotChecker.
func NewSnapshotChecker(address, token string) *SnapshotChecker {
	return &SnapshotChecker{
		address: address,
		token:   token,
		client:  &http.Client{Timeout: 15 * time.Second},
	}
}

// TakeSnapshot requests a snapshot from Vault and writes it to w.
// It returns metadata about the snapshot.
func (s *SnapshotChecker) TakeSnapshot(ctx context.Context, w io.Writer) (*SnapshotMeta, error) {
	url := fmt.Sprintf("%s/v1/sys/storage/raft/snapshot", s.address)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("snapshot: build request: %w", err)
	}
	req.Header.Set("X-Vault-Token", s.token)

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("snapshot: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("snapshot: unexpected status %d", resp.StatusCode)
	}

	n, err := io.Copy(w, resp.Body)
	if err != nil {
		return nil, fmt.Errorf("snapshot: write body: %w", err)
	}

	return &SnapshotMeta{
		TakenAt: time.Now().UTC(),
		Size:    n,
	}, nil
}
