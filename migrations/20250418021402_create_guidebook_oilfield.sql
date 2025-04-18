-- +goose Up
-- +goose StatementBegin
CREATE TABLE oilfield (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE oilfield;
-- +goose StatementEnd