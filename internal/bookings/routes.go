package bookings

import (
	"flash-ticket/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func BookingRoutes(router *gin.Engine, pool *pgxpool.Pool, rdb *redis.Client, jwtSecret string) {
	repo := NewBookingRepository(pool)
	service := NewBookingService(repo, rdb)
	handler := NewBookingHandler(service)
	bookingGroup := router.Group("/bookings")
	bookingGroup.POST("", middleware.Authenticate(jwtSecret), handler.HoldSeat)
}
