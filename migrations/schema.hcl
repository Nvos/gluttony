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

enum "unit" {
  schema = schema.public
  values = ["weight", "volume"]
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

  column "unit" {
    null = false
    type = enum.unit
  }

  primary_key {
    columns = [column.id]
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
}

table "recipes_ingredients" {
  schema = schema.public
  column "recipe_id" {
    type = integer
  }

  column "ingredient_id" {
    type = integer
  }

  column "note" {
    type = text
    default = ""
  }

  column "amount" {
    null = false
    type = integer
    comment = "depending on ingredient unit, either ml or g"
  }

  column "count" {
    null = false
    type = integer
    comment = "used to represent e.g. 3 apples, relevant for ingredients of weight type unit"
  }

  primary_key {
    columns = [column.recipe_id, column.ingredient_id]
  }
}