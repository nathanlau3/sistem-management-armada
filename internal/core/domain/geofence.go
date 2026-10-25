package domain

import "math"

const EventGeofenceEntry = "geofence_entry"

type Point struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Geofence struct {
	Center Point
	Radius float64
}

func (g Geofence) Contains(p Point) bool {
	return Distance(g.Center, p) <= g.Radius
}

type GeofenceEvent struct {
	VehicleID string `json:"vehicle_id"`
	Event     string `json:"event"`
	Location  Point  `json:"location"`
	Timestamp int64  `json:"timestamp"`
}

func Distance(a, b Point) float64 {
	const r = 6371000.0
	lat1 := a.Latitude * math.Pi / 180
	lat2 := b.Latitude * math.Pi / 180
	dLat := lat2 - lat1
	dLng := (b.Longitude - a.Longitude) * math.Pi / 180

	h := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1)*math.Cos(lat2)*math.Sin(dLng/2)*math.Sin(dLng/2)
	return 2 * r * math.Asin(math.Sqrt(h))
}
