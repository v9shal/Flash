package auth

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 1. The Request Struct (Your Zod equivalent)
type RegisterRequest struct {
	FullName string  `json:"full_name" binding:"required,min=2,max=100"`
	Email    string  `json:"email" binding:"required,email"`
	Password string  `json:"password" binding:"required,min=12,max=72"`
	Phone    *string `json:"phone" binding:"omitempty,e164"` // e164 is standard international phone format
}
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// 2. The Handler Struct
type AuthHandler struct {
	service   AuthService
	jwtSecret string
}

func NewAuthHandler(service AuthService, jwtSecret string) *AuthHandler {
	return &AuthHandler{service: service, jwtSecret: jwtSecret}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := h.service.RegisterUser(c.Request.Context(), req.FullName, req.Email, req.Password, req.Phone)
	if err != nil {
		if errors.Is(err, ErrDuplicateUser) {
			c.JSON(http.StatusConflict, gin.H{"error": "user with this email or phone already exists"})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
	}
	c.JSON(http.StatusCreated, gin.H{"message": "user created", "user": user})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, accessToken, refreshToken, err := h.service.LoginUser(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid email  or password"})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
	}
	c.SetCookie("refresh_token", refreshToken, 7*24*60*60, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "login successful", "accesstoken": accessToken, "user": user})

}
