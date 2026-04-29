package diff

import (
	"testing"
)

func TestCompare_Added(t *testing.T) {
	old := map[string]string{}
	new_ := map[string]string{"FOO": "bar"}

	r := Compare(old, new_)
	if len(r.Changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(r.Changes))
	}
	if r.Changes[0].Type != Added {
		t.Errorf("expected Added, got %s", r.Changes[0].Type)
	}
	if !r.HasChanges() {
		t.Error("expected HasChanges to be true")
	}
}

func TestCompare_Removed(t *testing.T) {
	old := map[string]string{"FOO": "bar"}
	new_ := map[string]string{}

	r := Compare(old, new_)
	if len(r.Changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(r.Changes))
	}
	if r.Changes[0].Type != Removed {
		t.Errorf("expected Removed, got %s", r.Changes[0].Type)
	}
}

func TestCompare_Modified(t *testing.T) {
	old := map[string]string{"FOO": "old"}
	new_ := map[string]string{"FOO": "new"}

	r := Compare(old, new_)
	if r.Changes[0].Type != Modified {
		t.Errorf("expected Modified, got %s", r.Changes[0].Type)
	}
	if r.Changes[0].OldValue != "old" || r.Changes[0].NewValue != "new" {
		t.Errorf("unexpected values: old=%s new=%s", r.Changes[0].OldValue, r.Changes[0].NewValue)
	}
}

func TestCompare_Unchanged(t *testing.T) {
	old := map[string]string{"FOO": "same"}
	new_ := map[string]string{"FOO": "same"}

	r := Compare(old, new_)
	if r.HasChanges() {
		t.Error("expected no changes")
	}
}

func TestCompare_SortedKeys(t *testing.T) {
	old := map[string]string{}
	new_ := map[string]string{"ZZZ": "1", "AAA": "2", "MMM": "3"}

	r := Compare(old, new_)
	keys := make([]string, len(r.Changes))
	for i, c := range r.Changes {
		keys[i] = c.Key
	}
	if keys[0] != "AAA" || keys[1] != "MMM" || keys[2] != "ZZZ" {
		t.Errorf("keys not sorted: %v", keys)
	}
}

func TestResult_Summary(t *testing.T) {
	old := map[string]string{"A": "1", "B": "old"}
	new_ := map[string]string{"B": "new", "C": "3"}

	r := Compare(old, new_)
	s := r.Summary()
	if s != "+1 added, -1 removed, ~1 modified" {
		t.Errorf("unexpected summary: %s", s)
	}
}
