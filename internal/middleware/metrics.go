package middleware

import (
	"flash-ticket/internal/platform/metrics"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func Metric(c *gin.Context) {
	start := time.Now()
	c.Next()
	duration := time.Since(start).Seconds()
	metrics.HTTPRequestDuration.WithLabelValues(c.Request.Method, c.FullPath()).Observe(duration)
	metrics.HTTPRequestsTotal.WithLabelValues(c.Request.Method, c.FullPath(), strconv.Itoa(c.Writer.Status())).Inc()
}
