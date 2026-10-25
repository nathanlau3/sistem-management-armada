package service

import (
	"context"
	"errors"

	"fleet/internal/core/domain"
	"fleet/internal/core/port"
)

type LocationService struct {
	repo      port.LocationRepository
	publisher port.EventPublisher
	geofence  domain.Geofence
}

func NewLocationService(repo port.LocationRepository, publisher port.EventPublisher, geofence domain.Geofence) *LocationService {
	return &LocationService{repo: repo, publisher: publisher, geofence: geofence}
}

func (s *LocationService) Record(ctx context.Context, loc domain.Location) error {
	if err := loc.Validate(); err != nil {
		return err
	}

	prev, err := s.repo.Latest(ctx, loc.VehicleID)
	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		return err
	}
	wasInside := err == nil && s.geofence.Contains(prev.Point())

	if err := s.repo.Save(ctx, loc); err != nil {
		return err
	}

	if wasInside || !s.geofence.Contains(loc.Point()) {
		return nil
	}

	return s.publisher.PublishGeofenceEvent(ctx, domain.GeofenceEvent{
		VehicleID: loc.VehicleID,
		Event:     domain.EventGeofenceEntry,
		Location:  loc.Point(),
		Timestamp: loc.Timestamp,
	})
}

func (s *LocationService) Latest(ctx context.Context, vehicleID string) (domain.Location, error) {
	return s.repo.Latest(ctx, vehicleID)
}

func (s *LocationService) History(ctx context.Context, vehicleID string, start, end int64) ([]domain.Location, error) {
	return s.repo.History(ctx, vehicleID, start, end)
}
