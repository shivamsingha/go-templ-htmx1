-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD PRIMARY KEY (id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP CONSTRAINT users_pkey;
-- +goose StatementEnd
