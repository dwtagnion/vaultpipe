// Package rotate implements secret rotation detection for vaultpipe.
//
// It compares secrets stored in HashiCorp Vault against the contents of an
// existing local .env file and reports which keys have been added, removed,
// or changed. This allows callers to decide whether a re-sync is necessary
// and to surface drift between Vault and the local environment.
//
// Basic usage:
//
//	client, _ := vault.NewClient(addr, token)
//	r := rotate.New(client)
//	diff, err := r.Diff("secret/data/myapp", ".env")
//	if err != nil {
//		log.Fatal(err)
//	}
//	if diff.HasChanges() {
//		fmt.Println("secrets have drifted — re-sync recommended")
//	}
package rotate
