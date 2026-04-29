package rotate_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/user/vaultpipe/rotate"
	"github.com/user/vaultpipe/vault"
)

func buildRotateServer(data map[string]interface{}) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		payload := map[string]interface{}{
			"data": map[string]interface{}{"data": data},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(payload)
	}))
}

func newTestClient(t *testing.T, serverURL string) *vault.Client {
	t.Helper()
	c, err := vault.NewClient(serverURL, "test-token")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	return c
}

func TestDiff_NoLocalFile(t *testing.T) {
	srv := buildRotateServer(map[string]interface{}{"API_KEY": "abc123"})
	defer srv.Close()

	r := rotate.New(newTestClient(t, srv.URL))
	diff, err := r.Diff("secret/data/app", "/nonexistent/.env")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(diff.Added) != 1 || diff.Added[0] != "API_KEY" {
		t.Errorf("expected API_KEY in Added, got %v", diff.Added)
	}
	if diff.HasChanges() == false {
		t.Error("expected HasChanges to be true")
	}
}

func TestDiff_ChangedKey(t *testing.T) {
	srv := buildRotateServer(map[string]interface{}{"DB_PASS": "newpass"})
	defer srv.Close()

	tmp := filepath.Join(t.TempDir(), ".env")
	os.WriteFile(tmp, []byte("DB_PASS=oldpass\n"), 0600)

	r := rotate.New(newTestClient(t, srv.URL))
	diff, err := r.Diff("secret/data/app", tmp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(diff.Changed) != 1 || diff.Changed[0] != "DB_PASS" {
		t.Errorf("expected DB_PASS in Changed, got %v", diff.Changed)
	}
}

func TestDiff_RemovedKey(t *testing.T) {
	srv := buildRotateServer(map[string]interface{}{})
	defer srv.Close()

	tmp := filepath.Join(t.TempDir(), ".env")
	os.WriteFile(tmp, []byte("OLD_KEY=value\n"), 0600)

	r := rotate.New(newTestClient(t, srv.URL))
	diff, err := r.Diff("secret/data/app", tmp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(diff.Removed) != 1 || diff.Removed[0] != "OLD_KEY" {
		t.Errorf("expected OLD_KEY in Removed, got %v", diff.Removed)
	}
}

func TestDiff_NoChanges(t *testing.T) {
	srv := buildRotateServer(map[string]interface{}{"TOKEN": "xyz"})
	defer srv.Close()

	tmp := filepath.Join(t.TempDir(), ".env")
	os.WriteFile(tmp, []byte("TOKEN=xyz\n"), 0600)

	r := rotate.New(newTestClient(t, srv.URL))
	diff, err := r.Diff("secret/data/app", tmp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if diff.HasChanges() {
		t.Errorf("expected no changes, got %+v", diff)
	}
}
