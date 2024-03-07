# Scripts
Generate migration

```shell
atlas migrate diff create_users --dir "file://migrations" --to "file://schema.hcl" --dev-url "sqlite://dev?mode=memory"
```
**Dev URL can be taken from https://atlasgo.io/concepts/dev-database**

# Atlas
In order for migrations to work atlas cli has to be installed globally take a look https://atlasgo.io/integrations/go-sdk

# Links
- https://github.com/sqlc-dev/sqlc/pull/2561/files