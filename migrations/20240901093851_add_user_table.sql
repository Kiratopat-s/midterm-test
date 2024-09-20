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

-- insert seed data
INSERT INTO users (username, password, position, first_name, last_name, photo_link) VALUES ('admin', '$2a$14$u4EmE.nz024ZNq7Bp8lTJ.esuLRT6DentEft.uO/CkXPaP7.aqM22', 'Admin', 'John', 'Doe', 'https://img.daisyui.com/images/stock/photo-1534528741775-53994a69daeb.webp');

-- +goose Down
DROP TABLE users;