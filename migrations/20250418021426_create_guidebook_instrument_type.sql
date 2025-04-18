-- +goose Up
-- +goose StatementBegin
CREATE TABLE instrument_type (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE instrument_type;
-- +goose StatementEnd