package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// SlackNotifier sends alert notifications to a Slack webhook.
type SlackNotifier struct {
	WebhookURL string
	client     *http.Client
}

type slackPayload struct {
	Text        string            `json:"text"`
	Attachments []slackAttachment `json:"attachments,omitempty"`
}

type slackAttachment struct {
	Color  string `json:"color"`
	Title  string `json:"title"`
	Text   string `json:"text"`
	Footer string `json:"footer"`
}

// NewSlackNotifier creates a SlackNotifier with the given webhook URL.
func NewSlackNotifier(webhookURL string) *SlackNotifier {
	return &SlackNotifier{
		WebhookURL: webhookURL,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

// Notify sends all provided alerts to the configured Slack webhook.
func (s *SlackNotifier) Notify(alerts []Alert) error {
	if len(alerts) == 0 {
		return nil
	}

	attachments := make([]slackAttachment, 0, len(alerts))
	for _, a := range alerts {
		color := "warning"
		if a.Level == Critical {
			color = "danger"
		}
		attachments = append(attachments, slackAttachment{
			Color:  color,
			Title:  fmt.Sprintf("[%s] %s", a.Level, a.LeaseID),
			Text:   a.Message,
			Footer: fmt.Sprintf("Expires: %s", a.ExpiresAt.Format(time.RFC3339)),
		})
	}

	payload := slackPayload{
		Text:        fmt.Sprintf(":rotating_light: VaultWatch: %d secret(s) expiring soon", len(alerts)),
		Attachments: attachments,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("slack: marshal payload: %w", err)
	}

	resp, err := s.client.Post(s.WebhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("slack: post webhook: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("slack: unexpected status %d", resp.StatusCode)
	}
	return nil
}
