// Package notify provides post-sync notification support,
// alerting configured channels when secrets are synced or rotated.
package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Event describes a sync or rotate event to be reported.
type Event struct {
	Timestamp time.Time `json:"timestamp"`
	Operation string    `json:"operation"` // "sync" | "rotate" | "diff"
	Role      string    `json:"role"`
	Keys      []string  `json:"keys"`
	Error     string    `json:"error,omitempty"`
}

// Notifier sends Event payloads to a webhook URL.
type Notifier struct {
	webhookURL string
	client     *http.Client
}

// New creates a Notifier targeting the given webhook URL.
// Returns an error if webhookURL is empty.
func New(webhookURL string) (*Notifier, error) {
	if webhookURL == "" {
		return nil, fmt.Errorf("notify: webhook URL must not be empty")
	}
	return &Notifier{
		webhookURL: webhookURL,
		client:     &http.Client{Timeout: 10 * time.Second},
	}, nil
}

// Send marshals the Event as JSON and POSTs it to the configured webhook.
func (n *Notifier) Send(ev Event) error {
	if ev.Timestamp.IsZero() {
		ev.Timestamp = time.Now().UTC()
	}

	body, err := json.Marshal(ev)
	if err != nil {
		return fmt.Errorf("notify: marshal event: %w", err)
	}

	resp, err := n.client.Post(n.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("notify: post to webhook: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("notify: webhook returned status %d", resp.StatusCode)
	}
	return nil
}
