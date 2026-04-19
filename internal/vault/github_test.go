package vault_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/eliziario/vaultwatch/internal/vault"
)

func newGitHubTestServer(t *testing.T, status int, keys []string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if status != http.StatusOK {
			w.WriteHeader(status)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{"keys": keys},
		})
	}))
}

func TestListGitHubTeams_Success(t *testing.T) {
	srv := newGitHubTestServer(t, http.StatusOK, []string{"dev", "ops"})
	defer srv.Close()

	c := vault.NewGitHubChecker(srv.URL, "tok", "")
	teams, err := c.ListTeams(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(teams) != 2 {
		t.Fatalf("expected 2 teams, got %d", len(teams))
	}
}

func TestListGitHubTeams_Empty(t *testing.T) {
	srv := newGitHubTestServer(t, http.StatusNotFound, nil)
	defer srv.Close()

	c := vault.NewGitHubChecker(srv.URL, "tok", "")
	teams, err := c.ListTeams(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(teams) != 0 {
		t.Fatalf("expected 0 teams, got %d", len(teams))
	}
}

func TestListGitHubTeams_HTTPError(t *testing.T) {
	srv := newGitHubTestServer(t, http.StatusInternalServerError, nil)
	defer srv.Close()

	c := vault.NewGitHubChecker(srv.URL, "tok", "")
	_, err := c.ListTeams(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestListGitHubTeams_InvalidURL(t *testing.T) {
	c := vault.NewGitHubChecker("http://127.0.0.1:0", "tok", "")
	_, err := c.ListTeams(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
}
