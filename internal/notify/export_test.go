// export_test.go exposes internal constructors for white-box testing.
package notify

// NewDispatcherFromNotifiers creates a Dispatcher directly from a slice of
// Notifier instances. Intended for use in tests only.
func NewDispatcherFromNotifiers(ns []Notifier) *Dispatcher {
	return &Dispatcher{notifiers: ns}
}
