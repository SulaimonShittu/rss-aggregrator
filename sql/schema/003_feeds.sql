-- +goose Up
CREATE TABLE feeds(
    id uuid PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    name text NOT NULL,
    url text unique not null,
    user_id uuid not null references users(id) on delete cascade
);

-- +goose Down
DROP TABLE feeds;