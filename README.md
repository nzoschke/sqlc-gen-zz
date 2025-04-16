# sqlc-gen-zz

[sqlc](https://sqlc.dev/) plugin for [zombiezen.com/go/sqlite](https://github.com/zombiezen/go-sqlite)

## Quick Start

Install sqlc and plugin, then `sqlc generate` in a project with a `sqlc.yaml` with the process plugin config:

```
brew install sqlc
go install github.com/nzoschke/sqlc-gen-zz@latest
sqlc generate
```

```yaml
version: "2"

plugins:
  - name: zz
    process:
      cmd: sqlc-gen-zz

sql:
  - engine: "sqlite"
    queries: "query.sql"
    schema: "schema.sql"
    codegen:
      - out: c
        plugin: zz
```

## SQLite Quirks

If a text field ends with `_at` it is converted to/from a `YYYY-MM-DD HH:MM:SS` string and a `time.Time` in support of:

```sql
created_at TEXT DEFAULT CURRENT_TIMESTAMP NOT NULL
```

## Development

```bash
go install ./... 
go generate ./...
go test ./...
```
