// Package sync orchestrates the end-to-end secret synchronisation
// workflow: read secrets from Vault, filter by role, and write to a
// local .env file.
package sync

import (
	"fmt"
	"log"

	"github.com/example/vaultpipe/config"
	"github.com/example/vaultpipe/env"
	"github.com/example/vaultpipe/filter"
	"github.com/example/vaultpipe/vault"
)

// Syncer coordinates the secret sync pipeline.
type Syncer struct {
	client *vault.Client
	filter *filter.Filter
	writer *env.Writer
	cfg    *config.Config
}

// New creates a Syncer from the provided configuration.
func New(cfg *config.Config) (*Syncer, error) {
	client, err := vault.NewClient(cfg.VaultAddr, cfg.VaultToken)
	if err != nil {
		return nil, fmt.Errorf("vault client: %w", err)
	}

	f, err := filter.New(cfg.RolesFile)
	if err != nil {
		return nil, fmt.Errorf("filter: %w", err)
	}

	w, err := env.NewWriter(cfg.OutputFile)
	if err != nil {
		return nil, fmt.Errorf("env writer: %w", err)
	}

	return &Syncer{client: client, filter: f, writer: w, cfg: cfg}, nil
}

// Run executes the sync: lists secrets at the configured path, applies
// role filtering, and writes matching key/value pairs to the .env file.
func (s *Syncer) Run(role string) error {
	keys, err := s.client.ListSecrets(s.cfg.SecretPath)
	if err != nil {
		return fmt.Errorf("list secrets at %q: %w", s.cfg.SecretPath, err)
	}

	log.Printf("found %d secret(s) under %s", len(keys), s.cfg.SecretPath)

	secrets := make(map[string]string, len(keys))
	for _, key := range keys {
		if !s.filter.Match(role, key) {
			log.Printf("skip %s (role=%s)", key, role)
			continue
		}

		data, err := s.client.ReadSecret(s.cfg.SecretPath + "/" + key)
		if err != nil {
			return fmt.Errorf("read secret %q: %w", key, err)
		}

		for k, v := range data {
			secrets[k] = v
		}
	}

	if err := s.writer.Write(secrets); err != nil {
		return fmt.Errorf("write env file: %w", err)
	}

	log.Printf("wrote %d key(s) to %s", len(secrets), s.cfg.OutputFile)
	return nil
}
