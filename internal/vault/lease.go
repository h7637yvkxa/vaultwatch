package vault

import (
	"time"
)

// LeaseStatus holds the computed expiry state of a single lease.
type LeaseStatus struct {
	LeaseID    string
	Path       string
	TTL        time.Duration
	ExpiresAt  time.Time
	IsExpiring bool
	IsCritical bool
}

// CheckExpiry evaluates a lease against the provided warn and critical
// thresholds and returns a populated LeaseStatus.
func CheckExpiry(leaseID, path string, expiresAt time.Time, warnBefore, criticalBefore time.Duration) LeaseStatus {
	now := time.Now()
	ttl := expiresAt.Sub(now)

	status := LeaseStatus{
		LeaseID:   leaseID,
		Path:      path,
		TTL:       ttl,
		ExpiresAt: expiresAt,
	}

	if ttl <= criticalBefore {
		status.IsExpiring = true
		status.IsCritical = true
	} else if ttl <= warnBefore {
		status.IsExpiring = true
	}

	return status
}

// FilterExpiring returns only those statuses where IsExpiring is true.
func FilterExpiring(statuses []LeaseStatus) []LeaseStatus {
	var expiring []LeaseStatus
	for _, s := range statuses {
		if s.IsExpiring {
			expiring = append(expiring, s)
		}
	}
	return expiring
}
