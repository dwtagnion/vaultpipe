package audit_test

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/yourusername/vaultpipe/audit"
)

func TestNewLogger_Stderr(t *testing.T) {
	l, err := audit.NewLogger("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if l == nil {
		t.Fatal("expected non-nil logger")
	}
}

func TestNewLogger_FileCreated(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "audit.log")

	l, err := audit.NewLogger(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if l == nil {
		t.Fatal("expected non-nil logger")
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatal("expected log file to be created")
	}
}

func TestLogSync_WritesJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "audit.log")

	l, err := audit.NewLogger(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	l.LogSync("secret/app", "DB_PASSWORD", "backend", "synced")

	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("open log: %v", err)
	}
	defer f.Close()

	var event audit.Event
	if err := json.NewDecoder(bufio.NewReader(f)).Decode(&event); err != nil {
		t.Fatalf("decode JSON: %v", err)
	}

	if event.Operation != "sync" {
		t.Errorf("operation: got %q, want %q", event.Operation, "sync")
	}
	if event.Key != "DB_PASSWORD" {
		t.Errorf("key: got %q, want %q", event.Key, "DB_PASSWORD")
	}
	if event.Role != "backend" {
		t.Errorf("role: got %q, want %q", event.Role, "backend")
	}
	if event.Status != "synced" {
		t.Errorf("status: got %q, want %q", event.Status, "synced")
	}
	if event.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestLogList_WritesJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "audit.log")

	l, _ := audit.NewLogger(path)
	l.LogList("secret/app", "ok", "found 3 keys")

	f, _ := os.Open(path)
	defer f.Close()

	var event audit.Event
	json.NewDecoder(f).Decode(&event) //nolint:errcheck

	if event.Operation != "list" {
		t.Errorf("operation: got %q, want %q", event.Operation, "list")
	}
	if event.Message != "found 3 keys" {
		t.Errorf("message: got %q, want %q", event.Message, "found 3 keys")
	}
}
