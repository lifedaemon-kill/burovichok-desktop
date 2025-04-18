-- +goose Up
-- +goose StatementBegin
CREATE TABLE reports (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    field_id INTEGER NOT NULL,          -- FK к oilfield
    well_number INTEGER NOT NULL,
    cluster_number INTEGER,             -- Может быть не у всех? Сделал необязательным
    horizon_id INTEGER NOT NULL,        -- FK к productive_horizon
    start_time DATETIME NOT NULL,
    end_time DATETIME NOT NULL,
    instrument_type_id INTEGER NOT NULL, -- FK к instrument_type
    instrument_number INTEGER,          -- Может быть не у всех? Сделал необязательным
    vdp_measured_depth REAL NOT NULL,   -- MD ВДП из модели TableFive
    density_oil REAL NOT NULL,          -- Плотность для дебита нефти
    density_liquid_stopped REAL NOT NULL, -- Плотность жидкости в простое
    density_liquid_working REAL NOT NULL, -- Плотность жидкости в работе

    FOREIGN KEY(field_id) REFERENCES oilfield(id),
    FOREIGN KEY(horizon_id) REFERENCES productive_horizon(id),
    FOREIGN KEY(instrument_type_id) REFERENCES instrument_type(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE reports;
-- +goose StatementEnd