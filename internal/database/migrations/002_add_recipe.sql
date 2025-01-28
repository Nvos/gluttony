-- +goose Up
CREATE TABLE recipes
(
    id                       INTEGER     NOT NULL,
    name                     TEXT UNIQUE NOT NULL,
    description              TEXT        NOT NULL DEFAULT '',
    instructions_markdown    TEXT        NOT NULL,
    thumbnail_url            TEXT        NOT NULL DEFAULT '',
    servings                 INTEGER     NOT NULL DEFAULT 1,
    cook_time_seconds        INTEGER     NOT NULL DEFAULT 0,
    preparation_time_seconds INTEGER     NOT NULL DEFAULT 0,
    source                   TEXT        NOT NULL DEFAULT '',
    created_at               DATETIME             DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at               DATETIME,

    PRIMARY KEY (id)
);

CREATE TABLE ingredients
(
    id   INTEGER     NOT NULL,
    name TEXT UNIQUE NOT NULL,

    PRIMARY KEY (id)
);

CREATE TABLE recipe_ingredients
(
    recipe_order  INTEGER NOT NULL,
    recipe_id     INTEGER NOT NULL REFERENCES recipes (id),
    ingredient_id INTEGER NOT NULL REFERENCES ingredients (id),
    unit          TEXT    NOT NULL DEFAULT 'g',
    quantity      INTEGER NOT NULL,

    PRIMARY KEY (recipe_id, ingredient_id)
);

CREATE TABLE recipe_tags
(
    recipe_order INTEGER NOT NULL,
    recipe_id    INTEGER NOT NULL REFERENCES recipes (id),
    tag_id       INTEGER NOT NULL REFERENCES tags (id),

    PRIMARY KEY (recipe_id, tag_id)
);

CREATE TABLE tags
(
    id   INTEGER     NOT NULL,
    name TEXT UNIQUE NOT NULL,

    PRIMARY KEY (id)
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
