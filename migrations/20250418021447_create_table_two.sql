-- +goose Up
-- +goose StatementBegin
CREATE TABLE table_two (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    report_id INTEGER NOT NULL, -- Внешний ключ к reports.id
    timestamp_tubing DATETIME NOT NULL,
    pressure_tubing REAL NOT NULL,
    timestamp_annulus DATETIME NOT NULL,
    pressure_annulus REAL NOT NULL,
    timestamp_linear DATETIME NOT NULL,
    pressure_linear REAL NOT NULL,

    FOREIGN KEY(report_id) REFERENCES reports(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE table_two;
-- +goose StatementEnd