-- +goose Up
CREATE TABLE items (
    id         bigserial PRIMARY KEY,
    title      VARCHAR(255) NOT NULL,
    amount     INT NOT NULL,
    quantity   INT NOT NULL,
    status     VARCHAR(20) NOT NULL,
    owner_id   INT NOT NULL
);

-- insert seed data
INSERT INTO items (title, amount, quantity, status, owner_id) VALUES ('item1', 100, 10, 'REJECTED', 1);
INSERT INTO items (title, amount, quantity, status, owner_id) VALUES ('item2', 200, 20, 'APPROVED', 1);
INSERT INTO items (title, amount, quantity, status, owner_id) VALUES ('item3', 300, 30, 'PENDING', 1);

-- +goose Down
DROP TABLE items;