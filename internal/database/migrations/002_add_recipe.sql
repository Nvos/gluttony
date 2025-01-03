-- +goose Up
CREATE TABLE recipes
(
    id          INTEGER                            NOT NULL,
    name        TEXT UNIQUE                        NOT NULL,
    description TEXT                               NOT NULL,
    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,

    PRIMARY KEY (id)
);

CREATE TABLE images
(
    id INTEGER NOT NULL,
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
    recipe_id     INTEGER NOT NULL REFERENCES recipes (id),
    ingredient_id INTEGER NOT NULL REFERENCES ingredients (id)
);

CREATE TABLE recipe_steps
(
    id        INTEGER NOT NULL NULL,
    recipe_id INTEGER NOT NULL REFERENCES recipes (id) ON DELETE CASCADE,
    position  INTEGER,

    PRIMARY KEY (id)
);

-- +goose Down
DROP TABLE recipes;
DROP TABLE recipe_steps;
