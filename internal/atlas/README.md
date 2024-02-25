# Scripts
Generate migration

```shell
atlas migrate diff create_users --dir "file://migrations" --to "file://schema.hcl" --dev-url "sqlite://dev?mode=memory"
```
**Dev URL can be taken from https://atlasgo.io/concepts/dev-database**