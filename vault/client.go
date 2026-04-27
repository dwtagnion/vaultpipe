package vault

import (
	"fmt"
	"os"

	vaultapi "github.com/hashicorp/vault/api"
)

// Client wraps the Vault API client with helper methods.
type Client struct {
	api    *vaultapi.Client
	Prefix string
}

// NewClient creates a new Vault client using environment variables or provided config.
func NewClient(addr, token string) (*Client, error) {
	config := vaultapi.DefaultConfig()

	if addr != "" {
		config.Address = addr
	} else if envAddr := os.Getenv("VAULT_ADDR"); envAddr != "" {
		config.Address = envAddr
	}

	client, err := vaultapi.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create vault client: %w", err)
	}

	if token != "" {
		client.SetToken(token)
	} else if envToken := os.Getenv("VAULT_TOKEN"); envToken != "" {
		client.SetToken(envToken)
	} else {
		return nil, fmt.Errorf("vault token not provided: set VAULT_TOKEN or use --token flag")
	}

	return &Client{api: client}, nil
}

// ReadSecret reads a KV v2 secret at the given path and returns its data map.
func (c *Client) ReadSecret(path string) (map[string]interface{}, error) {
	secret, err := c.api.Logical().Read(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read secret at %q: %w", path, err)
	}
	if secret == nil {
		return nil, fmt.Errorf("no secret found at path %q", path)
	}

	// KV v2 wraps data under secret.Data["data"]
	if data, ok := secret.Data["data"]; ok {
		if m, ok := data.(map[string]interface{}); ok {
			return m, nil
		}
	}

	// Fallback for KV v1
	return secret.Data, nil
}

// ListSecrets lists secret keys under the given path.
func (c *Client) ListSecrets(path string) ([]string, error) {
	secret, err := c.api.Logical().List(path)
	if err != nil {
		return nil, fmt.Errorf("failed to list secrets at %q: %w", path, err)
	}
	if secret == nil {
		return nil, fmt.Errorf("no secrets found at path %q", path)
	}

	keys, ok := secret.Data["keys"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected keys format at path %q", path)
	}

	result := make([]string, 0, len(keys))
	for _, k := range keys {
		if s, ok := k.(string); ok {
			result = append(result, s)
		}
	}
	return result, nil
}
