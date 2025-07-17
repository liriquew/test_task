-- +goose Up
-- +goose StatementBegin
INSERT INTO users (username, password, email, is_admin)
VALUES
    (
    'admin',
    'JDJhJDEwJGhLVm1ObW5EVEoyZ2ZzTlRQYkRBZHVzNVI2aFB3ck5wUnZmUnExLldBZ1F4WktWRHBWYy9h',
    'admin@admin.ru',
    true
    )
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- +goose StatementEnd
