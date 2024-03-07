-- Create "users" table
CREATE TABLE "users" (
  "id" serial NOT NULL,
  "name" character varying(100) NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "uq_users_name" to table: "users"
CREATE UNIQUE INDEX "uq_users_name" ON "users" ("name");
-- Create "recipes" table
CREATE TABLE "recipes" (
  "id" serial NOT NULL,
  "name" character varying(100) NOT NULL,
  "description" text NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "uq_recipes_name" to table: "recipes"
CREATE UNIQUE INDEX "uq_recipes_name" ON "recipes" ("name");
-- Create "recipe_steps" table
CREATE TABLE "recipe_steps" (
  "id" serial NOT NULL,
  "recipe_id" integer NOT NULL,
  "description" text NOT NULL,
  "order" integer NOT NULL,
  CONSTRAINT "fk_recipe" FOREIGN KEY ("recipe_id") REFERENCES "recipes" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "uq_recipe_steps_order" to table: "recipe_steps"
CREATE UNIQUE INDEX "uq_recipe_steps_order" ON "recipe_steps" ("recipe_id", "order");
