-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    username VARCHAR(127) UNIQUE NOT NULL,
    password CHAR(80) NOT NULL,
    email VARCHAR(127) NOT NULL,
    is_admin bool NOT NULL DEFAULT false
);

CREATE INDEX idx_users_username ON users (username);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS idx_users_username;
DROP DATABASE IF EXISTS users;

-- +goose StatementEnd
