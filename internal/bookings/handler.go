package bookings

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

type HoldBookingRequest struct {
	EventID         int64           `json:"event_id" binding:"required"`
	SeatID          int64           `json:"seat_id" binding:"required"`
	SeatNumber      string          `json:"seat_number" binding:"required"`
	ExpectedVersion int             `json:"expected_version" binding:"required,min=1"`
	Price           decimal.Decimal `json:"price" binding:"required"`
	IdempotencyKey  string          `json:"idempotency_key" binding:"required"`
}

type BookingHandler struct {
	service BookingService
}

func NewBookingHandler(service BookingService) *BookingHandler {
	return &BookingHandler{service: service}
}

func (h *BookingHandler) HoldSeat(c *gin.Context) {
	var req HoldBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDValue, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user id not found in context"})
		return
	}

	userID, ok := userIDValue.(int64)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user id in context"})
		return
	}

	booking, err := h.service.HoldSeat(
		c.Request.Context(),
		userID,
		req.EventID,
		req.SeatID,
		req.SeatNumber,
		req.ExpectedVersion,
		req.Price,
		req.IdempotencyKey,
	)
	if err != nil {
		switch {
		case errors.Is(err, ErrSeatAlreadyHeld), errors.Is(err, ErrUserAlreadyHasHold), errors.Is(err, ErrSeatConflit):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "seat held", "booking": booking})
}
