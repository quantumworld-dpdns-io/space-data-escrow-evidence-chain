package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func APIKey(required string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if required == "" {
			c.Next()
			return
		}
		if c.GetHeader("X-API-Key") != required {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		c.Next()
	}
}
