package seats

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SeatRepository interface {
	GetSeatsByEventID(ctx context.Context, eventID int64) ([]Seat, error)
}

type postgresRepository struct {
	pool *pgxpool.Pool
}

func NewSeatRepository(pool *pgxpool.Pool) SeatRepository {
	return &postgresRepository{
		pool: pool,
	}
}

func (r *postgresRepository) GetSeatsByEventID(ctx context.Context, eventID int64) ([]Seat, error) {
	query := `
		SELECT id, event_id, seat_number, price, status, held_by, expires_at, version FROM seats WHERE event_id = $1 ORDER BY id ASC
	`

	rows, err := r.pool.Query(ctx, query, eventID)
	if err != nil {
		return nil, fmt.Errorf("query seats by event: %w", err)
	}
	defer rows.Close()
	seats := make([]Seat, 0)
	for rows.Next() {
		var seat Seat
		if err := rows.Scan(&seat.ID, &seat.EventID, &seat.SeatNumber, &seat.Price, &seat.Status, &seat.HeldBy, &seat.ExpiresAt, &seat.Version); err != nil {
			return nil, err
		}
		seats = append(seats, seat)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return seats, nil
}
