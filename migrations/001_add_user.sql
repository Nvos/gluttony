-- +goose Up
CREATE TYPE users_role AS ENUM ('user', 'admin');

CREATE TABLE users
(
    id       INTEGER PRIMARY KEY GENERATED BY DEFAULT AS IDENTITY,
    username TEXT UNIQUE NOT NULL,
    password TEXT        NOT NULL,
    role     users_role  NOT NULL DEFAULT 'user'
);

-- +goose Down
DROP TABLE users;
