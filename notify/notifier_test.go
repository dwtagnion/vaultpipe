package notify_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/yourusername/vaultpipe/notify"
)

func TestNew_EmptyURL(t *testing.T) {
	_, err := notify.New("")
	if err == nil {
		t.Fatal("expected error for empty URL, got nil")
	}
}

func TestNew_ValidURL(t *testing.T) {
	n, err := notify.New("http://example.com/hook")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n == nil {
		t.Fatal("expected non-nil Notifier")
	}
}

func TestSend_PostsJSON(t *testing.T) {
	var received notify.Event

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("expected application/json, got %s", ct)
		}
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n, _ := notify.New(ts.URL)
	ev := notify.Event{
		Operation: "sync",
		Role:      "backend",
		Keys:      []string{"DB_HOST", "DB_PASS"},
	}

	if err := n.Send(ev); err != nil {
		t.Fatalf("Send returned error: %v", err)
	}
	if received.Operation != "sync" {
		t.Errorf("expected operation sync, got %s", received.Operation)
	}
	if len(received.Keys) != 2 {
		t.Errorf("expected 2 keys, got %d", len(received.Keys))
	}
	if received.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp to be set")
	}
}

func TestSend_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	n, _ := notify.New(ts.URL)
	ev := notify.Event{Operation: "rotate", Timestamp: time.Now().UTC()}

	if err := n.Send(ev); err == nil {
		t.Fatal("expected error for non-2xx status, got nil")
	}
}

func TestSend_UnreachableHost(t *testing.T) {
	n, _ := notify.New("http://127.0.0.1:19999/hook")
	ev := notify.Event{Operation: "diff"}
	if err := n.Send(ev); err == nil {
		t.Fatal("expected error for unreachable host, got nil")
	}
}
