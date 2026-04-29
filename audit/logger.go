// Package audit provides structured audit logging for secret sync operations.
package audit

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

// Event represents a single audit log entry.
type Event struct {
	Timestamp time.Time `json:"timestamp"`
	Operation string    `json:"operation"`
	Path      string    `json:"path"`
	Key       string    `json:"key,omitempty"`
	Role      string    `json:"role,omitempty"`
	Status    string    `json:"status"`
	Message   string    `json:"message,omitempty"`
}

// Logger writes audit events as newline-delimited JSON.
type Logger struct {
	w io.Writer
}

// NewLogger creates a Logger writing to the given path.
// Pass an empty path to write to stderr.
func NewLogger(path string) (*Logger, error) {
	if path == "" {
		return &Logger{w: os.Stderr}, nil
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return nil, fmt.Errorf("audit: open log file: %w", err)
	}
	return &Logger{w: f}, nil
}

// LogSync records a secret key sync event.
func (l *Logger) LogSync(path, key, role, status string) {
	l.write(Event{
		Timestamp: time.Now().UTC(),
		Operation: "sync",
		Path:      path,
		Key:       key,
		Role:      role,
		Status:    status,
	})
}

// LogList records a secret listing event.
func (l *Logger) LogList(path, status, message string) {
	l.write(Event{
		Timestamp: time.Now().UTC(),
		Operation: "list",
		Path:      path,
		Status:    status,
		Message:   message,
	})
}

func (l *Logger) write(e Event) {
	data, err := json.Marshal(e)
	if err != nil {
		return
	}
	_, _ = fmt.Fprintf(l.w, "%s\n", data)
}
