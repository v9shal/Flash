package main

import (
	"context"
	"flash-ticket/internal/auth"
	"flash-ticket/internal/bookings"
	"flash-ticket/internal/config"
	"flash-ticket/internal/jobs"
	"flash-ticket/internal/platform/database"
	"flash-ticket/internal/platform/redis"
	"flash-ticket/internal/seats"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	ctx := context.Background()
	pool, err := database.NewPostgresPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	client, err := redis.NewRedisClient(ctx, cfg.RedisURL)
	if err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}
	defer client.Close()
	defer pool.Close()
	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "OK",
		})
	})
	auth.RegisterRoutes(router, pool, cfg.JwtSecret)
	seats.SeatRoutes(router, pool, client)
	bookings.BookingRoutes(router, pool, client, cfg.JwtSecret)
	bookingRepo := bookings.NewBookingRepository(pool)
	jobs.StartSeatExpiryWorker(ctx, bookingRepo, 30*time.Second)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
