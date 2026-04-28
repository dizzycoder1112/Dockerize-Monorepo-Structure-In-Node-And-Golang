package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method
		rawQuery := c.Request.URL.RawQuery
		ip := c.ClientIP()

		c.Next()

		durationMs := float64(time.Since(start).Nanoseconds()) / 1e6
		statusCode := c.Writer.Status()

		query := ""
		if rawQuery != "" {
			query = " ?" + rawQuery
		}
		log.Printf("[%d] %s %s%s ip=%s dur=%.2fms", statusCode, method, path, query, ip, durationMs)
	}
}
