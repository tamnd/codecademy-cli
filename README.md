# codecademy

Browse and search the [Codecademy](https://www.codecademy.com) course catalog from the command line.

`codecademy` is a single pure-Go binary. No API key required.

## Install

```bash
go install github.com/tamnd/codecademy-cli/cmd/codecademy@latest
```

Or grab a prebuilt binary from the [releases](https://github.com/tamnd/codecademy-cli/releases), or run the container image:

```bash
docker run --rm ghcr.io/tamnd/codecademy:latest --help
```

## Usage

```bash
# List all 800+ courses
codecademy list

# List first 20 courses in table format
codecademy list -n 20 -o table

# Search by title, slug, or description
codecademy search "python"
codecademy search "machine learning" -n 10

# Output formats
codecademy list -o json
codecademy list -o csv -n 50
codecademy search "javascript" -o jsonl
```

## Commands

| Command | Description |
|---------|-------------|
| `list` | List all courses in the Codecademy catalog |
| `search <query>` | Search courses by title, slug, or description (case-insensitive) |
| `version` | Show version information |

## Global flags

```
-o, --output string    output format: table|json|jsonl|csv|tsv|url|raw (default "auto")
-n, --limit int        limit number of records (0 = all)
    --fields strings   comma-separated columns to include
    --no-header        omit header row
    --template string  Go text/template per record
    --timeout duration per-request timeout (default 30s)
    --delay duration   minimum spacing between requests
    --retries int      retry attempts on 429/5xx (default 3)
```

## License

Apache-2.0. See [LICENSE](LICENSE).
