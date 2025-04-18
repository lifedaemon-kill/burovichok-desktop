-- +goose Up
-- +goose StatementBegin
CREATE TABLE reports
(
    id                          INTEGER PRIMARY KEY AUTOINCREMENT,
    field_name                  TEXT     NOT NULL,
    field__number               INTEGER  NOT NULL,
    cluster_number              INTEGER,           -- Может быть не у всех
    horizon                     TEXT     NOT NULL,
    start_time                  DATETIME NOT NULL,
    end_time                    DATETIME NOT NULL,
    instrument_type             TEXT     NOT NULL,
    instrument_number           INTEGER,           -- Может быть не у всех
    measure_depth               REAL     NOT NULL,
    true_vertical_depth         REAL     NOT NULL,
    true_vertical_depth_sub_sea REAL     NOT NULL,
    vdp_measured_depth          REAL     NOT NULL, -- MD ВДП из модели TableFive
    vdp_true_vertical_depth     REAL,
    vdp_true_vertical_depth_sea REAL,
    diff_instrument_vdp         REAL,

    density_oil                 REAL     NOT NULL, -- Плотность для дебита нефти
    density_liquid_stopped      REAL     NOT NULL, -- Плотность жидкости в простое
    density_liquid_working      REAL     NOT NULL, -- Плотность жидкости в работе
    pressure_diff_stopped       REAL,
    pressure_diff_working       REAL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE reports;
-- +goose StatementEnd