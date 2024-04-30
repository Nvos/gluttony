-- +goose Up

-- configure polish dictionary
DROP TEXT SEARCH CONFIGURATION IF EXISTS polish cascade;
CREATE TEXT SEARCH DICTIONARY polish_ispell (
    Template = ispell,
    DictFile = polish,
    AffFile = polish,
    StopWords = polish
    );

CREATE TEXT SEARCH CONFIGURATION polish( COPY = pg_catalog.english);

ALTER TEXT SEARCH CONFIGURATION polish
    ALTER MAPPING FOR asciiword, asciihword, hword_asciipart, word, hword, hword_part
        WITH polish_ispell;

CREATE EXTENSION pg_trgm;
-- create "ingredients" table
CREATE TABLE "ingredients"
(
    "id"   serial NOT NULL,
    "name" jsonb  NOT NULL,
    PRIMARY KEY ("id")
);
-- create index "idx_ingredients_name_en" to table: "ingredients"
CREATE INDEX "idx_ingredients_name_en" ON "ingredients" ((to_tsvector('english'::regconfig, (name -> 'en'::text))));
-- create index "idx_ingredients_name_pl" to table: "ingredients"
CREATE INDEX "idx_ingredients_name_pl" ON "ingredients" ((to_tsvector('polish'::regconfig, (name -> 'pl'::text))));
-- create index "uq_ingredients_name" to table: "ingredients"
CREATE UNIQUE INDEX "uq_ingredients_name" ON "ingredients" ("name");
-- create "recipes" table
CREATE TABLE "recipes"
(
    "id"          serial NOT NULL,
    "name"        jsonb  NOT NULL,
    "description" jsonb  NOT NULL,
    "content"     jsonb  NOT NULL,
    PRIMARY KEY ("id")
);
-- create index "uq_recipes_name" to table: "recipes"
CREATE UNIQUE INDEX "uq_recipes_name" ON "recipes" ("name");
-- create "recipes_ingredients" table
CREATE TABLE "recipes_ingredients"
(
    "recipe_id"     integer NOT NULL,
    "ingredient_id" integer NOT NULL,
    PRIMARY KEY ("recipe_id", "ingredient_id")
);
-- create "users" table
CREATE TABLE "users"
(
    "id"       serial                 NOT NULL,
    "name"     character varying(100) NOT NULL,
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
-- reverse: create "recipes_ingredients" table
DROP TABLE "recipes_ingredients";
-- reverse: create index "uq_recipes_name" to table: "recipes"
DROP INDEX "uq_recipes_name";
-- reverse: create "recipes" table
DROP TABLE "recipes";
-- reverse: create index "uq_ingredients_name" to table: "ingredients"
DROP INDEX "uq_ingredients_name";
-- reverse: create index "idx_ingredients_name_pl" to table: "ingredients"
DROP INDEX "idx_ingredients_name_pl";
-- reverse: create index "idx_ingredients_name_en" to table: "ingredients"
DROP INDEX "idx_ingredients_name_en";
-- reverse: create "ingredients" table
DROP TABLE "ingredients";
