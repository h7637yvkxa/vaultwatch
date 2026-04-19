package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// KubernetesRole represents a Vault Kubernetes auth role.
type KubernetesRole struct {
	Name                 string   `json:"name"`
	BoundServiceAccounts []string `json:"bound_service_account_names"`
	BoundNamespaces      []string `json:"bound_service_account_namespaces"`
	TTL                  string   `json:"ttl"`
	Policies             []string `json:"policies"`
}

// KubernetesChecker lists Kubernetes auth roles from Vault.
type KubernetesChecker struct {
	address string
	token   string
	client  *http.Client
}

// NewKubernetesChecker creates a new KubernetesChecker.
func NewKubernetesChecker(address, token string) *KubernetesChecker {
	return &KubernetesChecker{
		address: address,
		token:   token,
		client:  &http.Client{},
	}
}

// ListKubernetesRoles returns all Kubernetes auth roles.
func (k *KubernetesChecker) ListKubernetesRoles(mountPath string) ([]KubernetesRole, error) {
	url := fmt.Sprintf("%s/v1/auth/%s/role?list=true", k.address, mountPath)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("kubernetes: build request: %w", err)
	}
	req.Header.Set("X-Vault-Token", k.token)

	resp, err := k.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("kubernetes: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return []KubernetesRole{}, nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("kubernetes: unexpected status %d", resp.StatusCode)
	}

	var envelope struct {
		Data struct {
			Keys []string `json:"keys"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
		return nil, fmt.Errorf("kubernetes: decode response: %w", err)
	}

	roles := make([]KubernetesRole, 0, len(envelope.Data.Keys))
	for _, name := range envelope.Data.Keys {
		roles = append(roles, KubernetesRole{Name: name})
	}
	return roles, nil
}
