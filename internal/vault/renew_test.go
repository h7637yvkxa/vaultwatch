package vault_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	vaultapi "github.com/hashicorp/vault/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/your-org/vaultwatch/internal/vault"
)

func newTestVaultServer(t *testing.T, renewDuration int) (*httptest.Server, *vaultapi.Client) {
	t.Helper()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut && r.URL.Path == "/v1/sys/leases/renew" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{ //nolint:errcheck
				"lease_id":       "secret/data/test/abc123",
				"lease_duration": renewDuration,
				"renewable":      true,
			})
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))

	cfg := vaultapi.DefaultConfig()
	cfg.Address = ts.URL
	client, err := vaultapi.NewClient(cfg)
	require.NoError(t, err)
	client.SetToken("test-token")
	return ts, client
}

func TestRenewLease_Success(t *testing.T) {
	ts, apiClient := newTestVaultServer(t, 3600)
	defer ts.Close()

	renewer := vault.NewRenewer(apiClient)
	result := renewer.RenewLease(context.Background(), "secret/data/test/abc123", 3600)

	assert.NoError(t, result.Error)
	assert.True(t, result.Renewed)
	assert.Equal(t, time.Hour, result.NewTTL)
	assert.Equal(t, "secret/data/test/abc123", result.LeaseID)
}

func TestRenewLease_EmptyID(t *testing.T) {
	ts, apiClient := newTestVaultServer(t, 3600)
	defer ts.Close()

	renewer := vault.NewRenewer(apiClient)
	result := renewer.RenewLease(context.Background(), "", 3600)

	assert.Error(t, result.Error)
	assert.False(t, result.Renewed)
}

func TestRenewExpiring_SkipsOK(t *testing.T) {
	ts, apiClient := newTestVaultServer(t, 3600)
	defer ts.Close()

	renewer := vault.NewRenewer(apiClient)
	statuses := []vault.LeaseStatus{
		{LeaseID: "ok-lease", Level: vault.LevelOK, TTL: 24 * time.Hour},
	}
	results := renewer.RenewExpiring(context.Background(), statuses, 3600)
	assert.Empty(t, results)
}

func TestRenewExpiring_RenewsWarningAndCritical(t *testing.T) {
	ts, apiClient := newTestVaultServer(t, 3600)
	defer ts.Close()

	renewer := vault.NewRenewer(apiClient)
	statuses := []vault.LeaseStatus{
		{LeaseID: "warn-lease", Level: vault.LevelWarning, TTL: 30 * time.Minute},
		{LeaseID: "crit-lease", Level: vault.LevelCritical, TTL: 5 * time.Minute},
		{LeaseID: "ok-lease", Level: vault.LevelOK, TTL: 24 * time.Hour},
	}
	results := renewer.RenewExpiring(context.Background(), statuses, 3600)
	assert.Len(t, results, 2)
}
