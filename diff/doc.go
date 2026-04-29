// Package diff compares a set of Vault secrets against an existing local
// .env file to produce a structured list of additions, removals and
// modifications.
//
// Basic usage:
//
//	// oldSecrets comes from parsing the current .env file.
//	// newSecrets comes from the Vault sync.
//	result := diff.Compare(oldSecrets, newSecrets)
//	if result.HasChanges() {
//		printer := diff.NewPrinter(os.Stdout, false)
//		printer.Print(result)
//	}
//
// The printer can optionally reveal secret values when showValues is true;
// take care not to enable this in shared or logged output.
package diff
