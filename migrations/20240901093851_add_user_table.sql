-- +goose Up
CREATE TABLE users (
    id          SERIAL PRIMARY KEY,
    username    VARCHAR(255) UNIQUE NOT NULL,
    password    VARCHAR(255) NOT NULL,
    position    VARCHAR(100),
    first_name  VARCHAR(100),
    last_name   VARCHAR(100),
    photo_link  TEXT
);

-- +goose Down
DROP TABLE users;