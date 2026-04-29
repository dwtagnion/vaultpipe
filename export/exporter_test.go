package export

import (
	"bytes"
	"strings"
	"testing"
)

var testSecrets = map[string]string{
	"DB_HOST":     "localhost",
	"DB_PASSWORD": "s3cr3t p@ss",
	"API_KEY":     "abc123",
}

func TestNew_InvalidFormat(t *testing.T) {
	_, err := New("toml", nil)
	if err == nil {
		t.Fatal("expected error for unsupported format")
	}
}

func TestNew_NilWriterDefaultsToStdout(t *testing.T) {
	e, err := New(FormatJSON, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if e.out == nil {
		t.Fatal("expected non-nil writer")
	}
}

func TestWrite_JSON(t *testing.T) {
	var buf bytes.Buffer
	e, _ := New(FormatJSON, &buf)
	if err := e.Write(testSecrets); err != nil {
		t.Fatalf("Write error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, `"API_KEY"`) {
		t.Error("expected API_KEY in JSON output")
	}
	if !strings.Contains(out, `"DB_HOST"`) {
		t.Error("expected DB_HOST in JSON output")
	}
}

func TestWrite_YAML(t *testing.T) {
	var buf bytes.Buffer
	e, _ := New(FormatYAML, &buf)
	if err := e.Write(testSecrets); err != nil {
		t.Fatalf("Write error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "API_KEY: abc123") {
		t.Errorf("expected API_KEY line in YAML, got:\n%s", out)
	}
	// value with space should be quoted
	if !strings.Contains(out, `"s3cr3t p@ss"`) {
		t.Errorf("expected quoted password in YAML, got:\n%s", out)
	}
}

func TestWrite_Shell(t *testing.T) {
	var buf bytes.Buffer
	e, _ := New(FormatShell, &buf)
	if err := e.Write(testSecrets); err != nil {
		t.Fatalf("Write error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "export API_KEY=") {
		t.Errorf("expected export statement, got:\n%s", out)
	}
	if !strings.Contains(out, "export DB_PASSWORD=") {
		t.Errorf("expected DB_PASSWORD export, got:\n%s", out)
	}
}

func TestWrite_SortedOutput(t *testing.T) {
	var buf bytes.Buffer
	e, _ := New(FormatYAML, &buf)
	_ = e.Write(testSecrets)
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if !strings.HasPrefix(lines[0], "API_KEY") {
		t.Errorf("expected API_KEY first (sorted), got %s", lines[0])
	}
}
