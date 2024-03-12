# Sqlc

Generate type-safe sql from queries

**Requires [Sqlc cli](https://docs.sqlc.dev/en/stable/overview/install.html)** to be available in path

**Run from root of project**

```shell
sqlc generate
```

# Atlas

Generate migrations via diffing `schema.hcl` and current state of database

**Requires [Atlas cli](https://atlasgo.io/getting-started/) to be available in path**

**Run from root of project**

```shell
atlas migrate diff --env local {REPLACE_MIGRATION_NAME}
```

# Buf

Generate grpc/connect api

1. Configure golang codegen https://connectrpc.com/docs/go/getting-started/
2. Configure typescript codegen (install npm packages globally) https://connectrpc.com/docs/web/generating-code

**Run from root or project**

```shell
buf generate
```
