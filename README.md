# sqlc-gen-zz

[sqlc](https://sqlc.dev/) plugin for
[zombiezen.com/go/sqlite](https://github.com/zombiezen/go-sqlite)

## Quick Start

Install sqlc and plugin, then `sqlc generate` in a project with a `sqlc.yaml`
with the process plugin config:

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

## Time Convention

If a text field ends with `_at` it is converted from/to a `YYYY-MM-DD HH:MM:SS`
string and a `time.Time` in support of:

```sql
created_at TEXT DEFAULT CURRENT_TIMESTAMP NOT NULL
```

## JSON Overrides

By default this will use []byte for JSON column type.

But like
[sqlc JSON overrides](https://docs.sqlc.dev/en/latest/reference/datatypes.html#json),
you can specify a struct and this will marshal/unmarshal the struct
automatically.

Note that overrides are passed in through the plugin options.

## Development

```bash
rm -rf pkg/sql/c
go install ./... 
go generate ./...
go test ./...
```
