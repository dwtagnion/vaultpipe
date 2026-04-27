package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func mockVaultServer(t *testing.T, path string, response map[string]interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Errorf("failed to encode mock response: %v", err)
		}
	}))
}

func TestNewClient_MissingToken(t *testing.T) {
	t.Setenv("VAULT_TOKEN", "")
	_, err := NewClient("http://localhost:8200", "")
	if err == nil {
		t.Fatal("expected error when token is missing, got nil")
	}
}

func TestNewClient_WithToken(t *testing.T) {
	client, err := NewClient("http://localhost:8200", "test-token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestReadSecret_KVv2(t *testing.T) {
	response := map[string]interface{}{
		"data": map[string]interface{}{
			"data": map[string]interface{}{
				"API_KEY": "secret-value",
				"DB_PASS": "hunter2",
			},
		},
	}
	server := mockVaultServer(t, "/v1/secret/data/myapp", response)
	defer server.Close()

	client, err := NewClient(server.URL, "test-token")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	data, err := client.ReadSecret("secret/data/myapp")
	if err != nil {
		t.Fatalf("unexpected error reading secret: %v", err)
	}

	if data["API_KEY"] != "secret-value" {
		t.Errorf("expected API_KEY=secret-value, got %v", data["API_KEY"])
	}
}

func TestListSecrets(t *testing.T) {
	response := map[string]interface{}{
		"data": map[string]interface{}{
			"keys": []interface{}{"app1", "app2", "app3"},
		},
	}
	server := mockVaultServer(t, "/v1/secret/metadata", response)
	defer server.Close()

	client, err := NewClient(server.URL, "test-token")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	keys, err := client.ListSecrets("secret/metadata")
	if err != nil {
		t.Fatalf("unexpected error listing secrets: %v", err)
	}

	if len(keys) != 3 {
		t.Errorf("expected 3 keys, got %d", len(keys))
	}
}
