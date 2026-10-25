CREATE EXTENSION IF NOT EXISTS timescaledb;

CREATE TABLE IF NOT EXISTS vehicle_locations (
    vehicle_id  TEXT             NOT NULL,
    latitude    DOUBLE PRECISION NOT NULL,
    longitude   DOUBLE PRECISION NOT NULL,
    timestamp   BIGINT           NOT NULL,
    PRIMARY KEY (vehicle_id, timestamp)
);

-- 1 chunk = 1 hari
SELECT create_hypertable('vehicle_locations', by_range('timestamp', 86400), if_not_exists => TRUE);

-- kolom waktu pakai unix detik, jadi policy perlu tahu "now" dalam detik juga
CREATE OR REPLACE FUNCTION unix_now() RETURNS BIGINT
LANGUAGE SQL STABLE AS $$ SELECT extract(epoch FROM now())::BIGINT $$;

SELECT set_integer_now_func('vehicle_locations', 'unix_now', replace_if_exists => TRUE);
