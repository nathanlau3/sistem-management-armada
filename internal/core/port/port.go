package port

import (
	"context"

	"fleet/internal/core/domain"
)

type LocationService interface {
	Record(ctx context.Context, loc domain.Location) error
	Latest(ctx context.Context, vehicleID string) (domain.Location, error)
	History(ctx context.Context, vehicleID string, start, end int64) ([]domain.Location, error)
}

type LocationRepository interface {
	Save(ctx context.Context, loc domain.Location) error
	Latest(ctx context.Context, vehicleID string) (domain.Location, error)
	History(ctx context.Context, vehicleID string, start, end int64) ([]domain.Location, error)
}

type EventPublisher interface {
	PublishGeofenceEvent(ctx context.Context, event domain.GeofenceEvent) error
}
