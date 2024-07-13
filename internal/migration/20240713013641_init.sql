-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS note(
  id serial PRIMARY KEY,
  title VARCHAR(32),
  content VARCHAR(32),
  active BOOLEAN,
  created_at TIMESTAMP,
  updated_at TIMESTAMP
);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS note;
-- +goose StatementEnd
