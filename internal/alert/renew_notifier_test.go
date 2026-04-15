package alert_test

import (
	"bytes"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/your-org/vaultwatch/internal/alert"
	"github.com/your-org/vaultwatch/internal/vault"
)

func TestRenewNotifier_NoResults(t *testing.T) {
	var buf bytes.Buffer
	n := alert.NewRenewNotifier(&buf)
	err := n.Notify(context.Background(), nil)
	require.NoError(t, err)
	assert.Contains(t, buf.String(), "no leases required renewal")
}

func TestRenewNotifier_AllSuccess(t *testing.T) {
	var buf bytes.Buffer
	n := alert.NewRenewNotifier(&buf)
	results := []vault.RenewResult{
		{LeaseID: "secret/db/creds/abc", NewTTL: time.Hour, Renewed: true},
		{LeaseID: "secret/db/creds/xyz", NewTTL: 30 * time.Minute, Renewed: true},
	}
	err := n.Notify(context.Background(), results)
	require.NoError(t, err)
	out := buf.String()
	assert.Contains(t, out, "OK")
	assert.Contains(t, out, "secret/db/creds/abc")
	assert.Contains(t, out, "1h0m0s")
}

func TestRenewNotifier_WithErrors(t *testing.T) {
	var buf bytes.Buffer
	n := alert.NewRenewNotifier(&buf)
	results := []vault.RenewResult{
		{LeaseID: "good-lease", NewTTL: time.Hour, Renewed: true},
		{LeaseID: "bad-lease", Error: errors.New("permission denied")},
	}
	err := n.Notify(context.Background(), results)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "permission denied")
	out := buf.String()
	assert.Contains(t, out, "FAILED")
	assert.Contains(t, out, "bad-lease")
	assert.Contains(t, out, "OK")
	assert.Contains(t, out, "good-lease")
}

func TestRenewNotifier_NilWriter_UsesStdout(t *testing.T) {
	// Should not panic when w is nil; uses os.Stdout.
	n := alert.NewRenewNotifier(nil)
	assert.NotNil(t, n)
}
