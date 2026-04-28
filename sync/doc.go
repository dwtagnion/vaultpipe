// Package sync provides the high-level orchestration layer for vaultpipe.
//
// It ties together the vault, filter, and env packages into a single
// Syncer type that:
//
//  1. Lists secret keys available at a configured Vault path.
//  2. Applies role-based filtering so only permitted keys are fetched.
//  3. Reads the allowed secrets from Vault (KV v2).
//  4. Writes the resulting key/value pairs to a local .env file,
//     creating a timestamped backup of any pre-existing file.
//
// Typical usage:
//
//	cfg, _ := config.Load("vaultpipe.yaml")
//	s, _ := sync.New(cfg)
//	if err := s.Run("production"); err != nil {
//		log.Fatal(err)
//	}
package sync
