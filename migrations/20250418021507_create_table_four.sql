-- +goose Up
-- +goose StatementBegin
CREATE TABLE table_four (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    report_id INTEGER NOT NULL, -- Внешний ключ к reports.id
    measured_depth REAL NOT NULL,
    true_vertical_depth REAL NOT NULL,
    true_vertical_depth_sub_sea REAL NOT NULL,

    FOREIGN KEY(report_id) REFERENCES reports(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE table_four;
-- +goose StatementEnd