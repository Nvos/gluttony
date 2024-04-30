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

table "ingredients" {
  schema = schema.public
  column "id" {
    null = false
    type = serial
  }

  column "name" {
    null = false
    type = jsonb
  }

  primary_key {
    columns = [column.id]
  }

  index "uq_ingredients_name" {
    columns = [column.name]
    unique = true
  }

  index "idx_ingredients_name_en" {
    on {
      expr = "to_tsvector('english', (name->'en'))"
    }
  }

  index "idx_ingredients_name_pl" {
    on {
      expr = "to_tsvector('polish', (name->'pl'))"
    }
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
    type = jsonb
  }

  column "description" {
    null = false
    type = jsonb
  }

  column "content" {
    null = false
    type = jsonb
  }

  primary_key {
    columns = [column.id]
  }

  index "uq_recipes_name" {
    columns = [column.name]
    unique = true
  }
}

table "recipes_ingredients" {
  schema = schema.public
  column "recipe_id" {
    null = false
    type = integer
  }

  column "ingredient_id" {
    null = false
    type = integer
  }

  primary_key {
    columns = [column.recipe_id, column.ingredient_id]
  }
}