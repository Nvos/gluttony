-- +goose Up
CREATE TABLE users
(
    id       INTEGER     NOT NULL,
    username TEXT UNIQUE NOT NULL,
    password TEXT        NOT NULL,
    role     INTEGER     NOT NULL DEFAULT 1,
    PRIMARY KEY (id)
);

-- +goose Down
DROP TABLE users;
