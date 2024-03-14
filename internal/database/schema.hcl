schema "public" {}

table "users" {
  schema = schema.public
  column "id" {
    null = false
    type = serial
  }
  column "name" {
    null = false
    type = varchar(100)
  }
  column "password" {
    null = false
    type = varchar(100)
  }
  primary_key {
    columns = [column.id]
  }
  index "uq_users_name" {
    columns = [column.name]
    unique = true
  }
}

table "recipe_steps" {
  schema = schema.public
  column "id" {
    null = false
    type = serial
  }
  column "recipe_id" {
    null = false
    type = int
  }
  column "description" {
    null = false
    type = text
  }
  column "order" {
    null = false
    type = int
  }
  index "uq_recipe_steps_order" {
    columns = [column.recipe_id, column.order]
    unique = true
  }
  foreign_key "fk_recipe" {
    columns = [column.recipe_id]
    ref_columns = [table.recipes.column.id]
    on_delete = CASCADE
    on_update = NO_ACTION
  }
}

table "recipes" {
  schema = schema.public
  column "id" {
    null = false
    type = serial
  }
  column "name" {
    null = false
    type = varchar(100)
  }
  column "description" {
    null = false
    type = text
  }
  primary_key {
    columns = [column.id]
  }
  index "uq_recipes_name" {
    columns = [column.name]
    unique = true
  }
}