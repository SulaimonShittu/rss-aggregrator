-- +goose Up

CREATE TABLE feed_follows(
    id uuid PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    user_id uuid not null references users(id) on delete cascade,
    feed_id uuid not null references feeds(id) on delete cascade,
    UNIQUE(user_id, feed_id)
);

-- +goose Down
DROP TABLE feed_follows;