package payments

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func PaymentRoutes(router *gin.Engine, pool *pgxpool.Pool, rdb *redis.Client, webhookSecret string) {
	repo := NewPaymentRepository(pool)
	service := NewPaymentService(repo, rdb, webhookSecret)
	handler := NewPaymenthandler(service)
	paymentGroup := router.Group("/payments")
	paymentGroup.POST("/webhook", handler.HandleWebhook)
}
