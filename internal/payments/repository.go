package payments

import (
	"context"
	"fmt"
	"time"

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
func (r *PostgresRepository) CreatePayment(ctx context.Context, bookingID int64, amount decimal.Decimal, idempotencyKey string) (*PaymentDetails, error) {
	tx, err := r.pool.Begin(ctx)
	var payment PaymentDetails
	if err != nil {

		return nil, err
	}
	defer tx.Rollback(ctx)

	insertPayment := `
		INSERT INTO Payments (
		booking_id,
		status,
		amount,
		idempotency_key
		created_at		
		)
		VALUES($1,$2,$3,$4,$5)
		RETURNING id,paid_at
	`
	err = tx.QueryRow(ctx, insertPayment, bookingID, string(StatusInitiated), amount, idempotencyKey, time.Now()).Scan(&payment.ID, &payment.PaidAt)
	if err != nil {
		return nil, err
	}
	payment.BookingId = bookingID
	payment.Amount = amount
	payment.IdempotencyKey = idempotencyKey
	payment.Status = StatusInitiated
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return &payment, nil
}
func (r *PostgresRepository) ConfirmPayment(
	ctx context.Context,
	bookingID int64,
	paymentID string,
	orderID string,
) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	var status string
	checkpayment := `
	SELECT status from payments where razorpay_payment_id=$1
	`
	err = tx.QueryRow(ctx, checkpayment, paymentID).Scan(&status)
	if err != nil {
		return err
	}

	if status == "success" {
		fmt.Println("already paid and booked")
		return nil
	}

	bookingQuery := `
		UPDATE bookings
		SET status = 'confirmed',
		WHERE id = $1
		  AND status = 'pending'
	`

	result, err := tx.Exec(ctx, bookingQuery, bookingID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("booking %d not found or already confirmed", bookingID)
	}

	seatQuery := `
		UPDATE seats s
		SET
			status = 'booked',
			expires_at = NULL,
			version = version + 1
		FROM booking_seats bs
		WHERE bs.booking_id = $1
		  AND bs.seat_id = s.id
		  AND s.status = 'held'
		  AND s.held_by = (
			  SELECT user_id
			  FROM bookings
			  WHERE id = $1
		  )
	`

	result, err = tx.Exec(ctx, seatQuery, bookingID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("no held seats found for booking %d", bookingID)
	}

	paymentQuery := `
		UPDATE payments
		SET
			status = 'success',
			razorpay_payment_id = $1,
			razorpay_order_id = $2,
			paid_at = NOW()
		WHERE booking_id = $3
		  AND status != 'paid'
	`

	result, err = tx.Exec(
		ctx,
		paymentQuery,
		paymentID,
		orderID,
		bookingID,
	)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("payment not found or already paid for booking %d", bookingID)
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}
