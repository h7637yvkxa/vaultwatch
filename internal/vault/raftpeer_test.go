package vault_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yourusername/vaultwatch/internal/vault"
)

func newRaftPeerTestServer(t *testing.T, status int, peers []vault.RaftPeer) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/sys/storage/raft/configuration" {
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(status)
		if status == http.StatusOK {
			payload := map[string]interface{}{
				"data": map[string]interface{}{
					"config": map[string]interface{}{
						"servers": peers,
					},
				},
			}
			_ = json.NewEncoder(w).Encode(payload)
		}
	}))
}

func TestListRaftPeers_Success(t *testing.T) {
	peers := []vault.RaftPeer{
		{NodeID: "node1", Address: "127.0.0.1:8201", Leader: true, Voter: true},
		{NodeID: "node2", Address: "127.0.0.2:8201", Leader: false, Voter: true},
	}
	srv := newRaftPeerTestServer(t, http.StatusOK, peers)
	defer srv.Close()

	c := newVaultClient(t, srv.URL)
	checker := vault.NewRaftPeerChecker(c)

	result, err := checker.ListRaftPeers(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Peers) != 2 {
		t.Fatalf("expected 2 peers, got %d", len(result.Peers))
	}
	if result.Peers[0].NodeID != "node1" {
		t.Errorf("expected node1, got %s", result.Peers[0].NodeID)
	}
	if !result.Peers[0].Leader {
		t.Error("expected first peer to be leader")
	}
}

func TestListRaftPeers_Empty(t *testing.T) {
	srv := newRaftPeerTestServer(t, http.StatusOK, []vault.RaftPeer{})
	defer srv.Close()

	c := newVaultClient(t, srv.URL)
	checker := vault.NewRaftPeerChecker(c)

	result, err := checker.ListRaftPeers(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Peers) != 0 {
		t.Errorf("expected 0 peers, got %d", len(result.Peers))
	}
}

func TestListRaftPeers_HTTPError(t *testing.T) {
	srv := newRaftPeerTestServer(t, http.StatusInternalServerError, nil)
	defer srv.Close()

	c := newVaultClient(t, srv.URL)
	checker := vault.NewRaftPeerChecker(c)

	_, err := checker.ListRaftPeers(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestListRaftPeers_InvalidURL(t *testing.T) {
	c := newVaultClient(t, "http://127.0.0.1:0")
	checker := vault.NewRaftPeerChecker(c)

	_, err := checker.ListRaftPeers(context.Background())
	if err == nil {
		t.Fatal("expected error for unreachable server")
	}
}
