package alert

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestSlackNotifier_Notify_Empty(t *testing.T) {
	notifier := NewSlackNotifier("http://example.com/hook")
	if err := notifier.Notify(nil); err != nil {
		t.Fatalf("expected no error for empty alerts, got %v", err)
	}
}

func TestSlackNotifier_Notify_Success(t *testing.T) {
	var received slackPayload
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("failed to decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	notifier := NewSlackNotifier(server.URL)
	alerts := []Alert{
		{
			LeaseID:   "secret/db/prod",
			Level:     Warning,
			Message:   "expires in 2h",
			ExpiresAt: time.Now().Add(2 * time.Hour),
		},
		{
			LeaseID:   "secret/api/key",
			Level:     Critical,
			Message:   "expires in 20m",
			ExpiresAt: time.Now().Add(20 * time.Minute),
		},
	}

	if err := notifier.Notify(alerts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(received.Attachments) != 2 {
		t.Errorf("expected 2 attachments, got %d", len(received.Attachments))
	}
	if received.Attachments[0].Color != "warning" {
		t.Errorf("expected color 'warning', got %s", received.Attachments[0].Color)
	}
	if received.Attachments[1].Color != "danger" {
		t.Errorf("expected color 'danger', got %s", received.Attachments[1].Color)
	}
}

func TestSlackNotifier_Notify_HTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	notifier := NewSlackNotifier(server.URL)
	alerts := []Alert{{LeaseID: "secret/test", Level: Warning, Message: "test", ExpiresAt: time.Now().Add(time.Hour)}}

	if err := notifier.Notify(alerts); err == nil {
		t.Fatal("expected error for non-200 response, got nil")
	}
}

func TestSlackNotifier_Notify_InvalidURL(t *testing.T) {
	notifier := NewSlackNotifier("http://127.0.0.1:0/invalid")
	alerts := []Alert{{LeaseID: "secret/test", Level: Critical, Message: "test", ExpiresAt: time.Now().Add(time.Minute)}}

	if err := notifier.Notify(alerts); err == nil {
		t.Fatal("expected connection error, got nil")
	}
}
