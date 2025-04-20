-- +goose Up
CREATE TABLE reports (
    id                          SERIAL PRIMARY KEY,
    field_name                  TEXT    NOT NULL,
    field_number                INTEGER NOT NULL,
    cluster_number              INTEGER NOT NULL,
    horizon                     TEXT    NOT NULL,
    start_time                  TIMESTAMP NOT NULL,
    end_time                    TIMESTAMP NOT NULL,
    instrument_type             TEXT    NOT NULL,
    instrument_number           INTEGER NOT NULL,
    measure_depth               DOUBLE PRECISION NOT NULL,
    true_vertical_depth         DOUBLE PRECISION,
    true_vertical_depth_sub_sea DOUBLE PRECISION,
    vdp_measured_depth          DOUBLE PRECISION NOT NULL,
    vdp_true_vertical_depth     DOUBLE PRECISION,
    vdp_true_vertical_depth_sea DOUBLE PRECISION,
    diff_instrument_vdp         DOUBLE PRECISION,
    density_oil                 DOUBLE PRECISION NOT NULL,
    density_liquid_stopped      DOUBLE PRECISION NOT NULL,
    density_liquid_working      DOUBLE PRECISION NOT NULL,
    pressure_diff_stopped       DOUBLE PRECISION,
    pressure_diff_working       DOUBLE PRECISION
);

-- +goose Down
DROP TABLE IF EXISTS reports;
