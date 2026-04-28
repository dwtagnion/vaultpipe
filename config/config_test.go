package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/your-org/vaultpipe/config"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "vaultpipe.yaml")
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatalf("writeTemp: %v", err)
	}
	return p
}

func TestLoad_ValidConfig(t *testing.T) {
	p := writeTemp(t, `
vault_addr: http://127.0.0.1:8200
vault_token: root
secret_path: secret/data/myapp
output_file: .env
role: backend
backup: true
`)
	cfg, err := config.Load(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Role != "backend" {
		t.Errorf("role: got %q, want %q", cfg.Role, "backend")
	}
	if !cfg.Backup {
		t.Error("expected backup to be true")
	}
}

func TestLoad_DefaultOutputFile(t *testing.T) {
	p := writeTemp(t, `
vault_addr: http://127.0.0.1:8200
vault_token: root
secret_path: secret/data/myapp
`)
	cfg, err := config.Load(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.OutputFile != ".env" {
		t.Errorf("output_file: got %q, want \".env\"", cfg.OutputFile)
	}
}

func TestLoad_MissingVaultAddr(t *testing.T) {
	p := writeTemp(t, `
vault_token: root
secret_path: secret/data/myapp
`)
	_, err := config.Load(p)
	if err == nil {
		t.Fatal("expected error for missing vault_addr")
	}
}

func TestLoad_TokenFromEnv(t *testing.T) {
	t.Setenv("VAULT_TOKEN", "env-token")
	p := writeTemp(t, `
vault_addr: http://127.0.0.1:8200
secret_path: secret/data/myapp
`)
	cfg, err := config.Load(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.VaultToken != "env-token" {
		t.Errorf("token: got %q, want %q", cfg.VaultToken, "env-token")
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := config.Load("/nonexistent/path/vaultpipe.yaml")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
