ALTER TABLE vehicle_locations SET (
    timescaledb.compress,
    timescaledb.compress_segmentby = 'vehicle_id'
);

-- kompres data > 7 hari
SELECT add_compression_policy('vehicle_locations', compress_after => 7 * 86400);

-- hapus data > 30 hari
SELECT add_retention_policy('vehicle_locations', drop_after => 30 * 86400);
