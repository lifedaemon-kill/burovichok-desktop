-- +goose Up
-- +goose StatementBegin
CREATE TABLE table_one (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    report_id INTEGER NOT NULL, -- Внешний ключ к reports.id
    timestamp DATETIME NOT NULL,
    pressure_depth REAL NOT NULL,
    temperature_depth REAL NOT NULL,

    FOREIGN KEY(report_id) REFERENCES reports(id) ON DELETE CASCADE -- При удалении отчета удалять связанные данные
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE table_one;
-- +goose StatementEnd