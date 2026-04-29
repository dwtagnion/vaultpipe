package template

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNew_EmptyTemplatePath(t *testing.T) {
	_, err := New("", "/tmp/out")
	if err == nil {
		t.Fatal("expected error for empty templatePath")
	}
}

func TestNew_EmptyOutputPath(t *testing.T) {
	_, err := New("/tmp/tmpl", "")
	if err == nil {
		t.Fatal("expected error for empty outputPath")
	}
}

func TestRender_BasicSecrets(t *testing.T) {
	dir := t.TempDir()

	tmplFile := filepath.Join(dir, "config.tmpl")
	content := "DB_HOST={{ index .Secrets \"db_host\" }}\nDB_PASS={{ index .Secrets \"db_pass\" }}\n"
	if err := os.WriteFile(tmplFile, []byte(content), 0600); err != nil {
		t.Fatalf("write template: %v", err)
	}

	outFile := filepath.Join(dir, "config.out")
	r, err := New(tmplFile, outFile)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	secrets := map[string]string{
		"db_host": "localhost",
		"db_pass": "s3cr3t",
	}
	if err := r.Render(secrets); err != nil {
		t.Fatalf("Render: %v", err)
	}

	got, err := os.ReadFile(outFile)
	if err != nil {
		t.Fatalf("read output: %v", err)
	}
	if !strings.Contains(string(got), "DB_HOST=localhost") {
		t.Errorf("expected DB_HOST=localhost in output, got: %s", got)
	}
	if !strings.Contains(string(got), "DB_PASS=s3cr3t") {
		t.Errorf("expected DB_PASS=s3cr3t in output, got: %s", got)
	}
}

func TestRender_TemplateFuncUpper(t *testing.T) {
	dir := t.TempDir()

	tmplFile := filepath.Join(dir, "func.tmpl")
	content := "{{ upper (index .Secrets \"env\") }}\n"
	if err := os.WriteFile(tmplFile, []byte(content), 0600); err != nil {
		t.Fatalf("write template: %v", err)
	}

	outFile := filepath.Join(dir, "func.out")
	r, _ := New(tmplFile, outFile)

	if err := r.Render(map[string]string{"env": "production"}); err != nil {
		t.Fatalf("Render: %v", err)
	}

	got, _ := os.ReadFile(outFile)
	if !strings.Contains(string(got), "PRODUCTION") {
		t.Errorf("expected PRODUCTION, got: %s", got)
	}
}

func TestRender_MissingTemplateFile(t *testing.T) {
	r, _ := New("/nonexistent/tmpl.txt", "/tmp/out.txt")
	if err := r.Render(nil); err == nil {
		t.Fatal("expected error for missing template file")
	}
}
