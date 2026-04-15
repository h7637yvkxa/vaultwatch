package vault

import "time"

// ExpiryStatus describes how urgently a lease needs attention.
type ExpiryStatus int

const (
	StatusOK      ExpiryStatus = iota // plenty of time remaining
	StatusWarning                     // within the warn-before window
	StatusCritical                    // expired or about to expire (<1 min)
)

// CheckExpiry evaluates a LeaseInfo against the configured warning threshold
// and returns the appropriate ExpiryStatus.
func CheckExpiry(lease LeaseInfo, warnBefore time.Duration) ExpiryStatus {
	switch {
	case lease.TTL <= time.Minute:
		return StatusCritical
	case lease.TTL <= warnBefore:
		return StatusWarning
	default:
		return StatusOK
	}
}

// FilterExpiring returns only those leases whose status is Warning or Critical.
func FilterExpiring(leases []LeaseInfo, warnBefore time.Duration) []LeaseInfo {
	var expiring []LeaseInfo
	for _, l := range leases {
		if CheckExpiry(l, warnBefore) != StatusOK {
			expiring = append(expiring, l)
		}
	}
	return expiring
}
