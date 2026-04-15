package vault

import (
	"testing"
	"time"
)

func makeLeaseWithTTL(ttl time.Duration) LeaseInfo {
	return LeaseInfo{
		LeaseID:  "test/lease/1",
		TTL:      ttl,
		ExpireAt: time.Now().Add(ttl),
	}
}

func TestCheckExpiry_OK(t *testing.T) {
	lease := makeLeaseWithTTL(48 * time.Hour)
	status := CheckExpiry(lease, 24*time.Hour)
	if status != StatusOK {
		t.Errorf("expected StatusOK, got %d", status)
	}
}

func TestCheckExpiry_Warning(t *testing.T) {
	lease := makeLeaseWithTTL(12 * time.Hour)
	status := CheckExpiry(lease, 24*time.Hour)
	if status != StatusWarning {
		t.Errorf("expected StatusWarning, got %d", status)
	}
}

func TestCheckExpiry_Critical(t *testing.T) {
	lease := makeLeaseWithTTL(30 * time.Second)
	status := CheckExpiry(lease, 24*time.Hour)
	if status != StatusCritical {
		t.Errorf("expected StatusCritical, got %d", status)
	}
}

func TestFilterExpiring_MixedLeases(t *testing.T) {
	warnBefore := 24 * time.Hour
	leases := []LeaseInfo{
		makeLeaseWithTTL(48 * time.Hour),  // OK
		makeLeaseWithTTL(12 * time.Hour),  // Warning
		makeLeaseWithTTL(30 * time.Second), // Critical
		makeLeaseWithTTL(72 * time.Hour),  // OK
	}

	expiring := FilterExpiring(leases, warnBefore)
	if len(expiring) != 2 {
		t.Errorf("expected 2 expiring leases, got %d", len(expiring))
	}
}

func TestFilterExpiring_NoneExpiring(t *testing.T) {
	warnBefore := 1 * time.Hour
	leases := []LeaseInfo{
		makeLeaseWithTTL(48 * time.Hour),
		makeLeaseWithTTL(24 * time.Hour),
	}

	expiring := FilterExpiring(leases, warnBefore)
	if len(expiring) != 0 {
		t.Errorf("expected 0 expiring leases, got %d", len(expiring))
	}
}

func TestFilterExpiring_AllExpiring(t *testing.T) {
	warnBefore := 72 * time.Hour
	leases := []LeaseInfo{
		makeLeaseWithTTL(1 * time.Hour),
		makeLeaseWithTTL(30 * time.Second),
	}

	expiring := FilterExpiring(leases, warnBefore)
	if len(expiring) != 2 {
		t.Errorf("expected 2 expiring leases, got %d", len(expiring))
	}
}
