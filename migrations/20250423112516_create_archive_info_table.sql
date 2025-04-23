-- migrations/20250423112516_create_archive_info_table.sql

-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS archive_info (
    object_name TEXT PRIMARY KEY,
    bucket_name TEXT NOT NULL,
    size BIGINT NOT NULL,
    etag TEXT,
    uploaded_at TIMESTAMPTZ DEFAULT NOW() NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS archive_info;
-- +goose StatementEnd