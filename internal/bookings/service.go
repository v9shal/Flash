package bookings

import (
	"context"
	"errors"
	"fmt"
	"time"

	redisPlatform "flash-ticket/internal/platform/redis"

	"github.com/redis/go-redis/v9"
	"github.com/shopspring/decimal"
)

var (
	ErrSeatAlreadyHeld    = errors.New("seat is already held by another user")
	ErrUserAlreadyHasHold = errors.New("user already has an active hold")
)

type BookingService interface {
	HoldSeat(ctx context.Context, userID int64, eventID int64, seatID int64, seatNumber string, expectedVersion int, price decimal.Decimal, idempotencyKey string) (*Booking, error)
}
type bookingService struct {
	repo  BookingRepository
	redis *redis.Client
}

func NewBookingService(
	repo BookingRepository,
	redis *redis.Client,
) BookingService {
	return &bookingService{
		repo:  repo,
		redis: redis,
	}
}

func (s *bookingService) HoldSeat(ctx context.Context, userID int64, eventID int64, seatID int64, seatNumber string, expectedVersion int, price decimal.Decimal, idempotencyKey string) (*Booking, error) {
	result, err := redisPlatform.RunHoldSeatScript(ctx, s.redis, eventID, seatNumber, userID)
	if err != nil {
		return nil, err
	}
	if result == 0 {
		return nil, ErrSeatAlreadyHeld
	} else if result == -1 {
		return nil, ErrUserAlreadyHasHold
	}
	expirayTime := time.Now().Add(10 * time.Minute)
	booking, err := s.repo.CreateHoldBooking(ctx, userID, eventID, seatID, expectedVersion, price, idempotencyKey, expirayTime)
	if err != nil {
		seatHoldKey := fmt.Sprintf("hold:%d:%s", eventID, seatNumber)
		userHoldKey := fmt.Sprintf("user:hold:%d", userID)
		s.redis.Del(ctx, seatHoldKey, userHoldKey)
		return nil, err

	}
	return booking, nil
}
