# vaultwatch

A CLI tool that monitors HashiCorp Vault secret expiry and sends configurable alerts before leases expire.

---

## Installation

```bash
go install github.com/yourusername/vaultwatch@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/vaultwatch.git
cd vaultwatch && go build -o vaultwatch .
```

---

## Usage

Set your Vault address and token, then run vaultwatch with a config file:

```bash
export VAULT_ADDR="https://vault.example.com"
export VAULT_TOKEN="s.your-token-here"

vaultwatch --config config.yaml
```

**Example `config.yaml`:**

```yaml
alert_threshold: 72h
secrets:
  - path: secret/data/my-app/db-credentials
  - path: secret/data/my-app/api-keys
notifications:
  slack:
    webhook_url: "https://hooks.slack.com/services/..."
  email:
    to: "ops-team@example.com"
```

vaultwatch will poll Vault at a configurable interval and send alerts when a secret lease is within the defined threshold of expiry.

**Flags:**

| Flag | Default | Description |
|------|---------|-------------|
| `--config` | `config.yaml` | Path to config file |
| `--interval` | `1h` | Polling interval |
| `--dry-run` | `false` | Log alerts without sending |

---

## License

MIT © [yourusername](https://github.com/yourusername)