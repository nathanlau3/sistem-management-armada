package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"fleet/internal/core/domain"
)

type LocationRepository struct {
	db *pgxpool.Pool
}

func NewLocationRepository(db *pgxpool.Pool) *LocationRepository {
	return &LocationRepository{db: db}
}

func (r *LocationRepository) Save(ctx context.Context, loc domain.Location) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO vehicle_locations (vehicle_id, latitude, longitude, timestamp) VALUES ($1, $2, $3, $4)
		ON CONFLICT (vehicle_id, timestamp) DO NOTHING`,
		loc.VehicleID, loc.Latitude, loc.Longitude, loc.Timestamp,
	)
	return err
}

func (r *LocationRepository) Latest(ctx context.Context, vehicleID string) (domain.Location, error) {
	var loc domain.Location
	err := r.db.QueryRow(ctx,
		`SELECT vehicle_id, latitude, longitude, timestamp
		FROM vehicle_locations
		WHERE vehicle_id = $1
		ORDER BY timestamp DESC
		LIMIT 1`,
		vehicleID,
	).Scan(&loc.VehicleID, &loc.Latitude, &loc.Longitude, &loc.Timestamp)
	if errors.Is(err, pgx.ErrNoRows) {
		return loc, domain.ErrNotFound
	}
	return loc, err
}

func (r *LocationRepository) History(ctx context.Context, vehicleID string, start, end int64) ([]domain.Location, error) {
	rows, err := r.db.Query(ctx,
		`SELECT vehicle_id, latitude, longitude, timestamp
		FROM vehicle_locations
		WHERE vehicle_id = $1 AND timestamp BETWEEN $2 AND $3
		ORDER BY timestamp`,
		vehicleID, start, end,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	locations := []domain.Location{}
	for rows.Next() {
		var loc domain.Location
		if err := rows.Scan(&loc.VehicleID, &loc.Latitude, &loc.Longitude, &loc.Timestamp); err != nil {
			return nil, err
		}
		locations = append(locations, loc)
	}
	return locations, rows.Err()
}
