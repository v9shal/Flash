package middleware

import (
	redisPlatform "flash-ticket/internal/platform/redis"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func RateLimiter(rdb *redis.Client, prefix string, limit int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 2. THIS inner function runs on EVERY incoming HTTP request!
		// Gin automatically passes the request's Context as the variable `c`.

		ip := c.ClientIP() // <--- `c` is the parameter of this inner function

		// ... your rate limiting logic here ...
		key := fmt.Sprintf("rate:%s:%s", prefix, ip)
		now := time.Now().UnixMilli()
		windowStart := now - window.Milliseconds()
		windowSec := int(window.Seconds())
		allowed, err := redisPlatform.RunRateLimiterScript(c.Request.Context(), rdb, key, now, windowStart, limit, windowSec)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error while check the requrest"})
			return
		}

		if !allowed {
			c.Header("Retry-After", strconv.Itoa(windowSec))
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "too many requests, please try again later"})
			return
		}
		c.Next()

	}

}
