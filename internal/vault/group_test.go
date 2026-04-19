package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newGroupTestServer(t *testing.T, status int, payload any) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		if payload != nil {
			_ = json.NewEncoder(w).Encode(payload)
		}
	}))
}

func TestListGroups_Success(t *testing.T) {
	payload := map[string]any{
		"data": map[string]any{
			"key_info": map[string]any{
				"abc123": map[string]any{
					"name": "admins",
					"type": "internal",
					"policies": []string{"default"},
				},
			},
		},
	}
	srv := newGroupTestServer(t, http.StatusOK, payload)
	defer srv.Close()

	c := newVaultClient(t, srv.URL)
	groups, err := NewGroupChecker(c).ListGroups()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(groups) != 1 {
		t.Fatalf("expected 1 group, got %d", len(groups))
	}
	if groups[0].Name != "admins" {
		t.Errorf("expected name admins, got %s", groups[0].Name)
	}
}

func TestListGroups_Empty(t *testing.T) {
	srv := newGroupTestServer(t, http.StatusNotFound, nil)
	defer srv.Close()

	c := newVaultClient(t, srv.URL)
	groups, err := NewGroupChecker(c).ListGroups()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(groups) != 0 {
		t.Errorf("expected 0 groups, got %d", len(groups))
	}
}

func TestListGroups_HTTPError(t *testing.T) {
	srv := newGroupTestServer(t, http.StatusInternalServerError, nil)
	defer srv.Close()

	c := newVaultClient(t, srv.URL)
	_, err := NewGroupChecker(c).ListGroups()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
