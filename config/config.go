package config

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/BurntSushi/toml"
)

// Config holds all runtime configuration for vaultpipe.
type Config struct {
	VaultAddr  string        `toml:"vault_addr"`
	VaultToken string        `toml:"vault_token"`
	SecretPath string        `toml:"secret_path"`
	OutputFile string        `toml:"output_file"`
	Role       string        `toml:"role"`
	AuditFile  string        `toml:"audit_file"`
	CacheFile  string        `toml:"cache_file"`
	CacheTTL   time.Duration `toml:"cache_ttl"`
}

// Load reads a TOML config file and merges environment variable overrides.
func Load(path string) (*Config, error) {
	cfg := &Config{
		OutputFile: ".env",
		CacheTTL:   5 * time.Minute,
	}

	if _, err := toml.DecodeFile(path, cfg); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("config: decode %s: %w", path, err)
		}
	}

	if v := os.Getenv("VAULT_ADDR"); v != "" {
		cfg.VaultAddr = v
	}
	if v := os.Getenv("VAULT_TOKEN"); v != "" {
		cfg.VaultToken = v
	}

	if cfg.VaultAddr == "" {
		return nil, errors.New("config: vault_addr is required")
	}
	if cfg.VaultToken == "" {
		return nil, errors.New("config: vault_token is required (set VAULT_TOKEN or vault_token in config)")
	}
	if cfg.SecretPath == "" {
		return nil, errors.New("config: secret_path is required")
	}

	if cfg.CacheFile == "" {
		cfg.CacheFile = ".vaultpipe_cache.json"
	}

	return cfg, nil
}
