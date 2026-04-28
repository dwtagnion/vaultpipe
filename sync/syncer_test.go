package sync_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/example/vaultpipe/config"
	"github.com/example/vaultpipe/sync"
)

// buildTestServer returns a minimal Vault-compatible HTTP server that
// serves a KV v2 list and a single secret.
func buildTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()

	// LIST secrets
	mux.HandleFunc("/v1/secret/metadata/app", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.URL.Query().Get("list") != "true" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{"keys": []string{"DB_PASSWORD"}},
		})
	})

	// READ secret
	mux.HandleFunc("/v1/secret/data/app/DB_PASSWORD", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"data": map[string]string{"DB_PASSWORD": "s3cr3t"},
			},
		})
	})

	return httptest.NewServer(mux)
}

func TestRun_WritesMatchingSecrets(t *testing.T) {
	srv := buildTestServer(t)
	t.Cleanup(srv.Close)

	dir := t.TempDir()
	rolesFile := filepath.Join(dir, "roles.yaml")
	_ = os.WriteFile(rolesFile, []byte("roles:\n  dev:\n    - DB_PASSWORD\n"), 0o644)
	outFile := filepath.Join(dir, ".env")

	cfg := &config.Config{
		VaultAddr:  srv.URL,
		VaultToken: "test-token",
		SecretPath: "secret/data/app",
		RolesFile:  rolesFile,
		OutputFile: outFile,
	}

	s, err := sync.New(cfg)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	if err := s.Run("dev"); err != nil {
		t.Fatalf("Run: %v", err)
	}

	data, err := os.ReadFile(outFile)
	if err != nil {
		t.Fatalf("read output: %v", err)
	}

	if !strings.Contains(string(data), "DB_PASSWORD") {
		t.Errorf("expected DB_PASSWORD in output, got:\n%s", data)
	}
}

func TestRun_FilterExcludesUnmatchedKeys(t *testing.T) {
	srv := buildTestServer(t)
	t.Cleanup(srv.Close)

	dir := t.TempDir()
	rolesFile := filepath.Join(dir, "roles.yaml")
	// role "readonly" has no matching patterns
	_ = os.WriteFile(rolesFile, []byte("roles:\n  readonly:\n    - UNRELATED_KEY\n"), 0o644)
	outFile := filepath.Join(dir, ".env")

	cfg := &config.Config{
		VaultAddr:  srv.URL,
		VaultToken: "test-token",
		SecretPath: "secret/data/app",
		RolesFile:  rolesFile,
		OutputFile: outFile,
	}

	s, err := sync.New(cfg)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	if err := s.Run("readonly"); err != nil {
		t.Fatalf("Run: %v", err)
	}

	data, _ := os.ReadFile(outFile)
	if strings.Contains(string(data), "DB_PASSWORD") {
		t.Errorf("DB_PASSWORD should have been filtered out, got:\n%s", data)
	}
}
