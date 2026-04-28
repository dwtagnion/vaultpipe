// Package config handles loading and validating vaultpipe configuration
// from YAML files and environment variables.
package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config holds the top-level vaultpipe configuration.
type Config struct {
	// VaultAddr is the address of the Vault server.
	VaultAddr string `yaml:"vault_addr"`

	// VaultToken is the token used to authenticate with Vault.
	// Can also be set via VAULT_TOKEN environment variable.
	VaultToken string `yaml:"vault_token"`

	// SecretPath is the Vault KV path to read secrets from.
	SecretPath string `yaml:"secret_path"`

	// OutputFile is the path to the .env file to write.
	OutputFile string `yaml:"output_file"`

	// Role defines which secret keys this instance is allowed to consume.
	Role string `yaml:"role"`

	// Backup controls whether a backup of the existing .env file is created.
	Backup bool `yaml:"backup"`
}

// Load reads a YAML config file from the given path and returns a Config.
// Environment variables take precedence over file values for sensitive fields.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("config: reading file %q: %w", path, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("config: parsing YAML: %w", err)
	}

	// Environment variable overrides.
	if v := os.Getenv("VAULT_ADDR"); v != "" {
		cfg.VaultAddr = v
	}
	if v := os.Getenv("VAULT_TOKEN"); v != "" {
		cfg.VaultToken = v
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// Validate checks that required fields are present.
func (c *Config) Validate() error {
	if c.VaultAddr == "" {
		return fmt.Errorf("config: vault_addr is required")
	}
	if c.VaultToken == "" {
		return fmt.Errorf("config: vault_token is required (set vault_token in config or VAULT_TOKEN env var)")
	}
	if c.SecretPath == "" {
		return fmt.Errorf("config: secret_path is required")
	}
	if c.OutputFile == "" {
		c.OutputFile = ".env"
	}
	return nil
}
