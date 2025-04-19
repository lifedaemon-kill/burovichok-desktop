-- +goose Up
-- +goose StatementBegin
CREATE TABLE instrument_type (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE instrument_type;
-- +goose StatementEnd
