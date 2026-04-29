// Package watch provides functionality to watch Vault secret paths
// for changes and trigger re-sync when secrets are updated.
package watch

import (
	"context"
	"log"
	"time"

	"github.com/yourusername/vaultpipe/sync"
	"github.com/yourusername/vaultpipe/vault"
)

// Watcher polls Vault secret paths at a configurable interval and
// triggers a sync whenever a change is detected.
type Watcher struct {
	client   *vault.Client
	syncer   *sync.Syncer
	paths    []string
	interval time.Duration
	logger   *log.Logger
}

// Config holds configuration for a Watcher.
type Config struct {
	Paths    []string
	Interval time.Duration
	Logger   *log.Logger
}

// New creates a new Watcher with the given Vault client, syncer, and config.
func New(client *vault.Client, syncer *sync.Syncer, cfg Config) *Watcher {
	interval := cfg.Interval
	if interval <= 0 {
		interval = 30 * time.Second
	}
	logger := cfg.Logger
	if logger == nil {
		logger = log.Default()
	}
	return &Watcher{
		client:   client,
		syncer:   syncer,
		paths:    cfg.Paths,
		interval: interval,
		logger:   logger,
	}
}

// Run starts the watch loop, polling at the configured interval until ctx is cancelled.
func (w *Watcher) Run(ctx context.Context) error {
	w.logger.Printf("watch: starting with interval %s on %d path(s)", w.interval, len(w.paths))
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	// Perform an initial sync immediately.
	if err := w.syncOnce(ctx); err != nil {
		w.logger.Printf("watch: initial sync error: %v", err)
	}

	for {
		select {
		case <-ctx.Done():
			w.logger.Println("watch: context cancelled, stopping")
			return ctx.Err()
		case <-ticker.C:
			if err := w.syncOnce(ctx); err != nil {
				w.logger.Printf("watch: sync error: %v", err)
			}
		}
	}
}

func (w *Watcher) syncOnce(ctx context.Context) error {
	w.logger.Println("watch: triggering sync")
	return w.syncer.Run(ctx)
}
