-- migrations/20250421200000_add_research_type_to_reports.sql

-- +goose Up
-- +goose StatementBegin
ALTER TABLE reports
ADD COLUMN research_type TEXT;
-- Сделаем его NOT NULL позже, если точно решим, что оно обязательное
-- Если оно точно обязательное СРАЗУ:
-- ALTER TABLE reports
-- ADD COLUMN research_type TEXT NOT NULL DEFAULT 'Unknown'; -- Или выбрать другое значение по умолчанию / убрать default и добавить проверку в коде
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE reports
DROP COLUMN research_type;
-- +goose StatementEnd