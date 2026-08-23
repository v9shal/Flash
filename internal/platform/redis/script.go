package redis

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/redis/go-redis/v9"
)

//go:embed scripts/hold_seat.lua
var holdSeatLua string

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
