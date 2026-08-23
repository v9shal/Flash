package bookings

import (
	"time"

	"github.com/shopspring/decimal"
)

type BookingStatus string

const (
	StatusPending   BookingStatus = "pending"
	StatusConfirmed BookingStatus = "confirmed"
	StatusFailed    BookingStatus = "failed"
	StatusCancelled BookingStatus = "cancelled"
	StatusExpired   BookingStatus = "expired"
)

type Booking struct {
	ID             int64           `json:"id"`
	UserID         int64           `json:"user_id"`
	EventID        int64           `json:"event_id"`
	Status         BookingStatus   `json:"status"`
	TotalAmount    decimal.Decimal `json:"total_amount"`
	ExpiresAt      time.Time       `json:"expires_at"`
	IdempotencyKey string          `json:"idempotency_key"`
	CreatedAt      time.Time       `json:"created_at"`
}
