package vault

import "net/http"

func NewGitHubCheckerFromParts(address, token, mount string, client *http.Client) *GitHubChecker {
	c := NewGitHubChecker(address, token, mount)
	c.client = client
	return c
}
