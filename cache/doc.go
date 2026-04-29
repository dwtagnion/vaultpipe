// Package cache provides a lightweight, file-backed, TTL-aware cache for
// Vault secret values. It is used by vaultpipe to avoid redundant Vault API
// calls during repeated sync or watch operations.
//
// Entries are stored as JSON on disk so that they survive process restarts
// within their TTL window. The cache file is written with mode 0600 to
// protect sensitive data at rest.
//
// Basic usage:
//
//	c, err := cache.New(".vaultpipe_cache.json", 5*time.Minute)
//	if err != nil { ... }
//
//	if val, ok := c.Get("secret/myapp"); ok {
//		// use cached value
//	} else {
//		// fetch from Vault, then:
//		_ = c.Set("secret/myapp", val)
//	}
package cache
