-- +goose Up
-- +goose StatementBegin
CREATE TABLE reports (
    id                          SERIAL PRIMARY KEY,
    field_name                  TEXT     NOT NULL,
    field_number                INTEGER  NOT NULL,
    cluster_number              INTEGER,
    horizon                     TEXT     NOT NULL,
    start_time                  TIMESTAMP NOT NULL,
    end_time                    TIMESTAMP NOT NULL,
    instrument_type             TEXT     NOT NULL,
    instrument_number           INTEGER,
    measure_depth               REAL     NOT NULL,
    true_vertical_depth         REAL     NOT NULL,
    true_vertical_depth_sub_sea REAL     NOT NULL,
    vdp_measured_depth          REAL     NOT NULL,
    vdp_true_vertical_depth     REAL,
    vdp_true_vertical_depth_sea REAL,
    diff_instrument_vdp         REAL,
    density_oil                 REAL     NOT NULL,
    density_liquid_stopped      REAL     NOT NULL,
    density_liquid_working      REAL     NOT NULL,
    pressure_diff_stopped       REAL,
    pressure_diff_working       REAL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE reports;
-- +goose StatementEnd
