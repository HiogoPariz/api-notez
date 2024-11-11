-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE TABLE IF NOT EXISTS session (
    id serial PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS session;
-- +goose StatementEnd

