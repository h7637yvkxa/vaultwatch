package vault

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type GitHubTeam struct {
	Name   string
	Policy string
}

type GitHubChecker struct {
	address string
	token   string
	client  *http.Client
	mount   string
}

func NewGitHubChecker(address, token, mount string) *GitHubChecker {
	if mount == "" {
		mount = "github"
	}
	return &GitHubChecker{address: address, token: token, client: &http.Client{}, mount: mount}
}

func (g *GitHubChecker) ListTeams(ctx context.Context) ([]GitHubTeam, error) {
	url := fmt.Sprintf("%s/v1/auth/%s/map/teams", g.address, g.mount)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Vault-Token", g.token)
	req.Header.Set("X-Vault-Request", "true")

	resp, err := g.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("vault returned status %d", resp.StatusCode)
	}

	var result struct {
		Data struct {
			Keys []string `json:"keys"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	teams := make([]GitHubTeam, 0, len(result.Data.Keys))
	for _, k := range result.Data.Keys {
		teams = append(teams, GitHubTeam{Name: k, Policy: "mapped"})
	}
	return teams, nil
}
