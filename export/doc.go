// Package export provides multi-format secret export capabilities for vaultpipe.
//
// It supports writing secrets retrieved from HashiCorp Vault in three formats:
//
//   - JSON  – a pretty-printed JSON object keyed by secret name.
//   - YAML  – a simple key: value YAML document (no external dependency).
//   - Shell – a series of `export KEY="value"` statements suitable for
//     sourcing directly into a shell session.
//
// All formats emit keys in lexicographic order for deterministic, diff-friendly
// output.
//
// Example:
//
//	exporter, err := export.New(export.FormatShell, os.Stdout)
//	if err != nil {
//		log.Fatal(err)
//	}
//	if err := exporter.Write(secrets); err != nil {
//		log.Fatal(err)
//	}
package export
