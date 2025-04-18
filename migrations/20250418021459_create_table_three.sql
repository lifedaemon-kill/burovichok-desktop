-- +goose Up
-- +goose StatementBegin
CREATE TABLE table_three (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    report_id INTEGER NOT NULL, -- Внешний ключ к reports.id
    timestamp DATETIME NOT NULL,
    flow_liquid REAL NOT NULL,
    water_cut REAL NOT NULL,
    flow_gas REAL NOT NULL,

    FOREIGN KEY(report_id) REFERENCES reports(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE table_three;
-- +goose StatementEnd