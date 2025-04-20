-- +goose Up
CREATE TABLE table_four (
    research_id                 UUID             NOT NULL,
    measure_depth               DOUBLE PRECISION NOT NULL,
    true_vertical_depth         DOUBLE PRECISION NOT NULL,
    true_vertical_depth_sub_sea DOUBLE PRECISION NOT NULL
);

-- +goose Down
DROP TABLE table_four;