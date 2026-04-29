# vaultpipe

> A CLI tool to sync secrets from HashiCorp Vault into local `.env` files with role-based filtering.

---

## Installation

```bash
go install github.com/yourusername/vaultpipe@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/vaultpipe.git
cd vaultpipe
go build -o vaultpipe .
```

---

## Usage

Set your Vault address and token, then run:

```bash
export VAULT_ADDR="https://vault.example.com"
export VAULT_TOKEN="s.yourtoken"

vaultpipe sync --path secret/myapp --role backend --output .env
```

This will pull all secrets from the specified Vault path, filter them by the given role, and write them to a `.env` file in the current directory.

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--path` | Vault secret path | _(required)_ |
| `--role` | Role filter to apply | `default` |
| `--output` | Output file path | `.env` |
| `--overwrite` | Overwrite existing file | `false` |
| `--dry-run` | Print secrets to stdout without writing to file | `false` |

### Example Output

```env
DB_HOST=db.internal
DB_PASSWORD=supersecret
API_KEY=abc123
```

---

## Environment Variables

Instead of passing flags each time, you can configure vaultpipe using environment variables:

| Variable | Equivalent Flag |
|----------|-----------------|
| `VAULT_ADDR` | Vault server address |
| `VAULT_TOKEN` | Vault authentication token |
| `VAULTPIPE_ROLE` | `--role` |
| `VAULTPIPE_OUTPUT` | `--output` |

---

## Requirements

- Go 1.21+
- HashiCorp Vault with a valid token or AppRole credentials

---

## License

MIT © 2024 yourusername
