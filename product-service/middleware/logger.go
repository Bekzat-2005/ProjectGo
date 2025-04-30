package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"time"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		requestID := uuid.New().String()
		c.Set("RequestID", requestID)

		// Выполняем обработку запроса
		c.Next()

		duration := time.Since(start)
		status := c.Writer.Status()

		log.Printf("[%s] [RequestID: %s] %s %s - %d - Duration: %v",
			start.Format(time.RFC3339),
			requestID,
			c.Request.Method,
			c.Request.URL.Path,
			status,
			duration,
		)
	}
}
