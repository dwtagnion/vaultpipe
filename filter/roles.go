// Package filter provides role-based filtering for Vault secrets.
package filter

import "strings"

// Role represents a named role with an associated set of secret key patterns.
type Role struct {
	Name     string
	Patterns []string
}

// Filter holds a collection of roles used to filter secrets.
type Filter struct {
	roles []Role
}

// New creates a new Filter from the given roles.
func New(roles []Role) *Filter {
	return &Filter{roles: roles}
}

// Match returns true if the given key matches any pattern in the specified role.
// If roleName is empty, all keys are considered a match.
func (f *Filter) Match(roleName, key string) bool {
	if roleName == "" {
		return true
	}
	for _, role := range f.roles {
		if role.Name != roleName {
			continue
		}
		for _, pattern := range role.Patterns {
			if matchPattern(pattern, key) {
				return true
			}
		}
	}
	return false
}

// RoleNames returns the list of configured role names.
func (f *Filter) RoleNames() []string {
	names := make([]string, 0, len(f.roles))
	for _, r := range f.roles {
		names = append(names, r.Name)
	}
	return names
}

// matchPattern checks whether key matches a simple glob pattern.
// Only a trailing '*' wildcard is supported (e.g. "DB_*").
func matchPattern(pattern, key string) bool {
	if strings.HasSuffix(pattern, "*") {
		return strings.HasPrefix(key, strings.TrimSuffix(pattern, "*"))
	}
	return pattern == key
}
