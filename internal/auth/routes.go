package auth

import (
	"flash-ticket/internal/middleware"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func RegisterRoutes(router *gin.Engine, pool *pgxpool.Pool, rdb *redis.Client, jwtSecret string) {
	repo := NewUserRepository(pool)
	service := NewAuthService(repo, jwtSecret)
	handler := NewAuthHandler(service, jwtSecret)
	authGroup := router.Group("/auth")

	authGroup.POST("/register", handler.Register)
	authGroup.POST("/login", middleware.RateLimiter(rdb, "login", 5, 1*time.Minute), handler.Login)
	authGroup.GET("/me", middleware.Authenticate(jwtSecret), handler.Me)

}
