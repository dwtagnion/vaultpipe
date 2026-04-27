// Package filter implements role-based filtering for Vault secret keys.
//
// A Role defines a named group of secret key patterns. Patterns support
// a simple trailing wildcard ('*') to match key prefixes, or exact string
// matching for precise control.
//
// Example usage:
//
//	f := filter.New([]filter.Role{
//		{Name: "backend", Patterns: []string{"DB_*", "REDIS_URL"}},
//		{Name: "frontend", Patterns: []string{"API_*"}},
//	})
//
//	// Check whether a key belongs to the backend role:
//	if f.Match("backend", secretKey) {
//		// include in .env output
//	}
//
// Passing an empty string as the role name bypasses filtering and
// matches every key, which is useful when no role restriction is needed.
package filter
