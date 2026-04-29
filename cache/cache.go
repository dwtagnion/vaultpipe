package cache

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Entry holds a cached secret value with an expiry timestamp.
type Entry struct {
	Value     map[string]string `json:"value"`
	ExpiresAt time.Time         `json:"expires_at"`
}

// Cache is a simple file-backed in-memory secret cache.
type Cache struct {
	mu      sync.RWMutex
	path    string
	ttl     time.Duration
	entries map[string]Entry
}

// New creates a Cache backed by the given file path with the specified TTL.
func New(path string, ttl time.Duration) (*Cache, error) {
	c := &Cache{
		path:    path,
		ttl:     ttl,
		entries: make(map[string]Entry),
	}
	if err := c.load(); err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	return c, nil
}

// Get returns the secret map for key if present and not expired.
func (c *Cache) Get(key string) (map[string]string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	e, ok := c.entries[key]
	if !ok || time.Now().After(e.ExpiresAt) {
		return nil, false
	}
	return e.Value, true
}

// Set stores the secret map for key and persists the cache to disk.
func (c *Cache) Set(key string, value map[string]string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[key] = Entry{
		Value:     value,
		ExpiresAt: time.Now().Add(c.ttl),
	}
	return c.save()
}

// Invalidate removes the entry for key from the cache.
func (c *Cache) Invalidate(key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.entries, key)
	return c.save()
}

func (c *Cache) load() error {
	data, err := os.ReadFile(c.path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &c.entries)
}

func (c *Cache) save() error {
	if err := os.MkdirAll(filepath.Dir(c.path), 0o700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(c.entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(c.path, data, 0o600)
}
