package vault

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newPolicyTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/v1/sys/policies/acl" && r.URL.RawQuery == "list=true":
			json.NewEncoder(w).Encode(map[string]any{
				"data": map[string]any{"keys": []string{"default", "root", "app-policy"}},
			})
		case r.URL.Path == "/v1/sys/policies/acl/app-policy":
			json.NewEncoder(w).Encode(map[string]any{
				"data": map[string]any{"policy": `path "secret/*" { capabilities = ["read"] }`},
			})
		case r.URL.Path == "/v1/sys/policies/acl/missing":
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"errors":["policy not found"]}`))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
}

func TestListPolicies_Success(t *testing.T) {
	srv := newPolicyTestServer(t)
	defer srv.Close()

	c := newVaultClient(t, srv.URL)
	pc := NewPolicyChecker(c)

	policies, err := pc.ListPolicies(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(policies) != 3 {
		t.Fatalf("expected 3 policies, got %d", len(policies))
	}
}

func TestGetPolicy_Success(t *testing.T) {
	srv := newPolicyTestServer(t)
	defer srv.Close()

	c := newVaultClient(t, srv.URL)
	pc := NewPolicyChecker(c)

	info, err := pc.GetPolicy(context.Background(), "app-policy")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info.Name != "app-policy" {
		t.Errorf("expected name app-policy, got %s", info.Name)
	}
	if info.Rules == "" {
		t.Error("expected non-empty rules")
	}
}

func TestGetPolicy_NotFound(t *testing.T) {
	srv := newPolicyTestServer(t)
	defer srv.Close()

	c := newVaultClient(t, srv.URL)
	pc := NewPolicyChecker(c)

	_, err := pc.GetPolicy(context.Background(), "missing")
	if err == nil {
		t.Fatal("expected error for missing policy")
	}
}
