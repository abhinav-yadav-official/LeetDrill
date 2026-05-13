-- +goose Up
-- +goose StatementBegin

ALTER TABLE users ADD COLUMN google_sub TEXT UNIQUE;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE users DROP COLUMN IF EXISTS google_sub;

-- +goose StatementEnd
