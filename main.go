package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const version = "0.1.0"

// rootCmd is the base command for vaultpipe
var rootCmd = &cobra.Command{
	Use:   "vaultpipe",
	Short: "Sync secrets from HashiCorp Vault into local .env files",
	Long: `vaultpipe is a CLI tool that connects to a HashiCorp Vault instance
and syncs secrets into local .env files with support for role-based filtering.

It allows developers to pull only the secrets relevant to their role or
environment, keeping local development environments in sync with Vault.`,
	Version: version,
}

// syncCmd pulls secrets from Vault and writes them to a .env file
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync secrets from Vault to a local .env file",
	Long:  `Connects to Vault, fetches secrets based on the configured path and role, and writes them to a .env file.`,
	RunE:  runSync,
}

// listCmd lists available secret paths in Vault
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available secret paths in Vault",
	Long:  `Lists all secret paths accessible under the configured Vault mount point.`,
	RunE:  runList,
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().String("vault-addr", "", "Vault server address (overrides VAULT_ADDR env var)")
	rootCmd.PersistentFlags().String("vault-token", "", "Vault token (overrides VAULT_TOKEN env var)")
	rootCmd.PersistentFlags().String("config", "", "Path to vaultpipe config file (default: .vaultpipe.yaml)")

	// Sync-specific flags
	syncCmd.Flags().StringP("output", "o", ".env", "Output .env file path")
	syncCmd.Flags().StringP("path", "p", "", "Vault secret path to sync from (required)")
	syncCmd.Flags().StringP("role", "r", "", "Role filter to apply when selecting secrets")
	syncCmd.Flags().Bool("overwrite", false, "Overwrite existing .env file without prompting")
	syncCmd.Flags().Bool("dry-run", false, "Print secrets to stdout instead of writing to file")

	// List-specific flags
	listCmd.Flags().StringP("mount", "m", "secret", "Vault KV mount point")

	// Register subcommands
	rootCmd.AddCommand(syncCmd)
	rootCmd.AddCommand(listCmd)
}

// runSync is the handler for the sync subcommand
func runSync(cmd *cobra.Command, args []string) error {
	path, err := cmd.Flags().GetString("path")
	if err != nil || path == "" {
		return fmt.Errorf("--path flag is required for sync")
	}

	// TODO: implement full sync logic in internal/sync package
	fmt.Printf("[vaultpipe] syncing secrets from path: %s\n", path)
	return nil
}

// runList is the handler for the list subcommand
func runList(cmd *cobra.Command, args []string) error {
	mount, _ := cmd.Flags().GetString("mount")

	// TODO: implement list logic in internal/vault package
	fmt.Printf("[vaultpipe] listing secrets under mount: %s\n", mount)
	return nil
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
