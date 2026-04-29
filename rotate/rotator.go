// Package rotate provides functionality to detect and handle secret
// rotation by comparing current Vault secrets against existing .env files.
package rotate

import (
	"fmt"
	"os"
	"strings"

	"github.com/user/vaultpipe/vault"
)

// DiffResult holds the outcome of comparing Vault secrets to a local .env file.
type DiffResult struct {
	Added   []string
	Removed []string
	Changed []string
}

// HasChanges returns true when any difference was detected.
func (d DiffResult) HasChanges() bool {
	return len(d.Added) > 0 || len(d.Removed) > 0 || len(d.Changed) > 0
}

// Rotator compares Vault secrets against a local .env file.
type Rotator struct {
	client *vault.Client
}

// New creates a new Rotator backed by the given Vault client.
func New(c *vault.Client) *Rotator {
	return &Rotator{client: c}
}

// Diff reads secrets at secretPath from Vault and compares them to the
// key=value pairs found in envFile, returning a DiffResult.
func (r *Rotator) Diff(secretPath, envFile string) (DiffResult, error) {
	vaultSecrets, err := r.client.ReadSecret(secretPath)
	if err != nil {
		return DiffResult{}, fmt.Errorf("rotate: read vault secret: %w", err)
	}

	localSecrets, err := parseEnvFile(envFile)
	if err != nil {
		return DiffResult{}, fmt.Errorf("rotate: parse env file: %w", err)
	}

	var result DiffResult

	for k, v := range vaultSecrets {
		local, exists := localSecrets[k]
		if !exists {
			result.Added = append(result.Added, k)
		} else if local != fmt.Sprintf("%v", v) {
			result.Changed = append(result.Changed, k)
		}
	}

	for k := range localSecrets {
		if _, exists := vaultSecrets[k]; !exists {
			result.Removed = append(result.Removed, k)
		}
	}

	return result, nil
}

// parseEnvFile reads a .env file and returns a map of key→value pairs.
func parseEnvFile(path string) (map[string]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return map[string]string{}, nil
		}
		return nil, err
	}

	result := make(map[string]string)
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		result[parts[0]] = strings.Trim(parts[1], `"`)
	}
	return result, nil
}
