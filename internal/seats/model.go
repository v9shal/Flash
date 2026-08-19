package seats

import (
	"time"

	"github.com/shopspring/decimal"
)

type SeatStatus string

const (
	StatusAvailable SeatStatus = "available"
	StatusHeld      SeatStatus = "held"
	StatusBooked    SeatStatus = "booked"
	StatusBlocked   SeatStatus = "blocked"
)

type Seat struct {
	ID         int64           `json:"id"`
	EventID    int64           `json:"event_id"`
	SeatNumber string          `json:"seat_number"`
	Price      decimal.Decimal `json:"price"`
	Status     SeatStatus      `json:"status"`
	HeldBy     *int64          `json:"held_by"`
	ExpiresAt  *time.Time      `json:"expires_at"`
	Version    int             `json:"version"`
}
