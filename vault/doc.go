// Package vault provides a thin wrapper around the HashiCorp Vault API client
// for use by the vaultpipe CLI tool.
//
// It supports:
//   - Creating an authenticated Vault client from environment variables or flags
//   - Reading KV v1 and KV v2 secrets by path
//   - Listing secret keys under a given path prefix
//
// # Usage
//
//	client, err := vault.NewClient("", "") // uses VAULT_ADDR and VAULT_TOKEN env vars
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	data, err := client.ReadSecret("secret/data/myapp/prod")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	for key, val := range data {
//	    fmt.Printf("%s=%v\n", key, val)
//	}
//
// # Authentication
//
// The client reads VAULT_TOKEN from the environment if no token is explicitly
// passed. VAULT_ADDR defaults to http://127.0.0.1:8200 when not set.
package vault
