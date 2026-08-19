package auth

import (
	"flash-ticket/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(router *gin.Engine, pool *pgxpool.Pool, jwtSecret string) {
	repo := NewUserRepository(pool)
	service := NewAuthService(repo, jwtSecret)
	handler := NewAuthHandler(service, jwtSecret)
	authGroup := router.Group("/auth")

	authGroup.POST("/register", handler.Register)
	authGroup.POST("/login", handler.Login)
	authGroup.GET("/me", middleware.Authenticate(jwtSecret), handler.Me)

}
