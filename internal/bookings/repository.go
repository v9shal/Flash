package bookings

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

var ErrSeatConflit = errors.New("seat is no longer available or version mismatch")

type BookingRepository interface {
	CreateHoldBooking(ctx context.Context, userID int64, eventID int64, seatID int64, expectedVersion int, price decimal.Decimal, idempotencyKey string, expiresAt time.Time) (*Booking, error)
}

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewBookingRepository(pool *pgxpool.Pool) BookingRepository {
	return &PostgresRepository{
		pool: pool,
	}
}

func (r *PostgresRepository) CreateHoldBooking(ctx context.Context, userID int64, eventID int64, seatID int64, expectedVersion int, price decimal.Decimal, idempotencyKey string, expiresAt time.Time) (*Booking, error) {
	tx, err := r.pool.Begin(ctx)
	var booking Booking
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	updateQuery := `
    UPDATE seats
    SET status = 'held',
        held_by = $1,
        expires_at = $2,
        version = version + 1
    WHERE id = $3
      AND status = 'available'
      AND version = $4
`
	cmdTag, err := tx.Exec(ctx, updateQuery, userID, expiresAt, seatID, expectedVersion)
	if err != nil {
		return nil, err
	}

	// THE BACKSTOP: Check if another writer took the seat
	if cmdTag.RowsAffected() == 0 {
		return nil, ErrSeatConflit // Sentinel error: Seat already held or version mismatch!
	}

	insertBooking := `
	
	INSERT INTO bookings(
		user_id,
		event_id,
		status,
		total_amount,
		expires_at,
		idempotency_key
	)
		VALUES($1,$2,$3,$4,$5,$6)
		RETURNING id,created_at
	`
	err = tx.QueryRow(ctx, insertBooking,
		userID,
		eventID,
		string(StatusPending),
		price,
		expiresAt,
		idempotencyKey,
	).Scan(&booking.ID, &booking.CreatedAt)
	if err != nil {
		return nil, err
	}
	booking.UserID = userID
	booking.EventID = eventID
	booking.Status = StatusPending
	booking.TotalAmount = price
	booking.ExpiresAt = expiresAt
	booking.IdempotencyKey = idempotencyKey
	insertBookingSeat := `
		INSERT INTO booking_seats(
			booking_id,
			seat_id,
			price_at_booking
		)
			VALUES($1,$2,$3)
	`
	_, err = tx.Exec(ctx, insertBookingSeat, booking.ID, seatID, price)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return &booking, nil

}
