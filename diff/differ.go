// Package diff provides utilities for comparing Vault secrets
// against existing local .env files to surface changes.
package diff

import (
	"fmt"
	"sort"
	"strings"
)

// ChangeType describes the kind of change detected for a key.
type ChangeType string

const (
	Added    ChangeType = "added"
	Removed  ChangeType = "removed"
	Modified ChangeType = "modified"
	Unchanged ChangeType = "unchanged"
)

// Change represents a single key-level difference.
type Change struct {
	Key      string
	Type     ChangeType
	OldValue string
	NewValue string
}

// Result holds the full diff between two secret maps.
type Result struct {
	Changes []Change
}

// HasChanges returns true when at least one non-unchanged entry exists.
func (r *Result) HasChanges() bool {
	for _, c := range r.Changes {
		if c.Type != Unchanged {
			return true
		}
	}
	return false
}

// Summary returns a human-readable one-line summary.
func (r *Result) Summary() string {
	var added, removed, modified int
	for _, c := range r.Changes {
		switch c.Type {
		case Added:
			added++
		case Removed:
			removed++
		case Modified:
			modified++
		}
	}
	return fmt.Sprintf("+%d added, -%d removed, ~%d modified", added, removed, modified)
}

// Compare computes the diff between oldSecrets (local) and newSecrets (vault).
func Compare(oldSecrets, newSecrets map[string]string) *Result {
	seen := make(map[string]bool)
	var changes []Change

	for k, newVal := range newSecrets {
		seen[k] = true
		if oldVal, ok := oldSecrets[k]; !ok {
			changes = append(changes, Change{Key: k, Type: Added, NewValue: newVal})
		} else if oldVal != newVal {
			changes = append(changes, Change{Key: k, Type: Modified, OldValue: oldVal, NewValue: newVal})
		} else {
			changes = append(changes, Change{Key: k, Type: Unchanged, OldValue: oldVal, NewValue: newVal})
		}
	}

	for k, oldVal := range oldSecrets {
		if !seen[k] {
			changes = append(changes, Change{Key: k, Type: Removed, OldValue: oldVal})
		}
	}

	sort.Slice(changes, func(i, j int) bool {
		return strings.ToLower(changes[i].Key) < strings.ToLower(changes[j].Key)
	})

	return &Result{Changes: changes}
}
