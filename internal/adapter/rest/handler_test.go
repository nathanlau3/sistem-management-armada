package rest

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"fleet/internal/core/domain"
)

type fakeService struct{}

func (fakeService) Record(context.Context, domain.Location) error { return nil }

func (fakeService) Latest(_ context.Context, vehicleID string) (domain.Location, error) {
	if vehicleID != "B1234ABC" {
		return domain.Location{}, domain.ErrNotFound
	}
	return domain.Location{VehicleID: vehicleID, Latitude: -6.2088, Longitude: 106.8456, Timestamp: 1715003456}, nil
}

func (fakeService) History(context.Context, string, int64, int64) ([]domain.Location, error) {
	return []domain.Location{}, nil
}

func TestRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	NewHandler(fakeService{}).Routes(r)

	tests := []struct {
		target string
		want   int
	}{
		{"/vehicles/B1234ABC/location", http.StatusOK},
		{"/vehicles/UNKNOWN/location", http.StatusNotFound},
		{"/vehicles/B1234ABC/history?start=1715000000&end=1715009999", http.StatusOK},
		{"/vehicles/B1234ABC/history?start=abc&end=1715009999", http.StatusBadRequest},
		{"/vehicles/B1234ABC/history?start=1715009999", http.StatusBadRequest},
		{"/vehicles/B1234ABC/history?start=10&end=1", http.StatusBadRequest},
	}

	for _, tt := range tests {
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, tt.target, nil))
		if rec.Code != tt.want {
			t.Errorf("%s: got %d, want %d", tt.target, rec.Code, tt.want)
		}
	}
}
