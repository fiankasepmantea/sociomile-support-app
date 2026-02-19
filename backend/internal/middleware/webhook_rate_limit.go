package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"backend/internal/cache"
)

func WebhookRateLimit(r *cache.Redis, limit int) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Redis optional
		if r == nil {
			c.Next()
			return
		}

		ip := c.ClientIP()
		key := "wh:" + ip

		n, err := r.Client.Incr(c, key).Result()
		if err != nil {
			c.Next()
			return
		}

		if n == 1 {
			r.Client.Expire(c, key, time.Minute)
		}

		if n > int64(limit) {
			c.AbortWithStatusJSON(429, gin.H{"error": "rate limit exceeded"})
			return
		}

		c.Next()
	}
}
