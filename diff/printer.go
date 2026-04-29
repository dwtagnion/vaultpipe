package diff

import (
	"fmt"
	"io"
	"os"
)

// Printer renders a diff Result to a writer in a human-friendly format.
type Printer struct {
	w        io.Writer
	showValues bool
}

// NewPrinter returns a Printer writing to w.
// If w is nil, os.Stdout is used.
// Set showValues to true to include old/new values in the output.
func NewPrinter(w io.Writer, showValues bool) *Printer {
	if w == nil {
		w = os.Stdout
	}
	return &Printer{w: w, showValues: showValues}
}

// Print writes the diff result to the configured writer.
func (p *Printer) Print(r *Result) {
	if !r.HasChanges() {
		fmt.Fprintln(p.w, "No changes detected.")
		return
	}

	for _, c := range r.Changes {
		switch c.Type {
		case Added:
			if p.showValues {
				fmt.Fprintf(p.w, "+ %s = %s\n", c.Key, c.NewValue)
			} else {
				fmt.Fprintf(p.w, "+ %s\n", c.Key)
			}
		case Removed:
			if p.showValues {
				fmt.Fprintf(p.w, "- %s = %s\n", c.Key, c.OldValue)
			} else {
				fmt.Fprintf(p.w, "- %s\n", c.Key)
			}
		case Modified:
			if p.showValues {
				fmt.Fprintf(p.w, "~ %s: %s -> %s\n", c.Key, c.OldValue, c.NewValue)
			} else {
				fmt.Fprintf(p.w, "~ %s\n", c.Key)
			}
		}
	}

	fmt.Fprintln(p.w, r.Summary())
}
