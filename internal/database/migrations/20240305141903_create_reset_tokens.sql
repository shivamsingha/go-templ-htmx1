-- +goose Up
-- +goose StatementBegin
CREATE TABLE reset_tokens (
    token text PRIMARY KEY,
    user_id uuid not null,
    created_at timestamp DEFAULT CURRENT_TIMESTAMP,
    expires_at timestamp DEFAULT (CURRENT_TIMESTAMP + INTERVAL '1 hour'),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS reset_tokens;
-- +goose StatementEnd