-- db/seed_data.sql
-- +goose up
INSERT INTO oilfield (name)
VALUES ('Test Field A'),
       ('Test Field B'),
       ('Test Field C')
ON CONFLICT (name) DO NOTHING;
INSERT INTO productive_horizon (name)
VALUES ('Horizon Alpha'),
       ('Horizon Beta'),
       ('Horizon Gamma')
ON CONFLICT (name) DO NOTHING;

INSERT INTO instrument_type (name)
VALUES ('Seismograph'),
       ('Well Log'),
       ('Core Sampler')
ON CONFLICT (name) DO NOTHING;

INSERT INTO research_type (name)
VALUES ('Geological Survey'),
       ('Geophysical Analysis'),
       ('Reservoir Simulation')
ON CONFLICT (name) DO NOTHING;
