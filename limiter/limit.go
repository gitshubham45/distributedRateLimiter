package limiter

import "github.com/gin-gonic/gin"

type RequestLimiter interface {
	Allow(c *gin.Context) bool
}

type PerHourLimit struct{}

func (l *PerHourLimit) Allow(c *gin.Context) bool {
	return true
}
