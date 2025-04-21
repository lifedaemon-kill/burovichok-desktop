-- db/seed_data.sql
INSERT INTO oilfield (name) VALUES
    ('Test Field A'),
    ('Test Field B'),
    ('Test Field C');

INSERT INTO productive_horizon (name) VALUES
    ('Horizon Alpha'),
    ('Horizon Beta'),
    ('Horizon Gamma');

INSERT INTO instrument_type (name) VALUES
    ('Seismograph'),
    ('Well Log'),
    ('Core Sampler');

INSERT INTO research_type (name) VALUES
    ('Geological Survey'),
    ('Geophysical Analysis'),
    ('Reservoir Simulation');
