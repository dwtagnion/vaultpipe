// Package export provides functionality to export secrets from Vault
// into multiple output formats such as JSON, YAML, and shell scripts.
package export

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

// Format represents the output format for exported secrets.
type Format string

const (
	FormatJSON  Format = "json"
	FormatYAML  Format = "yaml"
	FormatShell Format = "shell"
)

// Exporter writes secrets to a given writer in the specified format.
type Exporter struct {
	format Format
	out    io.Writer
}

// New creates a new Exporter writing to out in the given format.
// If out is nil, os.Stdout is used.
func New(format Format, out io.Writer) (*Exporter, error) {
	switch format {
	case FormatJSON, FormatYAML, FormatShell:
		// valid
	default:
		return nil, fmt.Errorf("unsupported export format %q", format)
	}
	if out == nil {
		out = os.Stdout
	}
	return &Exporter{format: format, out: out}, nil
}

// Write outputs the provided secrets map in the configured format.
func (e *Exporter) Write(secrets map[string]string) error {
	keys := sortedKeys(secrets)
	switch e.format {
	case FormatJSON:
		return e.writeJSON(secrets, keys)
	case FormatYAML:
		return e.writeYAML(secrets, keys)
	case FormatShell:
		return e.writeShell(secrets, keys)
	}
	return nil
}

func (e *Exporter) writeJSON(secrets map[string]string, keys []string) error {
	ordered := make(map[string]string, len(keys))
	for _, k := range keys {
		ordered[k] = secrets[k]
	}
	enc := json.NewEncoder(e.out)
	enc.SetIndent("", "  ")
	return enc.Encode(ordered)
}

func (e *Exporter) writeYAML(secrets map[string]string, keys []string) error {
	for _, k := range keys {
		v := secrets[k]
		if strings.ContainsAny(v, " :\n") {
			v = fmt.Sprintf("%q", v)
		}
		if _, err := fmt.Fprintf(e.out, "%s: %s\n", k, v); err != nil {
			return err
		}
	}
	return nil
}

func (e *Exporter) writeShell(secrets map[string]string, keys []string) error {
	for _, k := range keys {
		v := secrets[k]
		if _, err := fmt.Fprintf(e.out, "export %s=%q\n", k, v); err != nil {
			return err
		}
	}
	return nil
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
