package cache

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func tempCachePath(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, "cache.json")
}

func TestGet_MissOnEmpty(t *testing.T) {
	c, err := New(tempCachePath(t), time.Minute)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	_, ok := c.Get("secret/app")
	if ok {
		t.Fatal("expected miss on empty cache")
	}
}

func TestSet_ThenGet(t *testing.T) {
	c, err := New(tempCachePath(t), time.Minute)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	val := map[string]string{"DB_PASS": "secret"}
	if err := c.Set("secret/app", val); err != nil {
		t.Fatalf("Set: %v", err)
	}
	got, ok := c.Get("secret/app")
	if !ok {
		t.Fatal("expected hit after Set")
	}
	if got["DB_PASS"] != "secret" {
		t.Fatalf("unexpected value: %v", got)
	}
}

func TestGet_ExpiredEntry(t *testing.T) {
	c, err := New(tempCachePath(t), -time.Second)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	_ = c.Set("secret/app", map[string]string{"K": "V"})
	_, ok := c.Get("secret/app")
	if ok {
		t.Fatal("expected miss for expired entry")
	}
}

func TestInvalidate_RemovesEntry(t *testing.T) {
	c, err := New(tempCachePath(t), time.Minute)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	_ = c.Set("secret/app", map[string]string{"K": "V"})
	if err := c.Invalidate("secret/app"); err != nil {
		t.Fatalf("Invalidate: %v", err)
	}
	_, ok := c.Get("secret/app")
	if ok {
		t.Fatal("expected miss after Invalidate")
	}
}

func TestNew_LoadsExistingFile(t *testing.T) {
	path := tempCachePath(t)
	c1, _ := New(path, time.Minute)
	_ = c1.Set("secret/db", map[string]string{"PW": "abc"})

	c2, err := New(path, time.Minute)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	got, ok := c2.Get("secret/db")
	if !ok || got["PW"] != "abc" {
		t.Fatalf("expected persisted entry, got ok=%v val=%v", ok, got)
	}
}

func TestNew_MissingFileOK(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nonexistent", "cache.json")
	_, err := New(path, time.Minute)
	// directory doesn't exist yet — New should not error on missing file
	if err != nil && !os.IsNotExist(err) {
		t.Fatalf("unexpected error: %v", err)
	}
}
