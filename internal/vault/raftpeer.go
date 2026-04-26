package vault

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// RaftPeer represents a single member of the Raft cluster.
type RaftPeer struct {
	NodeID          string `json:"node_id"`
	Address         string `json:"address"`
	Leader          bool   `json:"leader"`
	ProtocolVersion string `json:"protocol_version"`
	Voter           bool   `json:"voter"`
}

// RaftPeerResult holds the list of Raft peers returned by Vault.
type RaftPeerResult struct {
	Peers []RaftPeer
}

// RaftPeerChecker fetches Raft cluster configuration from Vault.
type RaftPeerChecker struct {
	client *http.Client
	base   string
	token  string
}

// NewRaftPeerChecker constructs a RaftPeerChecker using the provided Vault client.
func NewRaftPeerChecker(c *Client) *RaftPeerChecker {
	return &RaftPeerChecker{
		client: c.HTTP,
		base:   c.Address,
		token:  c.Token,
	}
}

// ListRaftPeers retrieves the current Raft peer list from Vault.
func (r *RaftPeerChecker) ListRaftPeers(ctx context.Context) (*RaftPeerResult, error) {
	url := fmt.Sprintf("%s/v1/sys/storage/raft/configuration", r.base)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("raft peer: build request: %w", err)
	}
	req.Header.Set("X-Vault-Token", r.token)

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("raft peer: request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("raft peer: unexpected status %d", resp.StatusCode)
	}

	var payload struct {
		Data struct {
			Config struct {
				Servers []RaftPeer `json:"servers"`
			} `json:"config"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("raft peer: decode: %w", err)
	}

	return &RaftPeerResult{Peers: payload.Data.Config.Servers}, nil
}
