-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id uuid DEFAULT gen_random_uuid(),
    email text not null unique,
    name text not null,
    password text not null,
    email_verified boolean default false,
    created_at timestamp default CURRENT_TIMESTAMP,
    updated_at timestamp default CURRENT_TIMESTAMP
);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd