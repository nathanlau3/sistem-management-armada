package rest

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"fleet/internal/core/domain"
	"fleet/internal/core/port"
)

type Handler struct {
	svc port.LocationService
}

func NewHandler(svc port.LocationService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Routes(r *gin.Engine) {
	r.GET("/vehicles/:vehicle_id/location", h.latestLocation)
	r.GET("/vehicles/:vehicle_id/history", h.history)
}

func (h *Handler) latestLocation(c *gin.Context) {
	loc, err := h.svc.Latest(c.Request.Context(), c.Param("vehicle_id"))
	if errors.Is(err, domain.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "location not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, loc)
}

func (h *Handler) history(c *gin.Context) {
	start, err1 := strconv.ParseInt(c.Query("start"), 10, 64)
	end, err2 := strconv.ParseInt(c.Query("end"), 10, 64)
	if err1 != nil || err2 != nil || start > end {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start or end"})
		return
	}

	locations, err := h.svc.History(c.Request.Context(), c.Param("vehicle_id"), start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, locations)
}
