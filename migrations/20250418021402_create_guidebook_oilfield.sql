-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS oilfield (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE oilfield;
-- +goose StatementEnd
