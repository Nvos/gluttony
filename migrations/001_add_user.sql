-- +goose Up
CREATE TYPE users_role AS ENUM ('user', 'admin');

CREATE TABLE users
(
    id       SERIAL      NOT NULL,
    username TEXT UNIQUE NOT NULL,
    password TEXT        NOT NULL,
    role     users_role  NOT NULL DEFAULT 'admin',
    PRIMARY KEY (id)
);

-- +goose Down
DROP TABLE users;
