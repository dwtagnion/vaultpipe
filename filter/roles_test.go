package filter_test

import (
	"testing"

	"github.com/yourorg/vaultpipe/filter"
)

func testFilter() *filter.Filter {
	return filter.New([]filter.Role{
		{Name: "backend", Patterns: []string{"DB_*", "REDIS_URL"}},
		{Name: "frontend", Patterns: []string{"API_*", "PUBLIC_KEY"}},
	})
}

func TestMatch_EmptyRole(t *testing.T) {
	f := testFilter()
	if !f.Match("", "ANY_KEY") {
		t.Error("empty role should match any key")
	}
}

func TestMatch_WildcardPattern(t *testing.T) {
	f := testFilter()
	if !f.Match("backend", "DB_HOST") {
		t.Error("DB_HOST should match DB_* for backend role")
	}
	if !f.Match("backend", "DB_PORT") {
		t.Error("DB_PORT should match DB_* for backend role")
	}
}

func TestMatch_ExactPattern(t *testing.T) {
	f := testFilter()
	if !f.Match("backend", "REDIS_URL") {
		t.Error("REDIS_URL should match exact pattern for backend role")
	}
	if f.Match("backend", "REDIS_HOST") {
		t.Error("REDIS_HOST should NOT match backend role")
	}
}

func TestMatch_WrongRole(t *testing.T) {
	f := testFilter()
	if f.Match("frontend", "DB_HOST") {
		t.Error("DB_HOST should NOT match frontend role")
	}
}

func TestMatch_UnknownRole(t *testing.T) {
	f := testFilter()
	if f.Match("unknown", "DB_HOST") {
		t.Error("unknown role should not match any key")
	}
}

func TestRoleNames(t *testing.T) {
	f := testFilter()
	names := f.RoleNames()
	if len(names) != 2 {
		t.Fatalf("expected 2 role names, got %d", len(names))
	}
	if names[0] != "backend" || names[1] != "frontend" {
		t.Errorf("unexpected role names: %v", names)
	}
}
