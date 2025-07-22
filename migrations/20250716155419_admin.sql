-- +goose Up
-- +goose StatementBegin
INSERT INTO users (username, password, email, is_admin)
VALUES
    (
    'admin',
    'JDJhJDEwJFZpdWxBdHNCVHVDazlzLmY3WVMwTGU4LzUya2p6Li9sMUl3QW1HQW5tdVZ4U28vajNrS3hp',
    'admin@admin.ru',
    true
    )
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- +goose StatementEnd
