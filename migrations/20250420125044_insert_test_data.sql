-- +goose Up
-- +goose StatementBegin
INSERT INTO oilfield (name) VALUES
    ('Test Field A'),
    ('Test Field B'),
    ('Test Field C');
-- +goose StatementEnd

-- +goose StatementBegin
INSERT INTO productive_horizon (name) VALUES
    ('Horizon Alpha'),
    ('Horizon Beta'),
    ('Horizon Gamma');
-- +goose StatementEnd

-- +goose StatementBegin
INSERT INTO instrument_type (name) VALUES
    ('Seismograph'),
    ('Well Log'),
    ('Core Sampler');
-- +goose StatementEnd

-- +goose StatementBegin
INSERT INTO research_type (name) VALUES
    ('Geological Survey'),
    ('Geophysical Analysis'),
    ('Reservoir Simulation');
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DELETE FROM oilfield
WHERE name IN ('Test Field A', 'Test Field B', 'Test Field C');
-- +goose StatementEnd

-- +goose StatementBegin
DELETE FROM productive_horizon
WHERE name IN ('Horizon Alpha', 'Horizon Beta', 'Horizon Gamma');
-- +goose StatementEnd

-- +goose StatementBegin
DELETE FROM instrument_type
WHERE name IN ('Seismograph', 'Well Log', 'Core Sampler');
-- +goose StatementEnd

-- +goose StatementBegin
DELETE FROM research_type
WHERE name IN ('Geological Survey', 'Geophysical Analysis', 'Reservoir Simulation');
-- +goose StatementEnd
