schema "public" {}

table "users" {
  schema = schema.public
  column "id" {
    null = false
    type = int
  }
  column "name" {
    null = false
    type = varchar(100)
  }
  primary_key {
    columns = [column.id]
  }
  index "uq_name" {
    columns = [column.name]
  }
}

table "recipes" {
  schema = schema.public
  column "id" {
    null = false
    type = int
  }
  column "name" {
    null = null
    type = varchar(100)
  }
  primary_key {
    columns = [column.id]
  }
  index "uq_name" {
    columns = [column.name]
  }
}