package domain

import "errors"

var ErrNotFound = errors.New("not found")

type Location struct {
	VehicleID string  `json:"vehicle_id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timestamp int64   `json:"timestamp"`
}

func (l Location) Validate() error {
	if l.VehicleID == "" {
		return errors.New("vehicle_id is required")
	}
	if l.Latitude < -90 || l.Latitude > 90 {
		return errors.New("invalid latitude")
	}
	if l.Longitude < -180 || l.Longitude > 180 {
		return errors.New("invalid longitude")
	}
	if l.Timestamp <= 0 {
		return errors.New("invalid timestamp")
	}
	return nil
}

func (l Location) Point() Point {
	return Point{Latitude: l.Latitude, Longitude: l.Longitude}
}
