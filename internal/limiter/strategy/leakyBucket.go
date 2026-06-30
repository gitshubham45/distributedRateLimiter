package strategy

import (
	"rate-limiter/internal/config"

	"github.com/gin-gonic/gin"
)

type LeakyBucket struct{}

func (l *LeakyBucket) Allow(c *gin.Context, config config.StrategyConfig) bool {
	return true
}
