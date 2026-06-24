# WatchDocs CLI

A lightweight CLI tool that scans a project for dependencies and returns their official documentation URLs as structured JSON — built for AI agents and developers.

## Install

```bash
npx watchdocs scan
```

No installation needed. Works via `npx` anywhere Node.js is available.

---

## Commands

### `watchdocs scan`

Recursively scans the current directory for manifest files, parses dependencies, and resolves official doc URLs from their respective registries.

```bash
watchdocs scan
watchdocs scan --path ./frontend
watchdocs scan --ecosystem npm,go
watchdocs scan --slim
```

| Flag | Description |
|---|---|
| `--path <dir>` | Target a specific directory instead of cwd |
| `--ecosystem <list>` | Filter by ecosystem(s), comma-separated |
| `--slim` | Return only `name` and `docUrl` per result (saves tokens) |
| `--clear-cache` | Clear the local cache before scanning |
| `--no-cache` | Disable reading and writing to the local cache |
| `--format, -f <format>` | Output format: `json` (default) or `list` |

---

### `watchdocs lookup`

Look up a single package by name without needing a manifest file.

```bash
watchdocs lookup express --ecosystem npm
watchdocs lookup github.com/spf13/cobra --ecosystem go
watchdocs lookup requests --ecosystem pip --slim
```

| Flag | Description |
|---|---|
| `--ecosystem <eco>` | Required — one of: `npm`, `go`, `pip`, `cargo`, `pub`, `maven` |
| `--slim` | Return only `name` and `docUrl` |
| `--clear-cache` | Clear the local cache before looking up |
| `--no-cache` | Disable reading and writing to the local cache |
| `--format, -f <format>` | Output format: `json` (default) or `list` |

---

## Output

```json
{
  "scanned": ["go.mod", "package.json"],
  "total": 5,
  "results": [
    {
      "name": "express",
      "version": "^4.18.0",
      "ecosystem": "npm",
      "type": "prod",
      "docUrl": "https://expressjs.com",
      "status": "resolved"
    }
  ]
}
```

`status` is either `"resolved"` or `"not_found"`.

---

## Supported Ecosystems

| Ecosystem | Manifest files | Registry |
|---|---|---|
| `npm` | `package.json` | registry.npmjs.org |
| `go` | `go.mod` | proxy.golang.org / pkg.go.dev |
| `pip` | `requirements.txt`, `pyproject.toml`, `uv.lock` | pypi.org |
| `cargo` | `Cargo.toml` | crates.io |
| `pub` | `pubspec.yaml` | pub.dev |
| `maven` | `pom.xml` | search.maven.org |
| `nuget` | `*.csproj`, `packages.config`, `Directory.Packages.props` | api.nuget.org |

---

## Building from Source

### Prerequisites
- [Go](https://go.dev/doc/install) 1.21+

```bash
git clone https://github.com/prohv/watchdocs-cli
cd watchdocs-cli
go build -o watchdocs .
```

---

## Project Structure

```
watchdocs-cli/
├── cmd/                  # CLI commands (scan, lookup)
├── internal/
│   ├── scanner/          # Recursive manifest file detection
│   ├── parser/           # Per-ecosystem manifest parsers
│   ├── resolver/         # Per-ecosystem online resolvers
│   └── models/           # Shared types (Dependency, DocResult)
├── bin/
│   └── watchdocs.js      # npm binary shim
├── binaries/             # Pre-built platform binaries (not in git)
├── package.json
└── main.go
```
