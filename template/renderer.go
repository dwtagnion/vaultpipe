// Package template provides functionality for rendering secrets into
// user-defined template files, allowing flexible secret injection beyond
// standard .env format.
package template

import (
	"fmt"
	"os"
	"strings"
	gotemplate "text/template"
)

// Renderer renders secrets into a Go template file and writes the result
// to a destination path.
type Renderer struct {
	templatePath string
	outputPath   string
}

// New creates a new Renderer that reads from templatePath and writes
// rendered output to outputPath.
func New(templatePath, outputPath string) (*Renderer, error) {
	if templatePath == "" {
		return nil, fmt.Errorf("template: templatePath must not be empty")
	}
	if outputPath == "" {
		return nil, fmt.Errorf("template: outputPath must not be empty")
	}
	return &Renderer{
		templatePath: templatePath,
		outputPath:   outputPath,
	}, nil
}

// Render executes the template at r.templatePath with the provided secrets map
// and writes the result to r.outputPath.
// The template receives a map[string]string named .Secrets.
func (r *Renderer) Render(secrets map[string]string) error {
	tmplBytes, err := os.ReadFile(r.templatePath)
	if err != nil {
		return fmt.Errorf("template: reading template file: %w", err)
	}

	funcMap := gotemplate.FuncMap{
		"upper": strings.ToUpper,
		"lower": strings.ToLower,
		"default": func(def, val string) string {
			if val == "" {
				return def
			}
			return val
		},
	}

	tmpl, err := gotemplate.New("vaultpipe").Funcs(funcMap).Parse(string(tmplBytes))
	if err != nil {
		return fmt.Errorf("template: parsing template: %w", err)
	}

	out, err := os.Create(r.outputPath)
	if err != nil {
		return fmt.Errorf("template: creating output file: %w", err)
	}
	defer out.Close()

	data := map[string]interface{}{
		"Secrets": secrets,
	}

	if err := tmpl.Execute(out, data); err != nil {
		return fmt.Errorf("template: executing template: %w", err)
	}
	return nil
}
