package service

import (
	"context"
	"fleet/internal/core/domain"
	"fleet/internal/core/port"
	"testing"
)

type memRepo struct {
	data []domain.Location
}

type mockEventPublisher struct {
	events []domain.GeofenceEvent
}

var _ port.LocationRepository = (*memRepo)(nil)

var _ port.EventPublisher = (*mockEventPublisher)(nil)

func (m *memRepo) Save(_ context.Context, location domain.Location) error {
	m.data = append(m.data, location)
	return nil
}

func (m *memRepo) Latest(_ context.Context, vehicleID string) (domain.Location, error) {
	for i := len(m.data) - 1; i >= 0; i-- {
		if m.data[i].VehicleID == vehicleID {
			return m.data[i], nil
		}
	}
	return domain.Location{}, domain.ErrNotFound
}

func (m *memRepo) History(_ context.Context, vehicleID string, start, end int64) ([]domain.Location, error) {
	var out []domain.Location
	for _, l := range m.data {
		if l.VehicleID == vehicleID && l.Timestamp >= start && l.Timestamp <= end {
			out = append(out, l)
		}
	}
	return out, nil
}

func (e *mockEventPublisher) PublishGeofenceEvent(_ context.Context, geofenceEvent domain.GeofenceEvent) error {
	e.events = append(e.events, geofenceEvent)
	return nil
}

func TestRecordPublishesOnlyOnEntry(t *testing.T) {
	repo := &memRepo{}
	pub := &mockEventPublisher{}
	fence := domain.Geofence{
		Center: domain.Point{
			Latitude:  -6.2088,
			Longitude: 106.8456,
		},
		Radius: 50,
	}
	svc := NewLocationService(repo, pub, fence)

	steps := []struct {
		lat, lng float64
	}{
		{-6.2100, 106.8456},
		{-6.2092, 106.8456},
		{-6.2088, 106.8456},
		{-6.2100, 106.8456},
		{-6.2090, 106.8456},
	}

	for i, s := range steps {
		location := domain.Location{
			VehicleID: "B1234XYZ",
			Latitude:  s.lat,
			Longitude: s.lng,
			Timestamp: int64(1715003456 + i),
		}
		if err := svc.Record(context.Background(), location); err != nil {
			t.Fatalf("Record error at step %d: %v", i, err)
		}
	}

	if len(repo.data) != len(steps) {
		t.Fatalf("expected %d saved locations, got %d", len(steps), len(repo.data))
	}
	if len(pub.events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(pub.events))
	}
	if pub.events[0].Timestamp != 1715003457 {
		t.Errorf("expected first event at step 1, got %+v", pub.events[0])
	}
	if pub.events[0].Event != domain.EventGeofenceEntry {
		t.Errorf("expected event %q, got %q", domain.EventGeofenceEntry, pub.events[0].Event)
	}
}
