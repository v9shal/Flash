package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Authenticate(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {

		// 1. Extract Authorization header
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{"error": "authorization header missing"},
			)
			return
		}

		// 2. Check Bearer token
		parts := strings.SplitN(authHeader, " ", 2)

		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{"error": "invalid authorization header format"},
			)
			return
		}

		tokenString := parts[1]

		// 3. Parse and verify token
		token, err := jwt.Parse(
			tokenString,
			func(token *jwt.Token) (interface{}, error) {

				// Make sure the token uses the expected signing algorithm
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf(
						"unexpected signing method: %v",
						token.Header["alg"],
					)
				}

				return []byte(jwtSecret), nil
			},
		)

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{"error": "invalid or expired token"},
			)
			return
		}

		// 4. Extract claims
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			// jwt.MapClaims parses numbers as float64 by default
			userID := int64(claims["user_id"].(float64))

			// Attach it to the Gin context for the next function to use!
			c.Set("userID", userID)
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
			return
		}

		c.Next()
	}
}
