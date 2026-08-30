package main

import (
	"context"
	"errors"
	"flash-ticket/internal/auth"
	"flash-ticket/internal/bookings"
	"flash-ticket/internal/config"
	"flash-ticket/internal/jobs"
	"flash-ticket/internal/middleware"
	"flash-ticket/internal/payments"
	"flash-ticket/internal/platform/database"
	"flash-ticket/internal/platform/redis"
	"flash-ticket/internal/seats"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/gin-gonic/gin"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
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
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(
		middleware.RequestID, // (or middleware.RequestID)
		middleware.StructuredLogger(),
		middleware.Metric,
	)
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "OK",
		})
	})
	auth.RegisterRoutes(router, pool, client, cfg.JwtSecret)
	seats.SeatRoutes(router, pool, client)
	bookings.BookingRoutes(router, pool, client, cfg.JwtSecret)
	bookingRepo := bookings.NewBookingRepository(pool)
	payments.PaymentRoutes(router, pool, client, cfg.PaymentSecret)
	jobs.StartSeatExpiryWorker(ctx, bookingRepo, 30*time.Second)
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen error: %s\n", err)
		}
	}()
	<-ctx.Done()
	slog.Info("Shutdown signal received, shutting down gracefully...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
	}

	slog.Info("Server exited cleanly")
}
