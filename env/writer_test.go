package env

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWrite_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")

	w := NewWriter(path, false)
	secrets := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
	}

	if err := w.Write(secrets); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("file not created: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "DB_HOST=localhost") {
		t.Errorf("expected DB_HOST in output, got:\n%s", content)
	}
	if !strings.Contains(content, "DB_PORT=5432") {
		t.Errorf("expected DB_PORT in output, got:\n%s", content)
	}
}

func TestWrite_SortedKeys(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")

	w := NewWriter(path, false)
	secrets := map[string]string{"Z_KEY": "z", "A_KEY": "a", "M_KEY": "m"}

	if err := w.Write(secrets); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(path)
	content := string(data)

	aPos := strings.Index(content, "A_KEY")
	mPos := strings.Index(content, "M_KEY")
	zPos := strings.Index(content, "Z_KEY")

	if !(aPos < mPos && mPos < zPos) {
		t.Errorf("keys not in sorted order: a=%d m=%d z=%d", aPos, mPos, zPos)
	}
}

func TestWrite_QuotesValueWithSpaces(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")

	w := NewWriter(path, false)
	if err := w.Write(map[string]string{"MSG": "hello world"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(path)
	if !strings.Contains(string(data), `MSG="hello world"`) {
		t.Errorf("expected quoted value, got: %s", string(data))
	}
}

func TestWrite_BackupCreated(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")

	// Write initial file.
	os.WriteFile(path, []byte("OLD=value\n"), 0600)

	w := NewWriter(path, true)
	if err := w.Write(map[string]string{"NEW": "value"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	entries, _ := os.ReadDir(dir)
	var backups []string
	for _, e := range entries {
		if strings.HasSuffix(e.Name(), ".bak") {
			backups = append(backups, e.Name())
		}
	}
	if len(backups) != 1 {
		t.Errorf("expected 1 backup file, found %d", len(backups))
	}
}
