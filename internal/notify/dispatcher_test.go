package notify_test

import (
	"context"
	"errors"
	"testing"

	"github.com/yourusername/vaultwatch/internal/alert"
	"github.com/yourusername/vaultwatch/internal/notify"
)

// stubNotifier is a test double that records calls and can return an error.
type stubNotifier struct {
	called bool
	err    error
}

func (s *stubNotifier) Notify(_ context.Context, _ []alert.Alert) error {
	s.called = true
	return s.err
}

func TestDispatcher_Dispatch_NoAlerts(t *testing.T) {
	stub := &stubNotifier{}
	d := notify.NewDispatcherFromNotifiers([]notify.Notifier{stub})

	if err := d.Dispatch(context.Background(), nil); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if stub.called {
		t.Error("notifier should not be called when there are no alerts")
	}
}

func TestDispatcher_Dispatch_CallsNotifiers(t *testing.T) {
	stub := &stubNotifier{}
	d := notify.NewDispatcherFromNotifiers([]notify.Notifier{stub})

	alerts := []alert.Alert{{Path: "secret/db", Level: alert.Warning}}
	if err := d.Dispatch(context.Background(), alerts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !stub.called {
		t.Error("expected notifier to be called")
	}
}

func TestDispatcher_Dispatch_CollectsErrors(t *testing.T) {
	stub1 := &stubNotifier{err: errors.New("channel down")}
	stub2 := &stubNotifier{}
	d := notify.NewDispatcherFromNotifiers([]notify.Notifier{stub1, stub2})

	alerts := []alert.Alert{{Path: "secret/api", Level: alert.Critical}}
	err := d.Dispatch(context.Background(), alerts)
	if err == nil {
		t.Fatal("expected error from failing notifier")
	}
	if !stub2.called {
		t.Error("second notifier should still be called despite first error")
	}
}

func TestDispatcher_Dispatch_NoNotifiers(t *testing.T) {
	// Dispatching with an empty notifier list should succeed silently.
	d := notify.NewDispatcherFromNotifiers([]notify.Notifier{})

	alerts := []alert.Alert{{Path: "secret/db", Level: alert.Warning}}
	if err := d.Dispatch(context.Background(), alerts); err != nil {
		t.Fatalf("expected no error with no notifiers, got %v", err)
	}
}
