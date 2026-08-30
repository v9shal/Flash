package redis

import (
	"context"
	_ "embed"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

//go:embed scripts/hold_seat.lua
var holdSeatLua string

//go:embed scripts/rate_limiter.lua
var rateLimtierLua string
var RateLimiter = redis.NewScript(rateLimtierLua)
var HoldSeatScript = redis.NewScript(holdSeatLua)

func RunHoldSeatScript(ctx context.Context, rdb *redis.Client, eventID int64, seatNumber string, userID int64) (int, error) {
	seatHoldKey := fmt.Sprintf("hold:%d:%s", eventID, seatNumber)
	userHoldKey := fmt.Sprintf("user:hold:%d", userID)

	result, err := HoldSeatScript.Run(ctx, rdb, []string{seatHoldKey, userHoldKey}, userID).Int()
	if err != nil {
		return 0, fmt.Errorf("execute hold seat script: %w", err)
	}

	return result, nil
}
func RunRateLimiterScript(ctx context.Context, rdb *redis.Client, key string, now int64, windowStart int64,
	limit int, windowSec int) (bool, error) {
	member := fmt.Sprintf("%d", time.Now().UnixNano())
	val, err := RateLimiter.Run(ctx, rdb, []string{key}, now, windowStart, limit, windowSec, member).Int()
	if err != nil {
		return false, fmt.Errorf("execute rate limiter script: %w", err)
	}
	return val == 1, nil
}
