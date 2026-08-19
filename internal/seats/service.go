package seats

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/singleflight"
)

type SeatService interface {
	GetSeats(ctx context.Context, eventId int64) ([]Seat, error)
}

type seatService struct {
	repo  SeatRepository
	redis *redis.Client
	sf    *singleflight.Group
}

func NewSeatService(repo SeatRepository, rd *redis.Client) SeatService {
	return &seatService{
		repo:  repo,
		redis: rd,
		sf:    &singleflight.Group{},
	}
}

func (s *seatService) GetSeats(ctx context.Context, eventID int64) ([]Seat, error) {
	cacheKey := fmt.Sprintf("event:%d:seats", eventID)
	val, err := s.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		seats := make([]Seat, 0)
		if err := json.Unmarshal([]byte(val), &seats); err == nil {

			return seats, err
		}
	}

	res, err, _ := s.sf.Do(cacheKey, func() (interface{}, error) {

		seats, err := s.repo.GetSeatsByEventID(ctx, eventID)
		if err != nil {
			return nil, err

		}
		jsonData, err := json.Marshal(seats)
		if err == nil {
			s.redis.Set(ctx, cacheKey, jsonData, 5*time.Minute)
		}
		return seats, nil
	})
	if err != nil {
		return nil, err
	}
	seats, ok := res.([]Seat)
	if !ok {
		return nil, fmt.Errorf("unexpected data type from singleflight")

	}
	return seats, nil

}
