package vault

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// AuthMethod represents a Vault auth method entry.
type AuthMethod struct {
	Path        string
	Type        string
	Description string
	Accessor    string
	Local       bool
	SealWrap    bool
	CreatedAt   time.Time
}

// AuthChecker lists enabled auth methods from Vault.
type AuthChecker struct {
	client *http.Client
	base   string
	token  string
}

// NewAuthChecker creates a new AuthChecker.
func NewAuthChecker(client *http.Client, base, token string) *AuthChecker {
	return &AuthChecker{client: client, base: base, token: token}
}

// ListAuthMethods returns all enabled auth methods.
func (a *AuthChecker) ListAuthMethods(ctx context.Context) ([]AuthMethod, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, a.base+"/v1/sys/auth", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Vault-Token", a.token)

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("auth list: unexpected status % struct {
		Data[string]struct{
			string ` `json:"description"`
			Accessor    string `json:"accessor"`
			Local       bool   `json:"local"`
			SealWrap    bool   `json:"seal_wrap"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}

	var methods []AuthMethod
	for path, m := range raw.Data {
		methods = append(methods, AuthMethod{
			Path:        path,
			Type:        m.Type,
			Description: m.Description,
			Accessor:    m.Accessor,
			Local:       m.Local,
			SealWrap:    m.SealWrap,
		})
	}
	return methods, nil
}
