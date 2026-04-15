package vault

import (
	"context"
	"fmt"
	"time"

	vaultapi "github.com/hashicorp/vault/api"
)

// RotateResult holds the outcome of a single secret rotation attempt.
type RotateResult struct {
	LeaseID   string
	NewLeaseID string
	RenewedAt time.Time
	Err       error
}

// Rotator handles forced rotation of Vault leases by revoking and re-reading.
type Rotator struct {
	client *vaultapi.Client
}

// NewRotator creates a Rotator backed by the given Vault client.
func NewRotator(client *vaultapi.Client) *Rotator {
	return &Rotator{client: client}
}

// RotateLease revokes the given lease and re-reads the secret path to obtain a
// fresh lease. It returns the new lease ID on success.
func (r *Rotator) RotateLease(ctx context.Context, leaseID, secretPath string) RotateResult {
	if leaseID == "" {
		return RotateResult{Err: fmt.Errorf("leaseID must not be empty")}
	}
	if secretPath == "" {
		return RotateResult{Err: fmt.Errorf("secretPath must not be empty")}
	}

	// Revoke the existing lease.
	if err := r.client.Sys().RevokeWithContext(ctx, leaseID); err != nil {
		return RotateResult{LeaseID: leaseID, Err: fmt.Errorf("revoke lease %s: %w", leaseID, err)}
	}

	// Re-read the secret to obtain a fresh lease.
	secret, err := r.client.Logical().ReadWithContext(ctx, secretPath)
	if err != nil {
		return RotateResult{LeaseID: leaseID, Err: fmt.Errorf("re-read secret %s: %w", secretPath, err)}
	}
	if secret == nil {
		return RotateResult{LeaseID: leaseID, Err: fmt.Errorf("no secret returned for path %s", secretPath)}
	}

	return RotateResult{
		LeaseID:    leaseID,
		NewLeaseID: secret.LeaseID,
		RenewedAt:  time.Now(),
	}
}

// RotateExpiring rotates all leases whose status is Warning or Critical.
func (r *Rotator) RotateExpiring(ctx context.Context, statuses []LeaseStatus, pathMap map[string]string) []RotateResult {
	var results []RotateResult
	for _, s := range statuses {
		if s.Level == LevelOK {
			continue
		}
		path, ok := pathMap[s.LeaseID]
		if !ok {
			results = append(results, RotateResult{
				LeaseID: s.LeaseID,
				Err:     fmt.Errorf("no secret path mapped for lease %s", s.LeaseID),
			})
			continue
		}
		results = append(results, r.RotateLease(ctx, s.LeaseID, path))
	}
	return results
}
