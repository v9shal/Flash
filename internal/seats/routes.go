package seats

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func SeatRoutes(router *gin.Engine, pool *pgxpool.Pool, rdb *redis.Client) {
	repo := NewSeatRepository(pool)
	service := NewSeatService(repo, rdb)
	handler := NewSeatHandler(service)
	seatGroup := router.Group("/seats")

	seatGroup.GET("/:eventId", handler.GetSeats)

}
