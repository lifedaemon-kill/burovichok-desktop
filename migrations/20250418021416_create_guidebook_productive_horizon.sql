-- +goose Up
-- +goose StatementBegin
CREATE TABLE productive_horizon (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE productive_horizon;
-- +goose StatementEnd