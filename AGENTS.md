# AGENTS.md

## watchdocs-cli

A CLI tool that scans a project for dependencies and returns their official documentation URLs as structured JSON.

## When to use

Use this instead of web searching for documentation URLs when working inside a project.

## Run

```bash
npx @watchdocs/cli scan
```

Scans the current directory recursively for manifest files and resolves doc URLs for all dependencies.

## Flags

### `scan`
| Flag | Example | When to use |
|---|---|---|
| `--path` | `scan --path ./frontend` | Target a specific sub-directory |
| `--ecosystem` | `scan --ecosystem npm,go` | Only resolve specific ecosystems |
| `--slim` | `scan --slim` | Save tokens — returns only `name` and `docUrl` |

### `lookup`
Single package lookup — no manifest needed.
```bash
npx @watchdocs/cli lookup <package> --ecosystem <eco>
npx @watchdocs/cli lookup express --ecosystem npm --slim
```
| Flag | When to use |
|---|---|
| `--ecosystem` | Required — specify ecosystem |
| `--slim` | Save tokens — returns only `name` and `docUrl` |

## Output

Returns a JSON object:

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

## Supported ecosystems

`npm` · `go` · `pip` · `cargo` · `pub` · `maven`

## Supported manifest files

`package.json` · `go.mod` · `requirements.txt` · `pyproject.toml` · `uv.lock` · `Cargo.toml` · `pom.xml` · `pubspec.yaml`
