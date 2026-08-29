package payments

import (
	"time"

	"github.com/shopspring/decimal"
)

type PaymentStatus string

const (
	StatusInitiated PaymentStatus = "initiated"
	StatusSuccess   PaymentStatus = "success"
	StatusFailed    PaymentStatus = "failed"
	StatusRefunded  PaymentStatus = "refunded"
)

type PaymentDetails struct {
	ID                int64           `json:"id"`
	BookingId         int64           `json:"booking_id"`
	RazorpayOrderID   *string         `json:"razorpay_order_id"`
	RazorpayPaymentID *string         `json:"razorpay_payment_id"`
	Status            PaymentStatus   `json:"status"`
	Amount            decimal.Decimal `json:"amount"`
	IdempotencyKey    string          `json:"idempotency_key"`
	PaidAt            *time.Time      `json:"paid_at"`
	CreatedAt         time.Time       `json:"created_at"`
}

type PaymentNotes struct {
	BookingID  string `json:"booking_id"`
	EventID    string `json:"event_id"`
	SeatNumber string `json:"seat_number"`
	UserID     string `json:"user_id"`
}
type PaymentEntitiy struct {
	ID      string       `json:"id"`
	OrderID string       `json:"order_id"`
	Amount  int64        `json:"amount"`
	Status  string       `json:"status"`
	Notes   PaymentNotes `json:"notes"`
}
type PaymentItem struct {
	Entity PaymentEntitiy `json:"entity"`
}
type WebhookPayload struct {
	Payment PaymentItem `json:"payment"`
}
type RazorpayWebhookEvent struct {
	Event   string         `json:"event"`
	Payload WebhookPayload `json:"payload"`
}
