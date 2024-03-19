-- +goose Up
-- create "users" table
CREATE TABLE "users" (
  "id" serial NOT NULL,
  "name" character varying(100) NOT NULL,
  "password" character varying(100) NOT NULL,
  PRIMARY KEY ("id")
);
-- create index "uq_users_name" to table: "users"
CREATE UNIQUE INDEX "uq_users_name" ON "users" ("name");
-- create "recipes" table
CREATE TABLE "recipes" (
  "id" serial NOT NULL,
  "name" character varying(100) NOT NULL,
  "description" text NOT NULL,
  PRIMARY KEY ("id")
);
-- create index "uq_recipes_name" to table: "recipes"
CREATE UNIQUE INDEX "uq_recipes_name" ON "recipes" ("name");
-- create "recipe_steps" table
CREATE TABLE "recipe_steps" (
  "id" serial NOT NULL,
  "recipe_id" integer NOT NULL,
  "description" text NOT NULL,
  "order" integer NOT NULL,
  CONSTRAINT "fk_recipe" FOREIGN KEY ("recipe_id") REFERENCES "recipes" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- create index "uq_recipe_steps_order" to table: "recipe_steps"
CREATE UNIQUE INDEX "uq_recipe_steps_order" ON "recipe_steps" ("recipe_id", "order");

-- +goose Down
-- reverse: create index "uq_recipe_steps_order" to table: "recipe_steps"
DROP INDEX "uq_recipe_steps_order";
-- reverse: create "recipe_steps" table
DROP TABLE "recipe_steps";
-- reverse: create index "uq_recipes_name" to table: "recipes"
DROP INDEX "uq_recipes_name";
-- reverse: create "recipes" table
DROP TABLE "recipes";
-- reverse: create index "uq_users_name" to table: "users"
DROP INDEX "uq_users_name";
-- reverse: create "users" table
DROP TABLE "users";
