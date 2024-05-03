-- +goose Up
-- create enum type "unit"
CREATE TYPE "unit" AS ENUM ('weight', 'volume');
-- create "ingredients" table
CREATE TABLE "ingredients" (
  "id" serial NOT NULL,
  "name" jsonb NOT NULL,
  "unit" "unit" NOT NULL,
  PRIMARY KEY ("id")
);
-- create index "idx_ingredients_name_en" to table: "ingredients"
CREATE INDEX "idx_ingredients_name_en" ON "ingredients" ((to_tsvector('english'::regconfig, (name -> 'en'::text))));
-- create index "idx_ingredients_name_pl" to table: "ingredients"
CREATE INDEX "idx_ingredients_name_pl" ON "ingredients" ((to_tsvector('polish'::regconfig, (name -> 'pl'::text))));
-- create "recipes" table
CREATE TABLE "recipes" (
  "id" serial NOT NULL,
  "name" jsonb NOT NULL,
  "description" jsonb NOT NULL,
  "content" jsonb NOT NULL,
  PRIMARY KEY ("id")
);
-- create "recipes_ingredients" table
CREATE TABLE "recipes_ingredients" (
  "recipe_id" integer NOT NULL,
  "ingredient_id" integer NOT NULL,
  "note" text NOT NULL DEFAULT '',
  "amount" integer NOT NULL,
  "count" integer NOT NULL,
  PRIMARY KEY ("recipe_id", "ingredient_id")
);
-- set comment to column: "amount" on table: "recipes_ingredients"
COMMENT ON COLUMN "recipes_ingredients" ."amount" IS 'depending on ingredient unit, either ml or g';
-- set comment to column: "count" on table: "recipes_ingredients"
COMMENT ON COLUMN "recipes_ingredients" ."count" IS 'used to represent e.g. 3 apples, relevant for ingredients of weight type unit';
-- create "users" table
CREATE TABLE "users" (
  "id" serial NOT NULL,
  "name" character varying(100) NOT NULL,
  "password" character varying(100) NOT NULL,
  PRIMARY KEY ("id")
);
-- create index "uq_users_name" to table: "users"
CREATE UNIQUE INDEX "uq_users_name" ON "users" ("name");

-- +goose Down
-- reverse: create index "uq_users_name" to table: "users"
DROP INDEX "uq_users_name";
-- reverse: create "users" table
DROP TABLE "users";
-- reverse: set comment to column: "count" on table: "recipes_ingredients"
COMMENT ON COLUMN "recipes_ingredients" ."count" IS '';
-- reverse: set comment to column: "amount" on table: "recipes_ingredients"
COMMENT ON COLUMN "recipes_ingredients" ."amount" IS '';
-- reverse: create "recipes_ingredients" table
DROP TABLE "recipes_ingredients";
-- reverse: create "recipes" table
DROP TABLE "recipes";
-- reverse: create index "idx_ingredients_name_pl" to table: "ingredients"
DROP INDEX "idx_ingredients_name_pl";
-- reverse: create index "idx_ingredients_name_en" to table: "ingredients"
DROP INDEX "idx_ingredients_name_en";
-- reverse: create "ingredients" table
DROP TABLE "ingredients";
-- reverse: create enum type "unit"
DROP TYPE "unit";
