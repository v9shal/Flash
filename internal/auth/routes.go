package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(router *gin.Engine, pool *pgxpool.Pool) {
	repo := NewUserRepository(pool)
	service := NewAuthService(repo)
	handler := NewAuthHandler(service)
	authGroup := router.Group("/auth")

	authGroup.POST("/register", handler.Register)

}
