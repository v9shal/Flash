package seats

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SeatHandler struct {
	service SeatService
}

func NewSeatHandler(service SeatService) *SeatHandler {
	return &SeatHandler{
		service: service,
	}
}
func (h *SeatHandler) GetSeats(c *gin.Context) {
	eventIDStr := c.Param("eventId")
	eventID, err := strconv.ParseInt(eventIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event id"})
		return
	}
	seats, err := h.service.GetSeats(c.Request.Context(), eventID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error while getting the seats"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "seats retrieved", "seats": seats})
}
