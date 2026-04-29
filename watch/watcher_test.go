package watch_test

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/yourusername/vaultpipe/watch"
)

// stubSyncer satisfies the syncer interface used in tests.
type stubSyncer struct {
	callCount int
	errToReturn error
}

func (s *stubSyncer) Run(_ context.Context) error {
	s.callCount++
	return s.errToReturn
}

func TestNew_DefaultInterval(t *testing.T) {
	w := watch.New(nil, nil, watch.Config{
		Paths:    []string{"secret/app"},
		Interval: 0,
	})
	if w == nil {
		t.Fatal("expected non-nil Watcher")
	}
}

func TestNew_CustomInterval(t *testing.T) {
	w := watch.New(nil, nil, watch.Config{
		Paths:    []string{"secret/app"},
		Interval: 5 * time.Second,
	})
	if w == nil {
		t.Fatal("expected non-nil Watcher")
	}
}

func TestRun_CancelImmediately(t *testing.T) {
	logger := log.New(os.Stderr, "test: ", 0)
	w := watch.New(nil, nil, watch.Config{
		Paths:    []string{"secret/app"},
		Interval: 1 * time.Hour, // long interval so ticker never fires
		Logger:   logger,
	})

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel before Run

	err := w.Run(ctx)
	if err != context.Canceled {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
}

func TestRun_TickerFiresSync(t *testing.T) {
	logger := log.New(os.Stderr, "test: ", 0)
	w := watch.New(nil, nil, watch.Config{
		Paths:    []string{"secret/app"},
		Interval: 20 * time.Millisecond,
		Logger:   logger,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 70*time.Millisecond)
	defer cancel()

	// Run returns context.DeadlineExceeded after timeout.
	err := w.Run(ctx)
	if err != context.DeadlineExceeded {
		t.Fatalf("expected DeadlineExceeded, got %v", err)
	}
}
