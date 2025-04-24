-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS research_type
(
    id   SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE research_type;
-- +goose StatementEnd
