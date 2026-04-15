package vault

import (
	"context"
	"fmt"
	"time"

	vaultapi "github.com/hashicorp/vault/api"
)

// RenewResult holds the outcome of a lease renewal attempt.
type RenewResult struct {
	LeaseID    string
	NewTTL     time.Duration
	Renewed    bool
	Error      error
}

// Renewer wraps a Vault client and provides lease renewal capabilities.
type Renewer struct {
	client *vaultapi.Client
}

// NewRenewer creates a Renewer from an existing Vault API client.
func NewRenewer(client *vaultapi.Client) *Renewer {
	return &Renewer{client: client}
}

// RenewLease attempts to renew a single lease by its ID, requesting the given
// increment in seconds. It returns a RenewResult describing the outcome.
func (r *Renewer) RenewLease(ctx context.Context, leaseID string, incrementSeconds int) RenewResult {
	if leaseID == "" {
		return RenewResult{
			LeaseID: leaseID,
			Error:   fmt.Errorf("lease ID must not be empty"),
		}
	}

	secret, err := r.client.Sys().RenewWithContext(ctx, leaseID, incrementSeconds)
	if err != nil {
		return RenewResult{
			LeaseID: leaseID,
			Error:   fmt.Errorf("renewing lease %q: %w", leaseID, err),
		}
	}

	return RenewResult{
		LeaseID: leaseID,
		NewTTL:  time.Duration(secret.LeaseDuration) * time.Second,
		Renewed: true,
	}
}

// RenewExpiring iterates over the provided LeaseStatus slice and renews any
// lease whose level is Warning or Critical. It returns all results.
func (r *Renewer) RenewExpiring(ctx context.Context, statuses []LeaseStatus, incrementSeconds int) []RenewResult {
	var results []RenewResult
	for _, s := range statuses {
		if s.Level == LevelOK {
			continue
		}
		result := r.RenewLease(ctx, s.LeaseID, incrementSeconds)
		results = append(results, result)
	}
	return results
}
