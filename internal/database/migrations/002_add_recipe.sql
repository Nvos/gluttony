-- +goose Up
CREATE TABLE recipes
(
    id                       SERIAL PRIMARY KEY,
    name                     TEXT UNIQUE NOT NULL,
    description              TEXT        NOT NULL DEFAULT '',
    instructions_markdown    TEXT        NOT NULL,
    thumbnail_url            TEXT        NOT NULL DEFAULT '',
    servings                 INTEGER     NOT NULL DEFAULT 1,
    cook_time_seconds        INTEGER     NOT NULL DEFAULT 0,
    preparation_time_seconds INTEGER     NOT NULL DEFAULT 0,
    source                   TEXT        NOT NULL DEFAULT '',
    created_at               TIMESTAMPTZ          default (now() at time zone 'utc'),
    updated_at               TIMESTAMPTZ,
    owner_id                 INTEGER     NOT NULL REFERENCES users (id)
);

CREATE TABLE ingredients
(
    id   SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL
);

CREATE TABLE recipe_ingredients
(
    recipe_order  INTEGER NOT NULL,
    recipe_id     INTEGER NOT NULL REFERENCES recipes (id),
    ingredient_id INTEGER NOT NULL REFERENCES ingredients (id),
    unit          TEXT    NOT NULL DEFAULT 'g',
    quantity      REAL    NOT NULL,
    note          TEXT    NOT NULL DEFAULT '',

    PRIMARY KEY (recipe_id, ingredient_id)
);

CREATE TABLE tags
(
    id   SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL
);

CREATE TABLE recipe_tags
(
    recipe_order INTEGER NOT NULL,
    recipe_id    INTEGER NOT NULL REFERENCES recipes (id),
    tag_id       INTEGER NOT NULL REFERENCES tags (id),

    PRIMARY KEY (recipe_id, tag_id)
);


CREATE TABLE recipe_nutrition
(
    recipe_id INTEGER NOT NULL REFERENCES recipes (id),

    calories  REAL    NOT NULL DEFAULT 0,
    fat       REAL    NOT NULL DEFAULT 0,
    carbs     REAL    NOT NULL DEFAULT 0,
    protein   REAL    NOT NULL DEFAULT 0,

    PRIMARY KEY (recipe_id)
);

-- +goose Down
DROP TABLE recipe_ingredients;
DROP TABLE ingredients;
DROP TABLE recipe_tags;
DROP TABLE tags;
DROP TABLE recipe_nutrition;
DROP TABLE recipes;
