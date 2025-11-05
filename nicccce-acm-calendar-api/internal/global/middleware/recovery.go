package middleware

import (
	"github.com/gin-gonic/gin"
	"nicccce-acm-calendar-api/internal/global/response"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer response.Recovery(c)
		c.Next()
	}
}
