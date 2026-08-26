package payments

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

type PaymentRepository interface {
	CreatePayment(ctx context.Context, bookingID int64, amount decimal.Decimal, idempotencyKey string) (*PaymentDetails, error)
	ConfirmPayment(ctx context.Context, bookingID int64, paymentID string, orderID string) error
}
type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPaymentRepository(pool *pgxpool.Pool) PaymentRepository {
	return &PostgresRepository{
		pool: pool,
	}
}
func CreatePayment(ctx context.Context, bookingID int64, amount decimal.Decimal, idempotencyKey string) (*PaymentDetails, error) {
	queryBooking := `
		UPDATE bookings
		SET Status=
	`
}
